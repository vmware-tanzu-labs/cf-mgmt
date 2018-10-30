package commands

import (
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type AddOrgToConfigurationCommand struct {
	BaseConfigCommand
	OrgName            string `long:"org" env:"ORG" description:"Org name to add" required:"true"`
	OrgBillingMgrGroup string `long:"org-billing-mgr-grp" env:"ORG_BILLING_MGR_GRP" description:"LDAP group for Org Billing Manager"`
	OrgMgrGroup        string `long:"org-mgr-grp" env:"ORG_MGR_GRP" description:"LDAP group for Org Manager"`
	OrgAuditorGroup    string `long:"org-auditor-grp" env:"ORG_AUDITOR_GRP" description:"LDAP group for Org Auditor"`
}

//Execute - adds a named org to the configuration
func (c *AddOrgToConfigurationCommand) Execute([]string) error {
	lo.G.Warning("This command has been deprecated use lastest cf-mgmt-config cli")
	orgConfig := &config.OrgConfig{
		Org:                  c.OrgName,
		BillingManagerGroup:  c.OrgBillingMgrGroup,
		ManagerGroup:         c.OrgMgrGroup,
		AuditorGroup:         c.OrgAuditorGroup,
		RemoveUsers:          true,
		RemovePrivateDomains: true,
	}
	spacesConfig := &config.Spaces{Org: orgConfig.Org, EnableDeleteSpaces: true}
	return config.NewManager(c.ConfigDirectory).AddOrgToConfig(orgConfig, spacesConfig)
}
