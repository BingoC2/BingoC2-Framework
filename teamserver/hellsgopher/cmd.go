package hellsgopher

import (
	"os/exec"
	"runtime"
)

// will run a command with either `bash -c` or `cmd /C` and will provide no output
func CmdNoOut(command string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("bash", "-c", command)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		return
	}

	cmd.CombinedOutput()
}
