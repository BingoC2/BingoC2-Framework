package management

import (
	"fmt"
	"time"
)

func PrintTimePretty(t time.Duration) string {
	return fmt.Sprintf("%dh%dm%ds", int(t.Hours()), int(t.Minutes()), int(t.Seconds()))
}
