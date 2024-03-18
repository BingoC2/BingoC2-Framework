package commands

import (
	bingo_errors "github.com/bingoc2/bingoc2-framework/teamserver/errors"
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
)

func RegisterMainCommands(app *grumble.App) {
	listenCmd := &grumble.Command{
		Name:    "listen",
		Help:    "start a listener",
		Aliases: []string{"l"},
	}
	app.AddCommand(listenCmd)

	listenCmd.AddCommand(&grumble.Command{
		Name: "http",
		Help: "start a http listener",
		Flags: func(f *grumble.Flags) {
			f.String("l", "lhost", "0.0.0.0", "host to listen on")
			f.String("p", "port", "80", "port to listen on")
			f.String("u", "uri", "index.php", "uri for agent to callback to (cannot be in use by any other listener)")
			f.String("n", "name", "blank", "name of the new listener")
		},
		Run: func(c *grumble.Context) error {
			return ListenHttp(c)
		},
	})

	listenCmd.AddCommand(&grumble.Command{
		Name: "dns",
		Help: "start a dns listener",
		Flags: func(f *grumble.Flags) {
			f.String("l", "lhost", "0.0.0.0", "host to listen on")
			f.String("p", "port", "53", "port to listen on")
			f.String("n", "name", "blank", "name of the new listener")
		},
		Run: func(c *grumble.Context) error {
			return bingo_errors.ErrCmdNotSupported
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "deafen",
		Help:    "stop a listener",
		Aliases: []string{"d"},
		Flags: func(f *grumble.Flags) {
			f.String("n", "name", "", "listener to stop")
		},
		Run: func(c *grumble.Context) error {
			return Deafen(c)
		},
	})

	spawnCMD := &grumble.Command{
		Name:    "spawn",
		Help:    "spawn a beacon",
		Aliases: []string{"s"},
	}
	app.AddCommand(spawnCMD)

	spawnCMD.AddCommand(&grumble.Command{
		Name: "http",
		Help: "spawn http beacon",
		Flags: func(f *grumble.Flags) {
			f.String("P", "path", "../../../bin/", "full path to spawn the beacon to")
			f.String("l", "listener", "", "name of the listener")
			f.String("r", "rhost", "", "host to callback to; if left empty, will default to host specified in listener (does not work if listener is 0.0.0.0) (ex: domain.com)")
			f.String("p", "port", "", "port to call back to; will override the listener port")
			f.Int("s", "sleep", 5, "interval between callbacks")
			f.Int("j", "jitter", 2, "range of random intervals for callback (ex: sleep is 5 and jitter is 2, beacon will callback between 3 and 7 seconds)")
			f.String("o", "os", "windows", "operating system to compile beacon for (supports: windows, linux)")
			f.String("a", "arch", "amd64", "architecture to compile beacon for (amd64, i386)")
			f.String("u", "user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36", "user agent for agents to use")
			f.Bool("d", "isDLL", false, "generate the agent as a DLL file")
		},
		Run: func(c *grumble.Context) error {
			return SpawnHTTP(c)
		},
	})

	spawnCMD.AddCommand(&grumble.Command{
		Name: "dns",
		Help: "spawn dns beacon",
		Flags: func(f *grumble.Flags) {
			f.String("n", "name", "blank", "name of the beacon")
			f.String("l", "listener", "", "name of the listener")
			f.String("s", "sleep", "5", "interval between callbacks")
			f.String("j", "jitter", "2", "range of random intervals for callback (ex. sleep is 5 and jitter is 2, beacon will callback between 3 and 7 seconds)")
			f.String("o", "os", "windows", "operating system to compile beacon for (supported: windows, linux)")
			f.String("a", "arch", "x64", "architecture to compile beacon for (supported: x64, x86)")
		},
		Run: func(c *grumble.Context) error {
			return bingo_errors.ErrCmdNotSupported
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "respawn",
		Help:    "respawn a beacon",
		Aliases: []string{"r"},
		Args: func(a *grumble.Args) {
			a.String("name", "name of agent to respawn")
		},
		Run: func(c *grumble.Context) error {
			return Respawn(c)
		},
	})

	listCMD := &grumble.Command{
		Name: "list",
		Help: "list beacons/listeners",
	}
	app.AddCommand(listCMD)

	listCMD.AddCommand(&grumble.Command{
		Name:    "beacons",
		Help:    "list beacons",
		Aliases: []string{"b"},
		Run: func(c *grumble.Context) error {
			return ListBeacons(c)
		},
	})

	listCMD.AddCommand(&grumble.Command{
		Name:    "listeners",
		Help:    "list listeners",
		Aliases: []string{"l"},
		Run: func(c *grumble.Context) error {
			return ListListeners(c)
		},
	})

	listCMD.AddCommand(&grumble.Command{
		Name: "sessions",
		Help: "list sessions",
		Flags: func(f *grumble.Flags) {
			f.Bool("a", "active", false, "only list sessions with active shells")
		},
		Aliases: []string{"s"},
		Run: func(c *grumble.Context) error {
			return ListSessions(c)
		},
	})

	listCMD.AddCommand(&grumble.Command{
		Name:    "all",
		Help:    "list listeners, beacons, and sessions",
		Aliases: []string{"a"},
		Run: func(c *grumble.Context) error {
			return ListAll(c)
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "use",
		Help: "drop into a tasking shell on a session",
		Args: func(a *grumble.Args) {
			a.String("agentid", "id of the the agent")
		},
		Run: func(c *grumble.Context) error {
			return Use(c)
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "tasks",
		Help: "see tasks currently in queue for one/all sessions",
		Args: func(a *grumble.Args) {
			a.String("agentid", "id of agent to view tasks for (leave blank for all)", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			return ListTasks(c)
		},
	})
}
