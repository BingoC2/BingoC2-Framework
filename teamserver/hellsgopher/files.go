package hellsgopher

import (
	"errors"
	"os"
)

// check if a file exists
// returns true if file exists
func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
