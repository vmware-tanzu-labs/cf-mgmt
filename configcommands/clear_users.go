package configcommands

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type ClearUsersCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
}

// Execute - updates org configuration`
func (c *ClearUsersCommand) Execute(args []string) error {
	c.initConfig()
	orgs, err := c.ConfigManager.GetOrgConfigs()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		org.Auditor.Users = []string{}
		org.Manager.Users = []string{}
		org.BillingManager.Users = []string{}

		org.Auditor.LDAPUsers = []string{}
		org.Manager.LDAPUsers = []string{}
		org.BillingManager.LDAPUsers = []string{}

		org.Auditor.SamlUsers = []string{}
		org.Manager.SamlUsers = []string{}
		org.BillingManager.SamlUsers = []string{}

		org.Auditor.LDAPGroups = []string{}
		org.Manager.LDAPGroups = []string{}
		org.BillingManager.LDAPGroups = []string{}

		org.Auditor.LDAPGroup = ""
		org.Manager.LDAPGroup = ""
		org.BillingManager.LDAPGroup = ""

		if err := c.ConfigManager.SaveOrgConfig(&org); err != nil {
			return err
		}
	}
	spaces, err := c.ConfigManager.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, space := range spaces {
		space.Auditor.Users = []string{}
		space.Manager.Users = []string{}
		space.Developer.Users = []string{}
		space.Supporter.Users = []string{}

		space.Auditor.LDAPUsers = []string{}
		space.Manager.LDAPUsers = []string{}
		space.Developer.LDAPUsers = []string{}
		space.Supporter.LDAPUsers = []string{}

		space.Auditor.SamlUsers = []string{}
		space.Manager.SamlUsers = []string{}
		space.Developer.SamlUsers = []string{}
		space.Supporter.SamlUsers = []string{}

		space.Auditor.LDAPGroups = []string{}
		space.Manager.LDAPGroups = []string{}
		space.Developer.LDAPGroups = []string{}
		space.Supporter.LDAPGroups = []string{}

		space.Auditor.LDAPGroup = ""
		space.Manager.LDAPGroup = ""
		space.Developer.LDAPGroup = ""
		space.Supporter.LDAPGroup = ""

		if err := c.ConfigManager.SaveSpaceConfig(&space); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClearUsersCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
