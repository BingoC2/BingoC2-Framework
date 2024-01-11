package startup

import "os"

// checks if data directory exists, if it doesn't create data dirs
func Startup() {
	if !checkIfDirExist("../_data") {
		createDir("../_data")
		createDir("../_data/listeners")
		createDir("../_data/sessions")
		createDir("../_data/beacons")
	}
}

// checks if a directory exists. returns true if it exists
func checkIfDirExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// creates a directory with perms 0777
func createDir(path string) {
	os.Mkdir(path, 0777)
}
