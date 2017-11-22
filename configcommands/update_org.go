package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateOrgConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                          string   `long:"org" description:"Org name" required:"true"`
	PrivateDomains                   []string `long:"private-domain" description:"Private Domain(s) to add, specify multiple times"`
	PrivateDomainsToRemove           []string `long:"private-domain-to-remove" description:"Private Domain(s) to remove, specify multiple times"`
	EnableRemovePrivateDomains       string   `long:"enable-remove-private-domains" description:"Enable removing private domains" choice:"true" choice:"false"`
	SharedPrivateDomains             []string `long:"shared-private-domain" description:"Shared Private Domain(s) to add, specify multiple times"`
	SharedPrivateDomainsToRemove     []string `long:"shared-private-domain-to-remove" description:"Shared Private Domain(s) to remove, specify multiple times"`
	EnableRemoveSharedPrivateDomains string   `long:"enable-remove-shared-private-domains" description:"Enable removing shared private domains" choice:"true" choice:"false"`
	DefaultIsolationSegment          string   `long:"default-isolation-segment" description:"Default isolation segment for org" `
	ClearDefaultIsolationSegment     bool     `long:"clear-default-isolation-segment" description:"Sets the default isolation segment to blank"`
	EnableRemoveUsers                string   `long:"enable-remove-users" description:"Enable removing users from the org" choice:"true" choice:"false"`
	Quota                            OrgQuota `group:"quota"`
	BillingManager                   UserRole `group:"billing-manager" namespace:"billing-manager"`
	Manager                          UserRole `group:"manager" namespace:"manager"`
	Auditor                          UserRole `group:"auditor" namespace:"auditor"`
}

//Execute - updates org configuration`
func (c *UpdateOrgConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	orgConfig, err := c.ConfigManager.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}
	errorString := ""

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
	c.updateUsers(orgConfig, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveOrgConfig(orgConfig); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org [%s] has been updated", c.OrgName))
	return nil
}

func (c *UpdateOrgConfigurationCommand) updateUsers(orgConfig *config.OrgConfig, errorString *string) {
	updateUsersBasedOnRole(&orgConfig.BillingManager, orgConfig.GetBillingManagerGroups(), &c.BillingManager, errorString)
	updateUsersBasedOnRole(&orgConfig.Auditor, orgConfig.GetAuditorGroups(), &c.Auditor, errorString)
	updateUsersBasedOnRole(&orgConfig.Manager, orgConfig.GetManagerGroups(), &c.Manager, errorString)

	orgConfig.BillingManagerGroup = ""
	orgConfig.ManagerGroup = ""
	orgConfig.AuditorGroup = ""
}

func (c *UpdateOrgConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
