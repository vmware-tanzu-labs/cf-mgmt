package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateSpaceQuotaConfigurationCommand struct {
	BaseConfigCommand
	OrgName                 string `long:"org" description:"Org name" required:"true"`
	SpaceName               string `long:"space" description:"Space name" required:"true"`
	EnableSpaceQuota        string `long:"enable-space-quota" description:"Enable the Space Quota in the config (TRUE or FALSE)" required:"true"`
	MemoryLimit             string `long:"memory-limit" description:"An Space's memory limit in Megabytes"`
	InstanceMemoryLimit     string `long:"instance-memory-limit" description:"Global Space Application instance memory limit in Megabytes"`
	TotalRoutes             string `long:"total-routes" description:"Total Routes capacity for a space"`
	TotalServices           string `long:"total-services" description:"Total Services capacity for a space"`
	PaidServicesAllowed     string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an space (TRUE or FALSE)"`
	TotalPrivateDomains     string `long:"total-private-domains" description:"Total Private Domain capacity for an Space"`
	TotalReservedRoutePorts string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Space"`
	TotalServiceKeys        string `long:"total-service-keys" description:"Total Service Keys capacity for an Space"`
	AppInstanceLimit        string `long:"app-instance-limit" description:"Total Service Keys capacity for an Space"`
}

//Execute - updates space quota configuration`
func (c *UpdateSpaceQuotaConfigurationCommand) Execute(args []string) error {
	cfg := config.NewManager(c.ConfigDirectory)
	spaceConfig, err := cfg.GetSpaceConfig(c.OrgName, c.SpaceName)
	if err != nil {
		return err
	}
	errorString := ""
	convertToBool("enable-org-quota", &spaceConfig.EnableSpaceQuota, c.EnableSpaceQuota, &errorString)
	convertToInt("memory-limit", &spaceConfig.MemoryLimit, c.MemoryLimit, &errorString)
	convertToInt("instance-memory-limit", &spaceConfig.InstanceMemoryLimit, c.InstanceMemoryLimit, &errorString)
	convertToInt("total-routes", &spaceConfig.TotalRoutes, c.TotalRoutes, &errorString)
	convertToInt("total-services", &spaceConfig.TotalServices, c.TotalServices, &errorString)
	convertToBool("paid-service-plans-allowed", &spaceConfig.PaidServicePlansAllowed, c.PaidServicesAllowed, &errorString)
	convertToInt("total-private-domains", &spaceConfig.TotalPrivateDomains, c.TotalPrivateDomains, &errorString)
	convertToInt("total-reserved-route-ports", &spaceConfig.TotalReservedRoutePorts, c.TotalReservedRoutePorts, &errorString)
	convertToInt("total-service-keys", &spaceConfig.TotalReservedRoutePorts, c.TotalServiceKeys, &errorString)
	convertToInt("app-instance-limit", &spaceConfig.AppInstanceLimit, c.AppInstanceLimit, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := cfg.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	fmt.Printf("The quota information has been updated for org/space %s/%s", c.OrgName, c.SpaceName)
	return nil
}
