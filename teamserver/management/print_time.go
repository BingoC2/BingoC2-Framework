package management

import (
	"fmt"
	"time"
)

func PrintTimePretty(t time.Duration) string {
	seconds := t.Seconds()

	// get the number of hours
	hours := int(seconds / 60 / 60)

	// subtract the number of hours from seconds
	seconds = seconds - float64(hours)

	// get minutes
	minutes := int(seconds / 60)

	// subtract the number of minutes from seconds
	seconds = seconds - float64(minutes)

	return fmt.Sprintf("%dh:%dm:%ds", hours, minutes, int(seconds))
}
