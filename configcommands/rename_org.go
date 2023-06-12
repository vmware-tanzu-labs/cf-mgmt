package configcommands

import (
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type RenameOrgConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName    string `long:"org" description:"Org name" required:"true"`
	NewOrgName string `long:"new-org" description:"Org name to rename to" required:"true"`
}

// Execute - updates org configuration`
func (c *RenameOrgConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	orgs, err := c.ConfigManager.Orgs()
	if err != nil {
		return err
	}
	orgConfig, err := c.ConfigManager.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}

	orgSpaces, err := c.ConfigManager.OrgSpaces(c.OrgName)
	if err != nil {
		return err
	}

	orgConfig.Org = c.NewOrgName
	orgConfig.OriginalOrg = c.OrgName
	if err := c.ConfigManager.RenameOrgConfig(orgConfig); err != nil {
		return err
	}

	orgs.Replace(orgConfig.OriginalOrg, orgConfig.Org)

	if err := c.ConfigManager.SaveOrgs(orgs); err != nil {
		return err
	}

	orgSpaces.Org = c.NewOrgName
	if err := c.ConfigManager.SaveOrgSpaces(orgSpaces); err != nil {
		return err
	}

	for _, spaceName := range orgSpaces.Spaces {
		spaceConfig, err := c.ConfigManager.GetSpaceConfig(orgConfig.Org, spaceName)
		if err != nil {
			return err
		}
		spaceConfig.Org = orgConfig.Org
		err = c.ConfigManager.SaveSpaceConfig(spaceConfig)
		if err != nil {
			return err
		}
	}
	fmt.Println(fmt.Sprintf("The org [%s] has been renamed to [%s]", c.OrgName, c.NewOrgName))
	return nil
}

func (c *RenameOrgConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
