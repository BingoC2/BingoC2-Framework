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

	// parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// run startup script
	startup.Startup()

	// start beacon file server
	go startup.StartFileServer(*lhost)
	go startup.StartLootServer(*lhost)

	// print startup banner
	startup.Banner()

	fmt.Printf("Starting teamserver on %s:%s\n", *lhost, *lport)

	// start teamserver
	cfg := &readline.Config{}
	readline.ListenRemote("tcp", fmt.Sprintf("%s:%s", *lhost, *lport), cfg, app.Handler)
}
