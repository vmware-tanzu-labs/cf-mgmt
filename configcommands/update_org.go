package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateOrgConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                      string   `long:"org" description:"Org name" required:"true"`
	PrivateDomains               []string `long:"private-domain" description:"Private Domain(s) to add, specify muliple times"`
	PrivateDomainsToRemove       []string `long:"private-domain-to-remove" description:"Private Domain(s) to remove, specify muliple times"`
	EnableRemovePrivateDomains   string   `long:"enable-remove-private-domains" description:"Enable removing private domains" choice:"true" choice:"false"`
	DefaultIsolationSegment      string   `long:"default-isolation-segment" description:"Default isolation segment for org" `
	ClearDefaultIsolationSegment bool     `long:"clear-default-isolation-segment" description:"Sets the default isolation segment to blank"`
	EnableRemoveUsers            string   `long:"enable-remove-users" description:"Enable removing users from the org" choice:"true" choice:"false"`
	Quota                        struct {
		EnableOrgQuota          string `long:"enable-org-quota" description:"Enable the Org Quota in the config" choice:"true" choice:"false"`
		MemoryLimit             string `long:"memory-limit" description:"An Org's memory limit in Megabytes"`
		InstanceMemoryLimit     string `long:"instance-memory-limit" description:"Global Org Application instance memory limit in Megabytes"`
		TotalRoutes             string `long:"total-routes" description:"Total Routes capacity for an Org"`
		TotalServices           string `long:"total-services" description:"Total Services capacity for an Org"`
		PaidServicesAllowed     string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an org" choice:"true" choice:"false"`
		TotalPrivateDomains     string `long:"total-private-domains" description:"Total Private Domain capacity for an Org"`
		TotalReservedRoutePorts string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Org"`
		TotalServiceKeys        string `long:"total-service-keys" description:"Total Service Keys capacity for an Org"`
		AppInstanceLimit        string `long:"app-instance-limit" description:"Total Service Keys capacity for an Org"`
	} `group:"quota"`
	BillingManager UserRole `group:"billing-manager" namespace:"billing-manager"`
	Manager        UserRole `group:"manager" namespace:"manager"`
	Auditor        UserRole `group:"auditor" namespace:"auditor"`
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
	c.updatePrivateDomainConfig(orgConfig, &errorString)
	c.updateQuotaConfig(orgConfig, &errorString)
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
	updateUsersBasedOnRole(&orgConfig.BillingManager, orgConfig.GetBillingManagerGroups(), &c.BillingManager)
	updateUsersBasedOnRole(&orgConfig.Auditor, orgConfig.GetAuditorGroups(), &c.Auditor)
	updateUsersBasedOnRole(&orgConfig.Manager, orgConfig.GetManagerGroups(), &c.Manager)

	orgConfig.BillingManagerGroup = ""
	orgConfig.ManagerGroup = ""
	orgConfig.AuditorGroup = ""
}

func (c *UpdateOrgConfigurationCommand) updatePrivateDomainConfig(orgConfig *config.OrgConfig, errorString *string) {
	orgConfig.PrivateDomains = removeFromSlice(append(orgConfig.PrivateDomains, c.PrivateDomains...), c.PrivateDomainsToRemove)
	convertToBool("enable-remove-private-domains", &orgConfig.RemovePrivateDomains, c.EnableRemovePrivateDomains, errorString)
}

func (c *UpdateOrgConfigurationCommand) updateQuotaConfig(orgConfig *config.OrgConfig, errorString *string) {
	convertToBool("enable-org-quota", &orgConfig.EnableOrgQuota, c.Quota.EnableOrgQuota, errorString)
	convertToInt("memory-limit", &orgConfig.MemoryLimit, c.Quota.MemoryLimit, errorString)
	convertToInt("instance-memory-limit", &orgConfig.InstanceMemoryLimit, c.Quota.InstanceMemoryLimit, errorString)
	convertToInt("total-routes", &orgConfig.TotalRoutes, c.Quota.TotalRoutes, errorString)
	convertToInt("total-services", &orgConfig.TotalServices, c.Quota.TotalServices, errorString)
	convertToBool("paid-service-plans-allowed", &orgConfig.PaidServicePlansAllowed, c.Quota.PaidServicesAllowed, errorString)
	convertToInt("total-private-domains", &orgConfig.TotalPrivateDomains, c.Quota.TotalPrivateDomains, errorString)
	convertToInt("total-reserved-route-ports", &orgConfig.TotalReservedRoutePorts, c.Quota.TotalReservedRoutePorts, errorString)
	convertToInt("total-service-keys", &orgConfig.TotalServiceKeys, c.Quota.TotalServiceKeys, errorString)
	convertToInt("app-instance-limit", &orgConfig.AppInstanceLimit, c.Quota.AppInstanceLimit, errorString)
}

func (c *UpdateOrgConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
