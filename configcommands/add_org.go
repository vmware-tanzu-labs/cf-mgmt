package configcommands

import (
	"github.com/pivotalservices/cf-mgmt/config"
)

type AddOrgToConfigurationCommand struct {
	BaseConfigCommand
	OrgName            string `long:"org" env:"ORG" description:"Org name to add"`
	OrgBillingMgrGroup string `long:"org-billing-mgr-grp" env:"ORG_BILLING_MGR_GRP" description:"LDAP group for Org Billing Manager"`
	OrgMgrGroup        string `long:"org-mgr-grp" env:"ORG_MGR_GRP" description:"LDAP group for Org Manager"`
	OrgAuditorGroup    string `long:"org-auditor-grp" env:"ORG_AUDITOR_GRP" description:"LDAP group for Org Auditor"`
}

//Execute - adds a named org to the configuration
func (c *AddOrgToConfigurationCommand) Execute([]string) error {
	orgConfig := &config.OrgConfig{
		Org:                 c.OrgName,
		BillingManagerGroup: c.OrgBillingMgrGroup,
		ManagerGroup:        c.OrgMgrGroup,
		AuditorGroup:        c.OrgAuditorGroup,
	}
	return config.NewManager(c.ConfigDirectory).AddOrgToConfig(orgConfig)
}
