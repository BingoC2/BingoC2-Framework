package yamlstructs

import (
	"os"
	"time"

	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	"gopkg.in/yaml.v3"
)

// structure for HTTP Listener data yaml file
type ListenerDataYaml struct {
	Channel string `yaml:"channel"`
	Name    string `yaml:"name"`
	LHOST   string `yaml:"lhost"`
	LPORT   string `yaml:"lport"`
	URI     string `yaml:"uri"`
	IsAlive bool   `yaml:"isalive"`
}

// return data from a HTTP listener data yaml file in LISTENER_DATA_DIR
func ReadListenerYaml(name string) (ListenerDataYaml, error) {
	var data ListenerDataYaml

	file := version.LISTENER_DATA_DIR + name + ".yaml"

	yFile, err := os.ReadFile(file)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

type SessionDataYaml struct {
	Name            string    `yaml:"name"`
	BeaconName      string    `yaml:"beaconName"`
	AgentID         string    `yaml:"agentid"`
	Hostname        string    `yaml:"hostname"`
	IP              string    `yaml:"ip"`
	Interfaces      []string  `yaml:"interfaces"`
	PWD             string    `yaml:"pwd"`
	ProcessPath     string    `yaml:"processpath"`
	ProcessName     string    `yaml:"processname"`
	ProcessID       int       `yaml:"pid"`
	ParentProcessID int       `yaml:"ppid"`
	Username        string    `yaml:"username"`
	UID             string    `yaml:"uid"`
	GID             string    `yaml:"gid"`
	OperatingSystem string    `yaml:"os"`
	OSVersion       string    `yaml:"version"`
	OSBuild         string    `yaml:"build"`
	Sleep           int       `yaml:"sleep"`
	Jitter          int       `yaml:"jitter"`
	Listener        string    `yaml:"listener"`
	FirstContact    time.Time `yaml:"firstcontact"`
	LastCallBack    time.Time `yaml:"lastcallback"`
	Tasks           []string  `yaml:"tasks"`
}

// return data from a yaml file stored in "./_data/sessions/"
func ReadSessionYaml(name string) (SessionDataYaml, error) {
	var data SessionDataYaml

	file := "./_data/sessions/" + name + ".yaml"
	yFile, err := os.ReadFile(file)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		return SessionDataYaml{}, err
	}

	return data, nil
}

type BeaconDataYaml struct {
	Name      string `yaml:"name"`
	Sleep     int    `yaml:"sleep"`
	Jitter    int    `yaml:"jitter"`
	Listener  string `yaml:"listener"`
	RHOST     string `yaml:"rhost"`
	UserAgent string `yaml:"useragent"`
	GOOS      string `yaml:"goos"`
	GOARCH    string `yaml:"goarch"`
	Key       []byte `yaml:"key"`
}

// return data from a yaml file stored in "./_data/beacons/"
func ReadBeaconYaml(name string) (BeaconDataYaml, error) {
	var data BeaconDataYaml

	file := "./_data/beacons/" + name + ".yaml"
	yFile, err := os.ReadFile(file)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		return BeaconDataYaml{}, err
	}

	return data, nil
}
