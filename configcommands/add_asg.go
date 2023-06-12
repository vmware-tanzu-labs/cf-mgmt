package configcommands

import (
	"errors"
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type AddASGToConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	ASGName  string `long:"asg" description:"ASG name" required:"true"`
	FilePath string `long:"path" description:"path to asg definition"`
	Override bool   `long:"override" description:"override current definition"`
	ASGType  string `long:"type" description:"Space asg or default asg" choice:"space" choice:"default" default:"space"`
}

// Execute - adds a named asg to the configuration
func (c *AddASGToConfigurationCommand) Execute([]string) error {
	lo.G.Warning("*** Deprecated *** - Use `asg` command instead for adding/updating asg configurations")
	c.initConfig()

	errorString := ""

	if errorString != "" {
		return errors.New(errorString)
	}

	if !c.Override {
		if c.ASGType == "space" {
			asgConfigs, err := c.ConfigManager.GetASGConfigs()
			if err != nil {
				return err
			}
			for _, asgConfig := range asgConfigs {
				if c.ASGName == asgConfig.Name {
					return errors.New(fmt.Sprintf("asg [%s] already exists if wanting to update use --override flag", c.ASGName))
				}
			}
		} else {
			asgConfigs, err := c.ConfigManager.GetDefaultASGConfigs()
			if err != nil {
				return err
			}
			for _, asgConfig := range asgConfigs {
				if c.ASGName == asgConfig.Name {
					return errors.New(fmt.Sprintf("asg [%s] already exists if wanting to update use --override flag", c.ASGName))
				}
			}
		}
	}

	securityGroupsBytes := []byte("[]")
	if c.FilePath != "" {
		bytes, err := config.LoadFileBytes(c.FilePath)
		if err != nil {
			return err
		}
		securityGroupsBytes = bytes
	}
	if c.ASGType == "space" {
		if err := config.NewManager(c.ConfigDirectory).AddSecurityGroup(c.ASGName, securityGroupsBytes); err != nil {
			return err
		}
	} else {
		if err := config.NewManager(c.ConfigDirectory).AddDefaultSecurityGroup(c.ASGName, securityGroupsBytes); err != nil {
			return err
		}
	}
	fmt.Println(fmt.Sprintf("The asg [%s] has been updated", c.ASGName))
	return nil
}

func (c *AddASGToConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
