package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/bingoc2/bingoc2-framework/teamserver/app"
	"github.com/bingoc2/bingoc2-framework/teamserver/startup"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	"github.com/desertbit/readline"
)

func main() {
	// create a new parser
	parser := argparse.NewParser("bingo", version.SLOGAN)

	// add arguments
	var lhost *string = parser.String("l", "lhost", &argparse.Options{Required: false, Help: "address to listen on", Default: "0.0.0.0"})
	var lport *string = parser.String("p", "lport", &argparse.Options{Required: false, Help: "port to listen on", Default: "4455"})
	//var flport *string = parser.String("f", "flport", &argparse.Options{Required: false, Help: "port for fileserver to listen on", Default: "4456"})
	//var binDir *string = parser.String("d", "bindir", &argparse.Options{Required: true, Help: "default path to generate agents to"})

	// parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// run startup script
	startup.Startup()

	// start beacon file server
	//go startup.StartFileServer(*lhost, *flport, *binDir)

	// print startup banner
	startup.Banner()

	fmt.Printf("Starting teamserver on %s:%s\n", *lhost, *lport)

	// start teamserver
	cfg := &readline.Config{}
	readline.ListenRemote("tcp", fmt.Sprintf("%s:%s", *lhost, *lport), cfg, app.Handler)
}
