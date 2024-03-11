package commands

import (
	"fmt"
	"os"

	grumble "github.com/bingoc2/bingoc2-framework/teamserver/grumble_modified"
	"github.com/bingoc2/bingoc2-framework/teamserver/version"
	yamlstructs "github.com/bingoc2/bingoc2-framework/teamserver/yaml_structs"
	"gopkg.in/yaml.v3"
)

func Deafen(c *grumble.Context) error {
	name := c.Flags.String("name")

	listenerData, err := yamlstructs.ReadListenerYaml(name)
	if err != nil {
		return err
	}

	listenerData.IsAlive = false

	yamlData, err := yaml.Marshal(&listenerData)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf(version.LISTENER_DATA_DIR+"%s.yaml", name), yamlData, 0666)
	if err != nil {
		return err
	}

	return nil
}
