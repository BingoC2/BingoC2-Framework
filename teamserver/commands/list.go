package commands

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"github.com/fatih/color"
)

var Green = color.New(color.FgGreen)
var Red = color.New(color.FgRed)

// insert value into a at index
func InsertIntoSlice(a []string, index int, value string) []string {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

// get a list of file names from a directory stripping extension
func GetNamesFromDir(directory string) []string {
	var sFiles []string

	files, _ := os.ReadDir(directory)

	for _, file := range files {
		sFiles = append(sFiles, strings.Split(file.Name(), ".")[0])
	}

	return sFiles
}

func ListBeacons(c *grumble.Context) error {
	beaconNames := GetNamesFromDir(version.BEACON_DATA_DIR)

	var sleepSlice []string
	var jitterSlice []string
	var listenerSlice []string
	var opsysSlice []string
	var archSlice []string

	for _, name := range beaconNames {
		beaconData, err := yamlstructs.ReadBeaconYaml(name)
		if err != nil {
			return err
		}

		sleepSlice = append(sleepSlice, strconv.Itoa(beaconData.Sleep))
		jitterSlice = append(jitterSlice, strconv.Itoa(beaconData.Jitter))
		listenerSlice = append(listenerSlice, beaconData.Listener)
		opsysSlice = append(opsysSlice, beaconData.GOOS)
		archSlice = append(archSlice, beaconData.GOARCH)
	}

	beaconNames = InsertIntoSlice(beaconNames, 0, " Name ")
	beaconNames = InsertIntoSlice(beaconNames, 1, "======")

	sleepSlice = InsertIntoSlice(sleepSlice, 0, " Sleep ")
	sleepSlice = InsertIntoSlice(sleepSlice, 1, "=======")

	jitterSlice = InsertIntoSlice(jitterSlice, 0, " Jitter ")
	jitterSlice = InsertIntoSlice(jitterSlice, 1, "========")

	listenerSlice = InsertIntoSlice(listenerSlice, 0, " Listener ")
	listenerSlice = InsertIntoSlice(listenerSlice, 1, "==========")

	opsysSlice = InsertIntoSlice(opsysSlice, 0, " Operating System ")
	opsysSlice = InsertIntoSlice(opsysSlice, 1, "==================")

	archSlice = InsertIntoSlice(archSlice, 0, " Architecture ")
	archSlice = InsertIntoSlice(archSlice, 1, "==============")

	outptuBuffer := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(outptuBuffer, 0, 0, 1, ' ', 0)
	for key := range beaconNames {
		fmt.Fprintln(writer, beaconNames[key]+"\t"+sleepSlice[key]+"\t"+jitterSlice[key]+"%\t"+listenerSlice[key]+"\t"+opsysSlice[key]+"\t"+archSlice[key])
	}
	writer.Flush()
	c.App.Println(outptuBuffer)

	return nil
}

func ListListeners(c *grumble.Context) error {
	listenerNames := GetNamesFromDir(version.LISTENER_DATA_DIR)

	var lhostSlice []string
	var lportSlice []string
	var uriSlice []string
	var isAliveSlice []string
	var chanSlice []string

	for _, name := range listenerNames {
		listenerData, err := yamlstructs.ReadListenerYaml(name)
		if err != nil {
			return err
		}

		lhostSlice = append(lhostSlice, listenerData.LHOST)
		lportSlice = append(lportSlice, listenerData.LPORT)
		uriSlice = append(uriSlice, listenerData.URI)
		isAliveSlice = append(isAliveSlice, strconv.FormatBool(listenerData.IsAlive))
		chanSlice = append(chanSlice, listenerData.Channel)
	}

	listenerNames = InsertIntoSlice(listenerNames, 0, " Name ")
	listenerNames = InsertIntoSlice(listenerNames, 1, "======")

	chanSlice = InsertIntoSlice(chanSlice, 0, " Channel ")
	chanSlice = InsertIntoSlice(chanSlice, 1, "=========")

	lhostSlice = InsertIntoSlice(lhostSlice, 0, " LHOST ")
	lhostSlice = InsertIntoSlice(lhostSlice, 1, "=======")

	lportSlice = InsertIntoSlice(lportSlice, 0, " LPORT ")
	lportSlice = InsertIntoSlice(lportSlice, 1, "=======")

	uriSlice = InsertIntoSlice(uriSlice, 0, " URI ")
	uriSlice = InsertIntoSlice(uriSlice, 1, "=====")

	isAliveSlice = InsertIntoSlice(isAliveSlice, 0, " Alive ")
	isAliveSlice = InsertIntoSlice(isAliveSlice, 1, "=======")

	outptuBuffer := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(outptuBuffer, 0, 0, 1, ' ', 0)
	for key := range listenerNames {
		if key == 0 || key == 1 {
			fmt.Fprintln(writer, listenerNames[key]+"\t"+chanSlice[key]+"\t"+lhostSlice[key]+"\t"+lportSlice[key]+"\t"+uriSlice[key]+"\t"+isAliveSlice[key])
		} else if isAliveSlice[key] == "false" {
			fmt.Fprintln(writer, listenerNames[key]+"\t"+chanSlice[key]+"\t"+lhostSlice[key]+"\t"+lportSlice[key]+"\t"+uriSlice[key]+"\t"+Red.Sprint(isAliveSlice[key]))
		} else if isAliveSlice[key] == "true" {
			fmt.Fprintln(writer, listenerNames[key]+"\t"+chanSlice[key]+"\t"+lhostSlice[key]+"\t"+lportSlice[key]+"\t"+uriSlice[key]+"\t"+Green.Sprint(isAliveSlice[key]))
		} else {
			fmt.Fprintln(writer, listenerNames[key]+"\t"+chanSlice[key]+"\t"+lhostSlice[key]+"\t"+lportSlice[key]+"\t"+uriSlice[key]+"\t"+isAliveSlice[key])
		}
	}
	writer.Flush()
	c.App.Println(outptuBuffer)

	return nil
}

func ListSessions(c *grumble.Context) error {
	sessionAgentIds := GetNamesFromDir(version.SESSION_DATA_DIR)

	var hostnameSlice []string = []string{" Hostname ", "=========="}
	var ipsSlice []string = []string{" IP ", "===="}
	var pidSlice []string = []string{" PID ", "====="}
	var ppidSlice []string = []string{" PPID ", "======"}
	var userSlice []string = []string{" User ", "======"}
	var opsysSlice []string = []string{" Operating System", "=================="}
	var sleepSlice []string = []string{" Sleep ", "======="}
	var jitterSlice []string = []string{" Jitter ", "========"}
	var listenerSlice []string = []string{" Listener ", "=========="}
	var lastCallBackSlice []string = []string{" Last Call Back ", "================"}
	var lastCallBackSliceDetailed []time.Time = []time.Time{time.Now(), time.Now()}

	for _, id := range sessionAgentIds {
		sessionData, err := yamlstructs.ReadSessionYaml(id)
		if err != nil {
			return err
		}

		hostnameSlice = append(hostnameSlice, sessionData.Hostname)
		ipsSlice = append(ipsSlice, sessionData.IP)
		pidSlice = append(pidSlice, strconv.Itoa(sessionData.ProcessID))
		ppidSlice = append(ppidSlice, strconv.Itoa(sessionData.ParentProcessID))
		userSlice = append(userSlice, sessionData.Username)
		opsysSlice = append(opsysSlice, sessionData.OperatingSystem)
		sleepSlice = append(sleepSlice, strconv.Itoa(sessionData.Sleep))
		jitterSlice = append(jitterSlice, strconv.Itoa(sessionData.Jitter))
		listenerSlice = append(listenerSlice, sessionData.Listener)
		lastCallBackSlice = append(lastCallBackSlice, sessionData.LastCallBack.Format("15:04:05"))
		lastCallBackSliceDetailed = append(lastCallBackSliceDetailed, sessionData.LastCallBack)
	}

	sessionAgentIds = InsertIntoSlice(sessionAgentIds, 0, " Agent ID ")
	sessionAgentIds = InsertIntoSlice(sessionAgentIds, 1, "==========")

	outputBuffer := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(outputBuffer, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
	for key := range sessionAgentIds {
		if key == 0 || key == 1 {
			fmt.Fprintln(writer, sessionAgentIds[key]+"\t"+hostnameSlice[key]+"\t"+ipsSlice[key]+"\t"+pidSlice[key]+"\t"+ppidSlice[key]+"\t"+userSlice[key]+"\t"+opsysSlice[key]+"\t"+sleepSlice[key]+"\t"+jitterSlice[key]+"\t"+listenerSlice[key]+"\t"+lastCallBackSlice[key])
		} else {
			// check for dead sessions
			sleepInt, _ := strconv.Atoi(sleepSlice[key])
			jitterInt, _ := strconv.Atoi(jitterSlice[key])
			maxCallBackTime := sleepInt + jitterInt + 1

			sessionData, err := yamlstructs.ReadSessionYaml(sessionAgentIds[key])
			if err != nil {
				return err
			}

			if time.Now().After(sessionData.LastCallBack.Add(time.Duration(maxCallBackTime) * time.Second)) {
				fmt.Fprintln(writer, sessionAgentIds[key]+"\t"+hostnameSlice[key]+"\t"+ipsSlice[key]+"\t"+pidSlice[key]+"\t"+ppidSlice[key]+"\t"+userSlice[key]+"\t"+opsysSlice[key]+"\t"+sleepSlice[key]+"\t"+jitterSlice[key]+"\t"+listenerSlice[key]+"\t"+Red.Sprint(lastCallBackSlice[key])+" ("+fmt.Sprint(int(time.Since(lastCallBackSliceDetailed[key]).Seconds()))+" seconds)")
			} else {
				if userSlice[key] == "root" || strings.Contains(userSlice[key], "SYSTEM") {
					fmt.Fprintln(writer, sessionAgentIds[key]+"\t"+hostnameSlice[key]+"\t"+ipsSlice[key]+"\t"+pidSlice[key]+"\t"+ppidSlice[key]+"\t"+userSlice[key]+"💀\t"+opsysSlice[key]+"\t"+sleepSlice[key]+"\t"+jitterSlice[key]+"\t"+listenerSlice[key]+"\t"+Green.Sprint(lastCallBackSlice[key])+" ("+fmt.Sprint(int(time.Since(lastCallBackSliceDetailed[key]).Seconds()))+" seconds)")
				} else {
					fmt.Fprintln(writer, sessionAgentIds[key]+"\t"+hostnameSlice[key]+"\t"+ipsSlice[key]+"\t"+pidSlice[key]+"\t"+ppidSlice[key]+"\t"+userSlice[key]+"\t"+opsysSlice[key]+"\t"+sleepSlice[key]+"\t"+jitterSlice[key]+"\t"+listenerSlice[key]+"\t"+Green.Sprint(lastCallBackSlice[key])+" ("+fmt.Sprint(int(time.Since(lastCallBackSliceDetailed[key]).Seconds()))+" seconds)")
				}
			}
		}
	}
	writer.Flush()
	c.App.Println(outputBuffer)

	return nil
}

func ListAll(c *grumble.Context) error {
	c.App.Println("=== Listeners ===\n")
	err := ListListeners(c)
	if err != nil {
		return err
	}

	c.App.Println("\n")

	c.App.Println("=== Beacons ===\n")
	err = ListBeacons(c)
	if err != nil {
		return err
	}

	c.App.Println("\n")

	c.App.Println("=== Sesssions ===\n")
	err = ListSessions(c)
	if err != nil {
		return err
	}

	return nil
}
