package startup

import (
	"fmt"

	"github.com/bingoc2/bingoc2-framework/teamserver/version"
)

var sBanner = `
██████╗ ██╗███╗   ██╗ ██████╗  ██████╗  ██████╗██████╗ 
██╔══██╗██║████╗  ██║██╔════╝ ██╔═══██╗██╔════╝╚════██╗
██████╔╝██║██╔██╗ ██║██║  ███╗██║   ██║██║      █████╔╝
██╔══██╗██║██║╚██╗██║██║   ██║██║   ██║██║     ██╔═══╝ 
██████╔╝██║██║ ╚████║╚██████╔╝╚██████╔╝╚██████╗███████╗
╚═════╝ ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝  ╚═════╝╚══════╝`

func Banner() {
	fmt.Println(sBanner)
	fmt.Println(version.SLOGAN)
	fmt.Println("Version:", version.VERSION)
}
