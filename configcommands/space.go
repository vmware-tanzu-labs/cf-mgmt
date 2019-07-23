package configcommands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pivotalservices/cf-mgmt/config"
)

type SpaceConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                     string     `long:"org" description:"Org name" required:"true"`
	SpaceName                   string     `long:"space" description:"Space name" required:"true"`
	AllowSSH                    string     `long:"allow-ssh" description:"Enable the application ssh" choice:"true" choice:"false"`
	AllowSSHUntil               string     `long:"allow-ssh-until" description:"Temporarily allow application ssh until options are Days (1D), Hours (5H), or Minutes (10M)"`
	EnableRemoveUsers           string     `long:"enable-remove-users" description:"Enable removing users from the space" choice:"true" choice:"false"`
	EnableSecurityGroup         string     `long:"enable-security-group" description:"Enable space level security group definitions" choice:"true" choice:"false"`
	EnableUnassignSecurityGroup string     `long:"enable-unassign-security-group" description:"Enable unassigning security groups not in config" choice:"true" choice:"false"`
	IsoSegment                  string     `long:"isolation-segment" description:"Isolation segment assigned to space"`
	ClearIsolationSegment       bool       `long:"clear-isolation-segment" description:"Sets the isolation segment to blank"`
	ASGs                        []string   `long:"named-asg" description:"Named asg(s) to assign to space, specify multiple times"`
	ASGsToRemove                []string   `long:"named-asg-to-remove" description:"Named asg(s) to remove, specify multiple times"`
	NamedQuota                  string     `long:"named-quota" description:"Named quota to assign to space"`
	ClearNamedQuota             bool       `long:"clear-named-quota" description:"Sets the named quota to blank"`
	Quota                       SpaceQuota `group:"quota"`
	Developer                   UserRole   `group:"developer" namespace:"developer"`
	Manager                     UserRole   `group:"manager" namespace:"manager"`
	Auditor                     UserRole   `group:"auditor" namespace:"auditor"`
}

//Execute - updates space configuration`
func (c *SpaceConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	var spaceConfig *config.SpaceConfig
	var err error
	var newSpace bool
	spaceConfig, err = c.ConfigManager.GetSpaceConfig(c.OrgName, c.SpaceName)
	if err != nil {
		spaceConfig = &config.SpaceConfig{
			Org:         c.OrgName,
			Space:       c.SpaceName,
			RemoveUsers: true,
		}
		newSpace = true
	} else {
		newSpace = false
	}

	asgConfigs, err := c.ConfigManager.GetASGConfigs()
	if err != nil {
		return err
	}

	if c.Quota.EnableSpaceQuota == "true" && c.NamedQuota != "" {
		return fmt.Errorf("cannot enable space quota and use named quotas")
	}

	errorString := ""

	convertToBool("enable-remove-users", &spaceConfig.RemoveUsers, c.EnableRemoveUsers, &errorString)
	convertToBool("enable-security-group", &spaceConfig.EnableSecurityGroup, c.EnableSecurityGroup, &errorString)
	convertToBool("enable-unassign-security-group", &spaceConfig.EnableUnassignSecurityGroup, c.EnableUnassignSecurityGroup, &errorString)
	if c.IsoSegment != "" {
		spaceConfig.IsoSegment = c.IsoSegment
	}
	if c.ClearIsolationSegment {
		spaceConfig.IsoSegment = ""
	}

	spaceConfig.ASGs = removeFromSlice(addToSlice(spaceConfig.ASGs, c.ASGs, &errorString), c.ASGsToRemove)
	validateASGsExist(asgConfigs, spaceConfig.ASGs, &errorString)
	updateSpaceQuotaConfig(spaceConfig, c.Quota, &errorString)

	if c.NamedQuota != "" {
		spaceConfig.NamedQuota = c.NamedQuota
	}
	if c.ClearNamedQuota {
		spaceConfig.NamedQuota = ""
	}

	c.updateUsers(spaceConfig, &errorString)
	c.sshConfig(spaceConfig, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	if newSpace {
		fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been updated", c.OrgName, c.SpaceName))
	} else {
		fmt.Println(fmt.Sprintf("The org/space [%s/%s] has been created", c.OrgName, c.SpaceName))
	}
	return nil
}

func (c *SpaceConfigurationCommand) sshConfig(spaceConfig *config.SpaceConfig, errorString *string) {
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

func (c *SpaceConfigurationCommand) updateUsers(spaceConfig *config.SpaceConfig, errorString *string) {
	updateUsersBasedOnRole(&spaceConfig.Developer, spaceConfig.GetDeveloperGroups(), &c.Developer, errorString)
	updateUsersBasedOnRole(&spaceConfig.Auditor, spaceConfig.GetAuditorGroups(), &c.Auditor, errorString)
	updateUsersBasedOnRole(&spaceConfig.Manager, spaceConfig.GetManagerGroups(), &c.Manager, errorString)

	spaceConfig.DeveloperGroup = ""
	spaceConfig.ManagerGroup = ""
	spaceConfig.AuditorGroup = ""
}

func (c *SpaceConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
