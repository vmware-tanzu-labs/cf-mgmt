package configcommands

import (
	"errors"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type SpaceNamedQuotaConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	Name  string          `long:"name" description:"Name of quota" required:"true"`
	Org   string          `long:"org" description:"Name of org" required:"true"`
	Quota NamedSpaceQuota `group:"quota"`
}

// Execute - updates space named quotas
func (c *SpaceNamedQuotaConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	spaceQuota, err := c.ConfigManager.GetSpaceQuota(c.Name, c.Org)
	if err != nil {
		return err
	}
	if spaceQuota == nil {
		spaceQuota = &config.SpaceQuota{
			Name:                       c.Name,
			Org:                        c.Org,
			TotalReservedRoutePorts:    "unlimited",
			TotalServiceKeys:           "unlimited",
			AppInstanceLimit:           "unlimited",
			AppTaskLimit:               "unlimited",
			InstanceMemoryLimit:        "unlimited",
			TotalRoutes:                "unlimited",
			TotalServices:              "unlimited",
			PaidServicePlansAllowed:    true,
			LogRateLimitBytesPerSecond: "unlimited",
		}
	}
	errorString := ""
	updateSpaceNamedQuotaConfig(spaceQuota, c.Quota, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}
	return c.ConfigManager.SaveSpaceQuota(spaceQuota)
}

func updateSpaceNamedQuotaConfig(namedSpaceQuota *config.SpaceQuota, spaceQuota NamedSpaceQuota, errorString *string) {
	convertToGB("memory-limit", &namedSpaceQuota.MemoryLimit, spaceQuota.MemoryLimit, config.UNLIMITED, errorString)
	convertToGB("instance-memory-limit", &namedSpaceQuota.InstanceMemoryLimit, spaceQuota.InstanceMemoryLimit, config.UNLIMITED, errorString)
	convertToFormattedInt("total-routes", &namedSpaceQuota.TotalRoutes, spaceQuota.TotalRoutes, config.UNLIMITED, errorString)
	convertToFormattedInt("total-services", &namedSpaceQuota.TotalServices, spaceQuota.TotalServices, config.UNLIMITED, errorString)
	convertToBool("paid-service-plans-allowed", &namedSpaceQuota.PaidServicePlansAllowed, spaceQuota.PaidServicesAllowed, errorString)
	convertToFormattedInt("total-reserved-route-ports", &namedSpaceQuota.TotalReservedRoutePorts, spaceQuota.TotalReservedRoutePorts, config.UNLIMITED, errorString)
	convertToFormattedInt("total-service-keys", &namedSpaceQuota.TotalServiceKeys, spaceQuota.TotalServiceKeys, config.UNLIMITED, errorString)
	convertToFormattedInt("app-instance-limit", &namedSpaceQuota.AppInstanceLimit, spaceQuota.AppInstanceLimit, config.UNLIMITED, errorString)
	convertToFormattedInt("app-task-limit", &namedSpaceQuota.AppTaskLimit, spaceQuota.AppTaskLimit, config.UNLIMITED, errorString)
	convertToFormattedInt("app-task-limit", &namedSpaceQuota.LogRateLimitBytesPerSecond, spaceQuota.LogRateLimitBytesPerSecond, config.UNLIMITED, errorString)
}

func (c *SpaceNamedQuotaConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
