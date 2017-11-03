package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateOrgQuotaConfigurationCommand struct {
	BaseConfigCommand
	OrgName                 string `long:"org" description:"Org name" required:"true"`
	EnableOrgQuota          string `long:"enable-org-quota" description:"Enable the Org Quota in the config (TRUE or FALSE)" required:"true"`
	MemoryLimit             string `long:"memory-limit" description:"An Org's memory limit in Megabytes"`
	InstanceMemoryLimit     string `long:"instance-memory-limit" description:"Global Org Application instance memory limit in Megabytes"`
	TotalRoutes             string `long:"total-routes" description:"Total Routes capacity for an Org"`
	TotalServices           string `long:"total-services" description:"Total Services capacity for an Org"`
	PaidServicesAllowed     string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an org (TRUE or FALSE)"`
	TotalPrivateDomains     string `long:"total-private-domains" description:"Total Private Domain capacity for an Org"`
	TotalReservedRoutePorts string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Org"`
	TotalServiceKeys        string `long:"total-service-keys" description:"Total Service Keys capacity for an Org"`
	AppInstanceLimit        string `long:"app-instance-limit" description:"Total Service Keys capacity for an Org"`
}

//Execute - updates org quota configuration`
func (c *UpdateOrgQuotaConfigurationCommand) Execute(args []string) error {
	cfg := config.NewManager(c.ConfigDirectory)
	orgConfig, err := cfg.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}
	errorString := ""
	convertToBool("enable-org-quota", &orgConfig.EnableOrgQuota, c.EnableOrgQuota, &errorString)
	convertToInt("memory-limit", &orgConfig.MemoryLimit, c.MemoryLimit, &errorString)
	convertToInt("instance-memory-limit", &orgConfig.InstanceMemoryLimit, c.InstanceMemoryLimit, &errorString)
	convertToInt("total-routes", &orgConfig.TotalRoutes, c.TotalRoutes, &errorString)
	convertToInt("total-services", &orgConfig.TotalServices, c.TotalServices, &errorString)
	convertToBool("paid-service-plans-allowed", &orgConfig.PaidServicePlansAllowed, c.PaidServicesAllowed, &errorString)
	convertToInt("total-private-domains", &orgConfig.TotalPrivateDomains, c.TotalPrivateDomains, &errorString)
	convertToInt("total-reserved-route-ports", &orgConfig.TotalReservedRoutePorts, c.TotalReservedRoutePorts, &errorString)
	convertToInt("total-service-keys", &orgConfig.TotalReservedRoutePorts, c.TotalServiceKeys, &errorString)
	convertToInt("app-instance-limit", &orgConfig.AppInstanceLimit, c.AppInstanceLimit, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := cfg.SaveOrgConfig(orgConfig); err != nil {
		return err
	}
	fmt.Printf("The quota information has been updated for org %s", c.OrgName)
	return nil
}
