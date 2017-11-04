package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateOrgConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                    string   `long:"org" description:"Org name" required:"true"`
	PrivateDomains             []string `long:"private-domain" description:"Private Domain(s) to add, specify muliple times"`
	PrivateDomainsToRemove     []string `long:"private-domain-to-remove" description:"Private Domain(s) to remove, specify muliple times"`
	EnableRemovePrivateDomains string   `long:"enable-remove-private-domains" description:"Enable removing private domains" choice:"true" choice:"false"`
	Quota                      struct {
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
	BillingManager User `group:"billing-manager" namespace:"billing-manager"`
	Manager        User `group:"manager" namespace:"manager"`
	Auditor        User `group:"auditor" namespace:"auditor"`
}

type User struct {
	LDAPUsers          []string `long:"ldap-user" description:"Ldap User to add, specify muliple times"`
	LDAPUsersToRemove  []string `long:"ldap-user-to-remove" description:"Ldap User to remove, specify muliple times"`
	Users              []string `long:"user" description:"User to add, specify muliple times"`
	UsersToRemove      []string `long:"user-to-remove" description:"User to remove, specify muliple times"`
	SamlUsers          []string `long:"saml-user" description:"SAML user to add, specify muliple times"`
	SamlUsersToRemove  []string `long:"saml-user-to-remove" description:"SAML user to remove, specify muliple times"`
	LDAPGroups         []string `long:"ldap-group" description:"User to add, specify muliple times"`
	LDAPGroupsToRemove []string `long:"ldap-group-to-remove" description:"User to remove, specify muliple times"`
}

//Execute - updates org quota configuration`
func (c *UpdateOrgConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	orgConfig, err := c.ConfigManager.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}
	errorString := ""
	c.updatePrivateDomainConfig(orgConfig, &errorString)
	c.updateQuotaConfig(orgConfig, &errorString)

	//TODO Map all the users.....

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveOrgConfig(orgConfig); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org [%s] has been updated", c.OrgName))
	return nil
}

func (c *UpdateOrgConfigurationCommand) updatePrivateDomainConfig(orgConfig *config.OrgConfig, errorString *string) {
	orgConfig.PrivateDomains = append(orgConfig.PrivateDomains, c.PrivateDomains...)
	privateDomainToRemove := SliceToMap(c.PrivateDomainsToRemove)
	var updatedPrivateDomains []string
	for _, privateDomain := range orgConfig.PrivateDomains {
		if _, ok := privateDomainToRemove[privateDomain]; !ok {
			updatedPrivateDomains = append(updatedPrivateDomains, privateDomain)
		}
	}
	orgConfig.PrivateDomains = updatedPrivateDomains
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

func SliceToMap(theSlice []string) map[string]string {
	theMap := make(map[string]string)
	for _, val := range theSlice {
		theMap[val] = val
	}
	return theMap
}

func (c *UpdateOrgConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
