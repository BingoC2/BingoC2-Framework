package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
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
