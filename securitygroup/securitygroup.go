package securitygroup

import (
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
	sgs, err := m.CloudController.ListNonDefaultSecurityGroups()
	if err != nil {
		return err
	}
	securityGroupConfigs, err := m.Cfg.GetASGConfigs()
	if err != nil {
		return err
	}
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
