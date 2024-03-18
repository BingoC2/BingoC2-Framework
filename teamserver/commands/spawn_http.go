package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	bingo_errors "github.com/bingoc2/bingoc2-framework/teamserver/errors"
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/hellsgopher"
	"github.com/bingoc2/bingoc2-framework/teamserver/logging"
	"github.com/bingoc2/bingoc2-framework/teamserver/management"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"gopkg.in/yaml.v3"
)

func SpawnHTTP(c *grumble.Context) error {
	// assign vars to flags
	path := c.Flags.String("path")
	listener := c.Flags.String("listener")
	rhost := c.Flags.String("rhost")
	port := c.Flags.String("port")
	sleep := c.Flags.Int("sleep")
	jitter := c.Flags.Int("jitter")
	opsys := c.Flags.String("os")
	arch := c.Flags.String("arch")
	useragent := c.Flags.String("user-agent")
	useragent = strings.ReplaceAll(useragent, " ", "....")
	useragent = strings.ReplaceAll(useragent, ")", "////")
	useragent = strings.ReplaceAll(useragent, "(", "****")
	useragent = strings.ReplaceAll(useragent, ",", "----")
	useragent = strings.ReplaceAll(useragent, ";", "++++")

	// read listener data
	file := version.LISTENER_DATA_DIR + listener + ".yaml"
	yFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var data yamlstructs.ListenerDataYaml
	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		return err
	}

	var RHOST string
	if rhost == "" {
		if data.LHOST == "0.0.0.0" {
			return bingo_errors.ErrInvalidRHOST
		}

		RHOST = data.LHOST
	} else {
		RHOST = rhost
	}

	// check for port
	var RPORT string
	if port == "" {
		RPORT = data.LPORT
	} else {
		RPORT = port
	}

	URI := data.URI

	// generate name
	name := management.GenerateName()

	// error handling
	// check valid listener
	if !hellsgopher.DoesFileExist(fmt.Sprintf(version.LISTENER_DATA_DIR+"%s.yaml", listener)) {
		return bingo_errors.ErrInvalidListener
	}

	// check valid os
	if opsys != "windows" && opsys != "linux" {
		return bingo_errors.ErrInvalidOS
	}

	// check valid arch
	if arch != "amd64" && arch != "i386" {
		return bingo_errors.ErrInvalidArch
	}

	// generate key
	key, err := hellsgopher.GenerateKey()
	if err != nil {
		return err
	}

	// write data to file
	newBeaconData := yamlstructs.BeaconDataYaml{
		Name:      name,
		Key:       key,
		Sleep:     sleep,
		Jitter:    jitter,
		Listener:  listener,
		RHOST:     RHOST,
		UserAgent: useragent,
		GOOS:      opsys,
		GOARCH:    arch,
	}

	yamlData, err := yaml.Marshal(&newBeaconData)
	if err != nil {
		return err
	}
	os.WriteFile(fmt.Sprintf(version.BEACON_DATA_DIR+"%s.yaml", name), yamlData, 0666)

	spawnHTTP(path, listener, RHOST, sleep, jitter, opsys, arch, useragent, key, name, c, RPORT, URI)

	return nil
}

func spawnHTTP(path string, listener string, rhost string, sleep int, jitter int, opsys string, arch string, useragent string, key []byte, name string, c *grumble.Context, rport string, uri string) {
	// build beacon
	if opsys == "windows" {
		name += ".exe"
	}

	keyA := key[:8]
	keyB := key[8:16]
	keyC := key[16:24]
	keyD := key[24:32]

	nKeyA := binary.BigEndian.Uint64(keyA)
	nKeyB := binary.BigEndian.Uint64(keyB)
	nKeyC := binary.BigEndian.Uint64(keyC)
	nKeyD := binary.BigEndian.Uint64(keyD)

	sKeyA := strconv.FormatUint(nKeyA, 10)
	sKeyB := strconv.FormatUint(nKeyB, 10)
	sKeyC := strconv.FormatUint(nKeyC, 10)
	sKeyD := strconv.FormatUint(nKeyD, 10)

	hellsgopher.CmdReturn(fmt.Sprintf("make -C ./agent/http/%s/ BEACON_NAME=%s KEYA=\"%s\" KEYB=\"%s\" KEYC=\"%s\" KEYD=\"%s\" sSLEEP=%d sJITTER=%d RHOST=%s RPORT=%s LISTENER_NAME=%s USERAGENT=%s URI=%s BINARY_NAME=%s cGOOS=%s cGOARCH=%s SPAWN_PATH=%s", opsys, name, sKeyA, sKeyB, sKeyC, sKeyD, sleep, jitter, rhost, rport, listener, useragent, uri, name, opsys, arch, path))

	logging.Okay(fmt.Sprintf("Successfully generated beacon (%s)", name), c)
}
