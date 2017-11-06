package configcommands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddSpaceToConfigurationCommand struct {
	BaseConfigCommand
	OrgName             string `long:"org" env:"ORG" description:"Org name to add"`
	SpaceName           string `long:"space" env:"space" description:"Space name to add"`
	SpaceDeveloperGroup string `long:"space-dev-grp" env:"SPACE_DEV_GRP" description:"LDAP group for Space Developer"`
	SpaceMgrGroup       string `long:"space-mgr-grp" env:"SPACE_MGR_GRP" description:"LDAP group for Space Manager"`
	SpaceAuditorGroup   string `long:"space-auditor-grp" env:"SPACE_AUDITOR_GRP" description:"LDAP group for Space Auditor"`
}

//Execute - adds a named space to the configuration
func (c *AddSpaceToConfigurationCommand) Execute([]string) error {
	if c.SpaceName == "" || c.OrgName == "" {
		return fmt.Errorf("Must provide org name and space name")
	}
	spaceConfig := &config.SpaceConfig{
		Org:            c.OrgName,
		Space:          c.SpaceName,
		DeveloperGroup: c.SpaceDeveloperGroup,
		ManagerGroup:   c.SpaceMgrGroup,
		AuditorGroup:   c.SpaceAuditorGroup,
	}

	return config.NewManager(c.ConfigDirectory).AddSpaceToConfig(spaceConfig)
}
