package configcommands

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type InitConfigurationCommand struct {
	BaseConfigCommand
}

// Execute - initializes cf-mgmt configuration
func (c *InitConfigurationCommand) Execute([]string) error {
	lo.G.Infof("Initializing config in directory %s", c.ConfigDirectory)
	configManager := config.NewManager(c.ConfigDirectory)
	return configManager.CreateConfigIfNotExists("ldap")
}
