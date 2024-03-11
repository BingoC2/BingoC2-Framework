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
	ProcessUser     string   `json:"processuser"`
	OperatingSystem string   `json:"os"`
	Sleep           int      `json:"sleep"`
	Jitter          int      `json:"json"`
	Listener        string   `json:"listener"`
}
