package management

import (
	"fmt"
	"time"
)

func PrintTimePretty(t time.Duration) string {
	hours := t.Hours()
	t = t - time.Duration(hours*time.Hour.Hours())

	minutes := t.Minutes()
	t = t - time.Duration(minutes*time.Hour.Minutes())

	seconds := t.Seconds()

	return fmt.Sprintf("%dh%dm%ds", int(hours), int(minutes), int(seconds))
}
