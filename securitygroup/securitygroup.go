package securitygroup

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string, cfg config.Reader) Manager {
	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)

	return &DefaultSecurityGroupManager{
		Cfg:             cfg,
		CloudController: cloudController,
		UtilsMgr:        utils.NewDefaultManager(),
	}

}

//CreateApplicationSecurityGroups -
func (m *DefaultSecurityGroupManager) CreateApplicationSecurityGroups(configDir string) error {

	sgs, err := m.CloudController.ListSecurityGroups()
	if err != nil {
		return fmt.Errorf("Could not list security groups")
	}

	securityGroupConfigs, err := m.Cfg.GetASGConfigs()

	for _, input := range securityGroupConfigs {
		sgName := input.Name

		// For every named security group
		// Check if it's a new group or Update
		if sgGUID, ok := sgs[sgName]; ok {
			lo.G.Info("Updating security group", sgName)
			if err := m.CloudController.UpdateSecurityGroup(sgGUID, sgName, input.Rules); err != nil {
				continue
			}
		} else {
			lo.G.Info("Creating security group", sgName)
			_, err := m.CloudController.CreateSecurityGroup(sgName, input.Rules)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
