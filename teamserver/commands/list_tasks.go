package commands

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/logging"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"gopkg.in/yaml.v3"
)

func ListTasks(c *grumble.Context) error {
	sessionAgentIds := GetNamesFromDir(version.SESSION_DATA_DIR)

	for _, agentID := range sessionAgentIds {
		err := listTask(agentID, c)
		if err != nil {
			return err
		}
		c.App.Println("")
	}

	return nil
}

func listTask(agentid string, c *grumble.Context) error {
	sessionData, err := yamlstructs.ReadSessionYaml(agentid)
	if err != nil {
		return err
	}

	var keys []string
	for key := range sessionData.Tasks {
		keys = append(keys, strconv.Itoa(key))
	}

	sessionData.Tasks = InsertIntoSlice(sessionData.Tasks, 0, " Task ")
	sessionData.Tasks = InsertIntoSlice(sessionData.Tasks, 1, "======")

	keys = InsertIntoSlice(keys, 0, " Task ID ")
	keys = InsertIntoSlice(keys, 1, "=========")

	c.App.Println(fmt.Sprintf("=== Tasks (%s) ===\n", agentid))
	outputBuffer := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(outputBuffer, 0, 0, 1, ' ', 0)
	for i := range sessionData.Tasks {
		fmt.Fprintln(writer, keys[i]+"\t"+sessionData.Tasks[i])
	}
	writer.Flush()
	c.App.Println(outputBuffer)

	return nil
}

func deleteTask(agentid string, c *grumble.Context) error {
	taskid := c.Args.Int("taskid")

	sessionData, err := yamlstructs.ReadSessionYaml(agentid)
	if err != nil {
		return err
	}

	sessionData.Tasks = remove(sessionData.Tasks, taskid)

	// write new data
	newSessionsData := yamlstructs.SessionDataYaml{
		Name:            sessionData.Name,
		BeaconName:      sessionData.BeaconName,
		AgentID:         agentid,
		Hostname:        sessionData.Hostname,
		IP:              sessionData.IP,
		Interfaces:      sessionData.Interfaces,
		PWD:             sessionData.PWD,
		ProcessPath:     sessionData.ProcessPath,
		ProcessName:     sessionData.ProcessName,
		ProcessID:       sessionData.ProcessID,
		ParentProcessID: sessionData.ParentProcessID,
		Username:        sessionData.Username,
		UID:             sessionData.UID,
		GID:             sessionData.GID,
		OperatingSystem: sessionData.OperatingSystem,
		OSVersion:       sessionData.OSVersion,
		OSBuild:         sessionData.OSBuild,
		Sleep:           sessionData.Sleep,
		Jitter:          sessionData.Jitter,
		Listener:        sessionData.Listener,
		FirstContact:    sessionData.FirstContact,
		LastCallBack:    sessionData.LastCallBack,
		Tasks:           sessionData.Tasks,
	}

	yamlData, err := yaml.Marshal(&newSessionsData)
	if err != nil {
		return err
	}

	os.WriteFile(fmt.Sprintf(version.SESSION_DATA_DIR+"%s.yaml", agentid), yamlData, 0666)

	logging.Okay("successfully removed task", c)

	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
