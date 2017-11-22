package commands

import (
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type InitConfigurationCommand struct {
	BaseConfigCommand
}

//Execute - initializes cf-mgmt configuration
func (c *InitConfigurationCommand) Execute([]string) error {
	lo.G.Warning("This command has been deprecated use lastest cf-mgmt-config cli")
	lo.G.Infof("Initializing config in directory %s", c.ConfigDirectory)
	configManager := config.NewManager(c.ConfigDirectory)
	return configManager.CreateConfigIfNotExists("ldap")
}
