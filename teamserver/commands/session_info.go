package commands

import (
	"fmt"
	"time"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
)

func SessionInfo(agentid string, c *grumble.Context) error {
	// get session data
	sessionData, err := yamlstructs.ReadSessionYaml(agentid)
	if err != nil {
		return err
	}

	c.App.Println("=== Session Info ===")
	c.App.Println("Agent ID :", sessionData.AgentID)
	c.App.Println("Sleep :", sessionData.Sleep)
	c.App.Println("Jitter :", sessionData.Jitter)
	c.App.Println(fmt.Sprintf("First Contact : %s (%s seconds)", fmt.Sprint(sessionData.FirstContact.Format(time.ANSIC)), fmt.Sprint(time.Since(sessionData.FirstContact).Seconds())))
	c.App.Println(fmt.Sprintf("Last Callback : %s (%s seconds)", fmt.Sprint(sessionData.LastCallBack.Format(time.ANSIC)), fmt.Sprint(time.Since(sessionData.LastCallBack).Seconds())))

	c.App.Println("\n=== System Info ===")
	c.App.Println("Hostname :", sessionData.Hostname)
	c.App.Println(fmt.Sprintf("OS : %s %s", sessionData.OperatingSystem, sessionData.OSVersion))
	c.App.Println("OS Build :", sessionData.OSBuild)

	c.App.Println("\n=== Process Info ===")
	c.App.Println("Process Name :", sessionData.ProcessName)
	c.App.Println("Process Path :", sessionData.ProcessPath)
	c.App.Println("PID :", sessionData.ProcessID)
	c.App.Println("PPID :", sessionData.ParentProcessID)

	c.App.Println("\n=== User Info ===")
	c.App.Println("Username :", sessionData.Username)
	c.App.Println("UID :", sessionData.UID)
	c.App.Println("GID :", sessionData.GID)

	return nil
}
