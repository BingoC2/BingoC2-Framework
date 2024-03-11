package logging

import (
	"fmt"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/fatih/color"
)

func Okay(s string, c *grumble.Context) {
	green := color.New(color.FgGreen)

	green.Println(fmt.Sprintf("[+] %s", s))
	c.App.Println(fmt.Sprintf("[+] %s", s))
}

func Error(s string, c *grumble.Context) {
	red := color.New(color.FgRed)

	red.Println(fmt.Sprintf("[-] %s", s))
	c.App.Println(fmt.Sprintf("[-] %s", s))
}
