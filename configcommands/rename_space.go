package configcommands

import (
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type RenameSpaceConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName      string `long:"org" description:"Org name" required:"true"`
	SpaceName    string `long:"space" description:"Space name" required:"true"`
	NewSpaceName string `long:"new-space" description:"Space name to rename to" required:"true"`
}

// Execute - renames space config
func (c *RenameSpaceConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	spaceConfig, err := c.ConfigManager.GetSpaceConfig(c.OrgName, c.SpaceName)
	if err != nil {
		return err
	}

	spaces, err := c.ConfigManager.OrgSpaces(c.OrgName)
	if err != nil {
		return err
	}
	spaceConfig.Space = c.NewSpaceName
	spaceConfig.OriginalSpace = c.SpaceName

	err = c.ConfigManager.RenameSpaceConfig(spaceConfig)
	if err != nil {
		return err
	}
	spaces.Replace(c.SpaceName, c.NewSpaceName)
	err = c.ConfigManager.SaveOrgSpaces(spaces)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been renamed to [%s/%s]", c.OrgName, c.SpaceName, c.OrgName, c.NewSpaceName))
	return nil
}

func (c *RenameSpaceConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
