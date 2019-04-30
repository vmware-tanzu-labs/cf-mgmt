package configcommands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddSpaceToConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                     string      `long:"org" description:"Org name" required:"true"`
	SpaceName                   string      `long:"space" description:"Space name" required:"true"`
	AllowSSH                    string      `long:"allow-ssh" description:"Enable the application ssh" choice:"true" choice:"false"`
	AllowSSHUntil               string      `long:"allow-ssh-until" description:"Temporarily allow application ssh until options are Days (1D), Hours (5H), or Minutes (10M)"`
	EnableSecurityGroup         string      `long:"enable-security-group" description:"Enable space level security group definitions" choice:"true" choice:"false"`
	EnableUnassignSecurityGroup string      `long:"enable-unassign-security-group" description:"Enable unassigning security groups not in config" choice:"true" choice:"false"`
	IsoSegment                  string      `long:"isolation-segment" description:"Isolation segment assigned to space"`
	ASGs                        []string    `long:"named-asg" description:"Named asg(s) to assign to space, specify multiple times"`
	NamedQuota                  string      `long:"named-quota" description:"Named quota to assign to space"`
	Quota                       SpaceQuota  `group:"quota"`
	Developer                   UserRoleAdd `group:"developer" namespace:"developer"`
	Manager                     UserRoleAdd `group:"manager" namespace:"manager"`
	Auditor                     UserRoleAdd `group:"auditor" namespace:"auditor"`
}

//Execute - adds a named space to the configuration
func (c *AddSpaceToConfigurationCommand) Execute([]string) error {
	c.initConfig()
	spaceConfig := &config.SpaceConfig{
		Org:   c.OrgName,
		Space: c.SpaceName,
	}

	if c.Quota.EnableSpaceQuota == "true" && c.NamedQuota != "" {
		return fmt.Errorf("cannot enable space quota and use named quotas")
	}

	asgConfigs, err := c.ConfigManager.GetASGConfigs()
	if err != nil {
		return err
	}
	errorString := ""

	spaceConfig.RemoveUsers = true

	convertToBool("enable-security-group", &spaceConfig.EnableSecurityGroup, c.EnableSecurityGroup, &errorString)
	convertToBool("enable-unassign-security-group", &spaceConfig.EnableUnassignSecurityGroup, c.EnableUnassignSecurityGroup, &errorString)
	if c.IsoSegment != "" {
		spaceConfig.IsoSegment = c.IsoSegment
	}

	updateSpaceQuotaConfig(spaceConfig, c.Quota, &errorString)
	spaceConfig.NamedQuota = c.NamedQuota

	spaceConfig.ASGs = addToSlice(spaceConfig.ASGs, c.ASGs, &errorString)
	validateASGsExist(asgConfigs, spaceConfig.ASGs, &errorString)
	c.updateUsers(spaceConfig, &errorString)

	c.sshConfig(spaceConfig, &errorString)
	if errorString != "" {
		return errors.New(errorString)
	}

	if err := config.NewManager(c.ConfigDirectory).AddSpaceToConfig(spaceConfig); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been updated", c.OrgName, c.SpaceName))
	return nil
}

func (c *AddSpaceToConfigurationCommand) sshConfig(spaceConfig *config.SpaceConfig, errorString *string) {
	if strings.EqualFold(c.AllowSSH, "true") && c.AllowSSHUntil != "" {
		*errorString += fmt.Sprintf("\nCannot set --allow-ssh and --allow-ssh-until")
		return
	}
	if strings.EqualFold(c.AllowSSH, "true") {
		spaceConfig.AllowSSH = true
		spaceConfig.AllowSSHUntil = ""
	} else {
		spaceConfig.AllowSSH = false
	}
	if c.AllowSSHUntil != "" {
		t, err := config.FutureTime(time.Now(), c.AllowSSHUntil)
		if err != nil {
			*errorString += fmt.Sprintf("\n%s", err.Error())
		}
		spaceConfig.AllowSSHUntil = t
		spaceConfig.AllowSSH = false
	}

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
