package initialization

type httpPostInitRequest struct {
	Hostname        string   `json:"hostname"`
	IP              string   `json:"ip"`
	Interfaces      []string `json:"interfaces"`
	ProcessPath     string   `json:"processpath"`
	PWD             string   `json:"pwd"`
	ProcessName     string   `json:"processname"`
	ProcessID       int      `json:"pid"`
	ParentProcessID int      `json:"ppid"`
	Username        string   `json:"username"`
	UID             string   `json:"uid"`
	GID             string   `json:"gid"`
	OperatingSystem string   `json:"os"`
	OSVersion       string   `json:"version"`
	OSBuild         string   `json:"build"`
	Sleep           int      `json:"sleep"`
	Jitter          int      `json:"json"`
	Listener        string   `json:"listener"`
}
