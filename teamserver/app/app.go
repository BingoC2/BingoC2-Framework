package app

import (
	"github.com/bingoc2/bingoc2-framework/teamserver/commands"
	bingo_errors "github.com/bingoc2/bingoc2-framework/teamserver/errors"
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	"github.com/desertbit/readline"
	"github.com/fatih/color"
)

var (
	Blue = color.New(color.FgBlue)
	Red  = color.New(color.FgRed)
)

func Handler(rl *readline.Instance) {
	var app = grumble.New(&grumble.Config{
		Name:        "bingoc2",
		Description: "lightweight \"stage zero\" command & control adversary emulation framework",
		InterruptHandler: func(a *grumble.App, count int) {

		},
		Prompt:                "[client] bingoc2> ",
		PromptColor:           Blue,
		ErrorColor:            Red,
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,
		HelpHeadlineColor:     Blue,
		ASCIILogoColor:        Blue,
	})

	app.SetPrintASCIILogo(func(a *grumble.App) {
		app.Println(`
██████╗ ██╗███╗   ██╗ ██████╗  ██████╗  ██████╗██████╗ 
██╔══██╗██║████╗  ██║██╔════╝ ██╔═══██╗██╔════╝╚════██╗
██████╔╝██║██╔██╗ ██║██║  ███╗██║   ██║██║      █████╔╝
██╔══██╗██║██║╚██╗██║██║   ██║██║   ██║██║     ██╔═══╝ 
██████╔╝██║██║ ╚████║╚██████╔╝╚██████╔╝╚██████╗███████╗
╚═════╝ ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝  ╚═════╝╚══════╝`)
		app.Println("Version:", version.VERSION)
	})

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
			return commands.ListenHttp(c)
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
			return commands.Deafen(c)
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
			f.String("p", "path", "../../../bin/", "full path to spawn the beacon to")
			f.String("l", "listener", "", "name of the listener")
			f.String("r", "rhost", "", "host to callback to; if left empty, will default to host specified in listener (does not work if listener is 0.0.0.0) (ex: domain.com)")
			f.Int("s", "sleep", 5, "interval between callbacks")
			f.Int("j", "jitter", 2, "range of random intervals for callback (ex: sleep is 5 and jitter is 2, beacon will callback between 3 and 7 seconds)")
			f.String("o", "os", "windows", "operating system to compile beacon for (supports: windows, linux)")
			f.String("a", "arch", "amd64", "architecture to compile beacon for (amd64, i386)")
			f.String("u", "user-agent", "bingoc2/1.0", "user agent for agents to use")
		},
		Run: func(c *grumble.Context) error {
			return commands.SpawnHTTP(c)
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
		Flags: func(f *grumble.Flags) {
			f.String("n", "name", "", "name of beacon to respawn")
		},
		Run: func(c *grumble.Context) error {
			return nil
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
			return commands.ListBeacons(c)
		},
	})

	listCMD.AddCommand(&grumble.Command{
		Name:    "listeners",
		Help:    "list listeners",
		Aliases: []string{"l"},
		Run: func(c *grumble.Context) error {
			return commands.ListListeners(c)
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
			return commands.ListSessions(c)
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "task",
		Help:    "send a single task to a session",
		Aliases: []string{"t"},
		Flags: func(f *grumble.Flags) {
			f.String("i", "agent-id", "", "agent to task")
			f.String("t", "task", "", "task to give to agent (full list available by using `task -l`)")
			f.Bool("l", "list", false, "list all available tasks")
		},
		Run: func(c *grumble.Context) error {
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "details",
		Help: "get details of a session",
		Flags: func(f *grumble.Flags) {
			f.String("i", "agent-id", "", "agent id of session to get details of")
		},
		Run: func(c *grumble.Context) error {
			return nil
		},
	})

	app.RunWithReadline(rl)
}
