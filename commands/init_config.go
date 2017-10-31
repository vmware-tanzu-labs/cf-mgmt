package commands

import (
	"github.com/pivotalservices/cf-mgmt/config"
)

type InitConfigurationCommand struct {
	BaseConfigCommand
}

//Execute - initializes cf-mgmt configuration
func (c *InitConfigurationCommand) Execute([]string) error {
	configManager := config.NewManager(c.ConfigDirectory)
	return configManager.CreateConfigIfNotExists("ldap")
}
