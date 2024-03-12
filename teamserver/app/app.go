package app

import (
	"github.com/bingoc2/bingoc2-framework/teamserver/commands"
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

	commands.RegisterMainCommands(app)

	app.RunWithReadline(rl)
}
