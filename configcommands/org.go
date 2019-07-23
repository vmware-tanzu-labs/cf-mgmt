package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type OrgConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                          string        `long:"org" description:"Org name" required:"true"`
	PrivateDomains                   []string      `long:"private-domain" description:"Private Domain(s) to add, specify multiple times"`
	PrivateDomainsToRemove           []string      `long:"private-domain-to-remove" description:"Private Domain(s) to remove, specify multiple times"`
	EnableRemovePrivateDomains       string        `long:"enable-remove-private-domains" description:"Enable removing private domains" choice:"true" choice:"false"`
	SharedPrivateDomains             []string      `long:"shared-private-domain" description:"Shared Private Domain(s) to add, specify multiple times"`
	SharedPrivateDomainsToRemove     []string      `long:"shared-private-domain-to-remove" description:"Shared Private Domain(s) to remove, specify multiple times"`
	EnableRemoveSharedPrivateDomains string        `long:"enable-remove-shared-private-domains" description:"Enable removing shared private domains" choice:"true" choice:"false"`
	EnableRemoveSpaces               string        `long:"enable-remove-spaces" description:"Enable removing spaces" choice:"true" choice:"false"`
	DefaultIsolationSegment          string        `long:"default-isolation-segment" description:"Default isolation segment for org" `
	ClearDefaultIsolationSegment     bool          `long:"clear-default-isolation-segment" description:"Sets the default isolation segment to blank"`
	EnableRemoveUsers                string        `long:"enable-remove-users" description:"Enable removing users from the org" choice:"true" choice:"false"`
	NamedQuota                       string        `long:"named-quota" description:"Named quota to assign to org"`
	ClearNamedQuota                  bool          `long:"clear-named-quota" description:"Sets the named quota to blank"`
	Quota                            OrgQuota      `group:"quota"`
	BillingManager                   UserRole      `group:"billing-manager" namespace:"billing-manager"`
	Manager                          UserRole      `group:"manager" namespace:"manager"`
	Auditor                          UserRole      `group:"auditor" namespace:"auditor"`
	ServiceAccess                    ServiceAccess `group:"service-access"`
}

//Execute - updates org configuration`
func (c *OrgConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	var orgConfig *config.OrgConfig
	var err error
	var orgSpaces *config.Spaces
	var newOrg bool
	orgConfig, err = c.ConfigManager.GetOrgConfig(c.OrgName)
	if err != nil {
		newOrg = true
		lo.G.Debugf("Org [%s] doesn't exist creating it", c.OrgName)
		orgConfig = &config.OrgConfig{
			Org:                        c.OrgName,
			RemoveUsers:                true,
			RemovePrivateDomains:       true,
			RemoveSharedPrivateDomains: true,
		}

	} else {
		newOrg = false
	}

	if c.Quota.EnableOrgQuota == "true" && c.NamedQuota != "" {
		return fmt.Errorf("cannot enable org quota and use named quotas")
	}

	orgSpaces, err = c.ConfigManager.OrgSpaces(c.OrgName)
	if err != nil {
		orgSpaces = &config.Spaces{
			Org:                c.OrgName,
			EnableDeleteSpaces: true,
		}
	}
	errorString := ""

	convertToBool("enable-remove-spaces", &orgSpaces.EnableDeleteSpaces, c.EnableRemoveSpaces, &errorString)
	if c.DefaultIsolationSegment != "" {
		orgConfig.DefaultIsoSegment = c.DefaultIsolationSegment
	}
	if c.ClearDefaultIsolationSegment {
		orgConfig.DefaultIsoSegment = ""
	}
	convertToBool("enable-remove-users", &orgConfig.RemoveUsers, c.EnableRemoveUsers, &errorString)
	orgConfig.PrivateDomains = removeFromSlice(addToSlice(orgConfig.PrivateDomains, c.PrivateDomains, &errorString), c.PrivateDomainsToRemove)
	convertToBool("enable-remove-private-domains", &orgConfig.RemovePrivateDomains, c.EnableRemovePrivateDomains, &errorString)

	orgConfig.SharedPrivateDomains = removeFromSlice(addToSlice(orgConfig.SharedPrivateDomains, c.SharedPrivateDomains, &errorString), c.SharedPrivateDomainsToRemove)
	convertToBool("enable-remove-shared-private-domains", &orgConfig.RemoveSharedPrivateDomains, c.EnableRemoveSharedPrivateDomains, &errorString)

	updateOrgQuotaConfig(orgConfig, c.Quota, &errorString)
	if c.NamedQuota != "" {
		orgConfig.NamedQuota = c.NamedQuota
	}
	if c.ClearNamedQuota {
		orgConfig.NamedQuota = ""
	}
	c.updateUsers(orgConfig, &errorString)

	if c.ServiceAccess.ServiceNameToRemove != "" {
		delete(orgConfig.ServiceAccess, c.ServiceAccess.ServiceNameToRemove)
	}

	if c.ServiceAccess.ServiceName != "" {
		if len(c.ServiceAccess.Plans) > 0 {
			orgConfig.ServiceAccess[c.ServiceAccess.ServiceName] = c.ServiceAccess.Plans
		} else {
			orgConfig.ServiceAccess[c.ServiceAccess.ServiceName] = []string{"*"}
		}
	}

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveOrgConfig(orgConfig); err != nil {
		return err
	}

	if err := c.ConfigManager.SaveOrgSpaces(orgSpaces); err != nil {
		return err
	}

	if newOrg {
		fmt.Println(fmt.Sprintf("The org [%s] has been created", c.OrgName))
	} else {
		fmt.Println(fmt.Sprintf("The org [%s] has been updated", c.OrgName))
	}
	return nil
}

func (c *OrgConfigurationCommand) updateUsers(orgConfig *config.OrgConfig, errorString *string) {
	updateUsersBasedOnRole(&orgConfig.BillingManager, orgConfig.GetBillingManagerGroups(), &c.BillingManager, errorString)
	updateUsersBasedOnRole(&orgConfig.Auditor, orgConfig.GetAuditorGroups(), &c.Auditor, errorString)
	updateUsersBasedOnRole(&orgConfig.Manager, orgConfig.GetManagerGroups(), &c.Manager, errorString)

	orgConfig.BillingManagerGroup = ""
	orgConfig.ManagerGroup = ""
	orgConfig.AuditorGroup = ""
}

func (c *OrgConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
