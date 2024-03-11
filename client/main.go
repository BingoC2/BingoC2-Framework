package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/desertbit/readline"
)

func main() {
	// create a new parser
	parser := argparse.NewParser("bingo", "lightweight command & control framework")

	// add arguments
	var rhost *string = parser.String("r", "rhost", &argparse.Options{Required: true, Help: "address to connect to"})
	var rport *string = parser.String("p", "rport", &argparse.Options{Required: true, Help: "port to connect to"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Connecting to teamserver...")
	if err := readline.DialRemote("tcp", fmt.Sprintf("%s:%s", *rhost, *rport)); err != nil {
		fmt.Println(fmt.Errorf("an error occurred: %s", err.Error()))
	}
}
