package hellsgopher

import (
	"errors"
	"io"
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

func CopyFile(sourcePath string, destinationPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
