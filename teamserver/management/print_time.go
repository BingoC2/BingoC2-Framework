package management

import (
	"fmt"
	"time"
)

func PrintTimePretty(t time.Duration) string {
	seconds := t.Seconds()
	fmt.Println(seconds)

	// get the number of hours
	hours := int(seconds / 60 / 60)
	fmt.Println(hours)

	// subtract the number of hours from seconds
	seconds = seconds - float64(hours)
	fmt.Println(seconds)

	// get minutes
	minutes := int(seconds / 60)
	fmt.Println(minutes)

	// subtract the number of minutes from seconds
	seconds = seconds - float64(minutes)
	fmt.Println(float64(minutes))
	fmt.Println(seconds)

	return fmt.Sprintf("%dh:%dm:%ds", hours, minutes, int(seconds))
}
