package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateSpaceConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName               string   `long:"org" description:"Org name" required:"true"`
	SpaceName             string   `long:"space" description:"Space name" required:"true"`
	AllowSSH              string   `long:"allow-ssh" description:"Enable the Space Quota in the config" choice:"true" choice:"false"`
	EnableRemoveUsers     string   `long:"enable-remove-users" description:"Enable removing users from the space" choice:"true" choice:"false"`
	IsoSegment            string   `long:"isolation-segment" description:"Isolation segment assigned to space"`
	ClearIsolationSegment bool     `long:"clear-isolation-segment" description:"Sets the isolation segment to blank"`
	ASGs                  []string `long:"named-asg" description:"Named asg(s) to assign to space, specify muliple times"`
	ASGsToRemove          []string `long:"named-asg-to-remove" description:"Named asg(s) to remove, specify muliple times"`
	Quota                 struct {
		EnableSpaceQuota        string `long:"enable-space-quota" description:"Enable the Space Quota in the config" choice:"true" choice:"false"`
		MemoryLimit             string `long:"memory-limit" description:"An Space's memory limit in Megabytes"`
		InstanceMemoryLimit     string `long:"instance-memory-limit" description:"Space Application instance memory limit in Megabytes"`
		TotalRoutes             string `long:"total-routes" description:"Total Routes capacity for an Space"`
		TotalServices           string `long:"total-services" description:"Total Services capacity for an Space"`
		PaidServicesAllowed     string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an Space" choice:"true" choice:"false"`
		TotalPrivateDomains     string `long:"total-private-domains" description:"Total Private Domain capacity for an Space"`
		TotalReservedRoutePorts string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Space"`
		TotalServiceKeys        string `long:"total-service-keys" description:"Total Service Keys capacity for an Space"`
		AppInstanceLimit        string `long:"app-instance-limit" description:"Total Service Keys capacity for an Space"`
	} `group:"quota"`
	Developer UserRole `group:"developer" namespace:"developer"`
	Manager   UserRole `group:"manager" namespace:"manager"`
	Auditor   UserRole `group:"auditor" namespace:"auditor"`
}

//Execute - updates space configuration`
func (c *UpdateSpaceConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	spaceConfig, err := c.ConfigManager.GetSpaceConfig(c.OrgName, c.SpaceName)
	if err != nil {
		return err
	}
	errorString := ""

	convertToBool("allow-ssh", &spaceConfig.AllowSSH, c.AllowSSH, &errorString)
	convertToBool("enable-remove-users", &spaceConfig.RemoveUsers, c.EnableRemoveUsers, &errorString)
	if c.IsoSegment != "" {
		spaceConfig.IsoSegment = c.IsoSegment
	}
	if c.ClearIsolationSegment {
		spaceConfig.IsoSegment = ""
	}

	spaceConfig.ASGs = removeFromSlice(append(spaceConfig.ASGs, c.ASGs...), c.ASGsToRemove)

	c.updateQuotaConfig(spaceConfig, &errorString)
	c.updateUsers(spaceConfig, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been updated", c.OrgName, c.SpaceName))
	return nil
}

func (c *UpdateSpaceConfigurationCommand) updateUsers(spaceConfig *config.SpaceConfig, errorString *string) {
	updateUsersBasedOnRole(&spaceConfig.Developer, spaceConfig.GetDeveloperGroups(), &c.Developer)
	updateUsersBasedOnRole(&spaceConfig.Auditor, spaceConfig.GetAuditorGroups(), &c.Auditor)
	updateUsersBasedOnRole(&spaceConfig.Manager, spaceConfig.GetManagerGroups(), &c.Manager)

	spaceConfig.DeveloperGroup = ""
	spaceConfig.ManagerGroup = ""
	spaceConfig.AuditorGroup = ""
}

func (c *UpdateSpaceConfigurationCommand) updateQuotaConfig(spaceConfig *config.SpaceConfig, errorString *string) {
	convertToBool("enable-space-quota", &spaceConfig.EnableSpaceQuota, c.Quota.EnableSpaceQuota, errorString)
	convertToInt("memory-limit", &spaceConfig.MemoryLimit, c.Quota.MemoryLimit, errorString)
	convertToInt("instance-memory-limit", &spaceConfig.InstanceMemoryLimit, c.Quota.InstanceMemoryLimit, errorString)
	convertToInt("total-routes", &spaceConfig.TotalRoutes, c.Quota.TotalRoutes, errorString)
	convertToInt("total-services", &spaceConfig.TotalServices, c.Quota.TotalServices, errorString)
	convertToBool("paid-service-plans-allowed", &spaceConfig.PaidServicePlansAllowed, c.Quota.PaidServicesAllowed, errorString)
	convertToInt("total-private-domains", &spaceConfig.TotalPrivateDomains, c.Quota.TotalPrivateDomains, errorString)
	convertToInt("total-reserved-route-ports", &spaceConfig.TotalReservedRoutePorts, c.Quota.TotalReservedRoutePorts, errorString)
	convertToInt("total-service-keys", &spaceConfig.TotalServiceKeys, c.Quota.TotalServiceKeys, errorString)
	convertToInt("app-instance-limit", &spaceConfig.AppInstanceLimit, c.Quota.AppInstanceLimit, errorString)
}

func (c *UpdateSpaceConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
