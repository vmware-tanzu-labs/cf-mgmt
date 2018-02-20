package securitygroup

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(cloudController cloudcontroller.Manager, cfg config.Reader) Manager {
	return &DefaultSecurityGroupManager{
		Cfg:             cfg,
		CloudController: cloudController,
	}

}

//CreateApplicationSecurityGroups -
func (m *DefaultSecurityGroupManager) CreateApplicationSecurityGroups() error {
	sgs, err := m.CloudController.ListSecurityGroups()
	if err != nil {
		return err
	}
	securityGroupConfigs, err := m.Cfg.GetASGConfigs()
	if err != nil {
		return err
	}
	defaultSecurityGroupConfigs, err := m.Cfg.GetDefaultASGConfigs()
	if err != nil {
		return err
	}
	err = m.processSecurityGroups(securityGroupConfigs, sgs)
	if err != nil {
		return err
	}
	err = m.processSecurityGroups(defaultSecurityGroupConfigs, sgs)
	if err != nil {
		return err
	}

	return nil
}

//AssignDefaultSecurityGroups -
func (m *DefaultSecurityGroupManager) AssignDefaultSecurityGroups() error {
	sgs, err := m.CloudController.ListSecurityGroups()
	if err != nil {
		return err
	}
	globalConfig, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	for _, runningGroup := range globalConfig.RunningSecurityGroups {
		if group, ok := sgs[runningGroup]; ok {
			if !group.DefaultRunning {
				lo.G.Infof("assigning security group %s as running security group", runningGroup)
				err = m.CloudController.AssignRunningSecurityGroup(group.GUID)
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("Security Group %s does not exist", runningGroup)
		}
	}

	for _, stagingGroup := range globalConfig.StagingSecurityGroups {
		if group, ok := sgs[stagingGroup]; ok {
			if !group.DefaultStaging {
				lo.G.Infof("assigning security group %s as staging security group", stagingGroup)
				err = m.CloudController.AssignStagingSecurityGroup(group.GUID)
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("Security Group %s does not exist", stagingGroup)
		}
	}

	if globalConfig.EnableUnassignSecurityGroups {
		for groupName, group := range sgs {
			if group.DefaultRunning && !m.Contains(globalConfig.RunningSecurityGroups, groupName) {
				lo.G.Infof("unassigning security group %s as running security group", groupName)
				err = m.CloudController.UnassignRunningSecurityGroup(group.GUID)
				if err != nil {
					return err
				}
			}
			if group.DefaultStaging && !m.Contains(globalConfig.StagingSecurityGroups, groupName) {
				lo.G.Infof("unassigning security group %s as staging security group", groupName)
				err = m.CloudController.UnassignStagingSecurityGroup(group.GUID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *DefaultSecurityGroupManager) Contains(list []string, groupName string) bool {
	groupNameToUpper := strings.ToUpper(groupName)
	for _, v := range list {
		if strings.ToUpper(v) == groupNameToUpper {
			return true
		}
	}
	return false
}

func (m *DefaultSecurityGroupManager) processSecurityGroups(securityGroupConfigs []config.ASGConfig, sgs map[string]cloudcontroller.SecurityGroupInfo) error {
	for _, input := range securityGroupConfigs {
		sgName := input.Name

		// For every named security group
		// Check if it's a new group or Update
		if sgInfo, ok := sgs[sgName]; ok {
			match, err := DoesJsonMatch(sgInfo.Rules, input.Rules)
			if err != nil {
				return err
			}
			if !match {
				lo.G.Info("Updating security group", sgName)
				if err := m.CloudController.UpdateSecurityGroup(sgInfo.GUID, sgName, input.Rules); err != nil {
					return err
				}
			}
		} else {
			lo.G.Info("Creating security group", sgName)
			if _, err := m.CloudController.CreateSecurityGroup(sgName, input.Rules); err != nil {
				return err
			}
		}
	}

	return nil
}
