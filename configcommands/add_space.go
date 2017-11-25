package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddSpaceToConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName             string      `long:"org" description:"Org name" required:"true"`
	SpaceName           string      `long:"space" description:"Space name" required:"true"`
	AllowSSH            string      `long:"allow-ssh" description:"Enable the Space Quota in the config" choice:"true" choice:"false"`
	EnableSecurityGroup string      `long:"enable-security-group" description:"Enable space level security group definitions" choice:"true" choice:"false"`
	IsoSegment          string      `long:"isolation-segment" description:"Isolation segment assigned to space"`
	ASGs                []string    `long:"named-asg" description:"Named asg(s) to assign to space, specify multiple times"`
	Quota               SpaceQuota  `group:"quota"`
	Developer           UserRoleAdd `group:"developer" namespace:"developer"`
	Manager             UserRoleAdd `group:"manager" namespace:"manager"`
	Auditor             UserRoleAdd `group:"auditor" namespace:"auditor"`
}

//Execute - adds a named space to the configuration
func (c *AddSpaceToConfigurationCommand) Execute([]string) error {
	spaceConfig := &config.SpaceConfig{
		Org:   c.OrgName,
		Space: c.SpaceName,
	}
	asgConfigs, err := c.ConfigManager.GetASGConfigs()
	if err != nil {
		return err
	}
	errorString := ""

	spaceConfig.RemoveUsers = true

	convertToBool("allow-ssh", &spaceConfig.AllowSSH, c.AllowSSH, &errorString)
	convertToBool("enable-security-group", &spaceConfig.EnableSecurityGroup, c.EnableSecurityGroup, &errorString)
	if c.IsoSegment != "" {
		spaceConfig.IsoSegment = c.IsoSegment
	}

	spaceConfig.ASGs = addToSlice(spaceConfig.ASGs, c.ASGs, &errorString)
	validateASGsExist(asgConfigs, spaceConfig.ASGs, &errorString)
	updateSpaceQuotaConfig(spaceConfig, c.Quota, &errorString)
	c.updateUsers(spaceConfig, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := config.NewManager(c.ConfigDirectory).AddSpaceToConfig(spaceConfig); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been updated", c.OrgName, c.SpaceName))
	return nil
}

func (c *AddSpaceToConfigurationCommand) updateUsers(spaceConfig *config.SpaceConfig, errorString *string) {
	addUsersBasedOnRole(&spaceConfig.Developer, spaceConfig.GetDeveloperGroups(), &c.Developer, errorString)
	addUsersBasedOnRole(&spaceConfig.Auditor, spaceConfig.GetAuditorGroups(), &c.Auditor, errorString)
	addUsersBasedOnRole(&spaceConfig.Manager, spaceConfig.GetManagerGroups(), &c.Manager, errorString)

	spaceConfig.DeveloperGroup = ""
	spaceConfig.ManagerGroup = ""
	spaceConfig.AuditorGroup = ""
}

func (c *AddSpaceToConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
