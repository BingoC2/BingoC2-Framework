package commands

import (
	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
)

func Respawn(c *grumble.Context) error {
	agentName := c.Args.String("name")

	agentData, err := yamlstructs.ReadBeaconYaml(agentName)
	if err != nil {
		return err
	}

	listenerData, err := yamlstructs.ReadListenerYaml(agentData.Listener)
	if err != nil {
		return err
	}

	spawnHTTP("../../../bin", agentData.Listener, agentData.RHOST, agentData.Sleep, agentData.Jitter, agentData.GOOS, agentData.GOARCH, agentData.UserAgent, agentData.Key, agentName, c, listenerData.LPORT, listenerData.URI)

	return nil
}
