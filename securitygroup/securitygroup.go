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

//func (m *DefaultSecurityGroupManager) FindSecurityGroup(securitygroupName string) (*cloudcontroller.SecurityGroup, error) {
func (m *DefaultSecurityGroupManager) FindSecurityGroup(securitygroupName string) (string, error) {
	//	ListSecurityGroups() (map[string]string, error)

	securityGroups, err := m.CloudController.ListSecurityGroups()
	if err != nil {
		return "", err
	}

	for _, theSecurityGroup := range securityGroups {
		if theSecurityGroup == securitygroupName {
			return theSecurityGroup, nil
		}
	}
	return "", fmt.Errorf("security group [%s] not found", securitygroupName)
}

//////////
/*
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		if !input.EnableSecurityGroup {
			continue
		}
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			continue
		}
		sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
		sgs, err := m.CloudController.ListSecurityGroups()
		if err != nil {
			continue
		}
		if sgGUID, ok := sgs[sgName]; ok {
			lo.G.Info("Updating security group", sgName)
			if err := m.CloudController.UpdateSecurityGroup(sgGUID, sgName, input.SecurityGroupContents); err != nil {
				continue
			}
			lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
			m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgGUID)
		} else {
			lo.G.Info("Creating security group", sgName)
			targetSGGUID, err := m.CloudController.CreateSecurityGroup(sgName, input.SecurityGroupContents)
			if err != nil {
				continue
			}
			lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
			m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, targetSGGUID)
		}

		// iterate through and assign named security groups to the space - ensuring that they are up to date is
		// done elsewhere.

		lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
		for _, securityGroupName := range input.ASGs {
			if sgGUID, ok := sgs[securityGroupName]; ok {
				lo.G.Info("Binding security group", securityGroupName, "to space", space.Entity.Name)
				m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgGUID)
			}
		}

	}
	return nil
}







*/

/*
//UpdateSpaces -
func (m *DefaultSecurityGroupManager) UpdateSpaces(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			continue
		}
		lo.G.Info("Processing space", space.Entity.Name)
		if input.AllowSSH != space.Entity.AllowSSH {
			if err := m.CloudController.UpdateSpaceSSH(input.AllowSSH, space.MetaData.GUID); err != nil {
				return err
			}
		}
	}
	return nil
}

//UpdateSpaceUsers -
func (m *DefaultSecurityGroupManager) UpdateSpaceUsers(configDir, ldapBindPassword string) error {
	config, err := m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	uaacUsers, err := m.UAACMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(config, &input, uaacUsers); err != nil {
			return err
		}
	}

	return nil
}

//FindSpace -
func (m *DefaultSecurityGroupManager) FindSpace(orgName, spaceName string) (*cloudcontroller.Space, error) {
	orgGUID, err := m.OrgMgr.GetOrgGUID(orgName)
	if err != nil {
		return nil, err
	}
	spaces, err := m.CloudController.ListSpaces(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, theSpace := range spaces {
		if theSpace.Entity.Name == spaceName {
			return theSpace, nil
		}
	}
	return nil, fmt.Errorf("space [%s] not found in org [%s]", spaceName, orgName)
}

//CreateSpaces -
func (m *DefaultSecurityGroupManager) CreateSpaces(configDir, ldapBindPassword string) error {
	configSpaceList, err := m.Cfg.Spaces()
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {
		if len(input.Spaces) == 0 {
			continue
		}
		orgGUID, err := m.OrgMgr.GetOrgGUID(input.Org)
		if err != nil {
			return err
		}
		spaces, err := m.CloudController.ListSpaces(orgGUID)
		if err != nil {
			continue
		}
		for _, spaceName := range input.Spaces {
			if m.doesSpaceExist(spaces, spaceName) {
				lo.G.Infof("[%s] space already exists", spaceName)
				continue
			}
			lo.G.Infof("Creating [%s] space in [%s] org", spaceName, input.Org)
			if err = m.CloudController.CreateSpace(spaceName, orgGUID); err != nil {
				lo.G.Error(err)
				return err
			}
			if err = m.UpdateSpaceWithDefaults(configDir, spaceName, input.Org, ldapBindPassword); err != nil {
				lo.G.Error(err)
				return err
			}
		}
	}
	return nil
}

func (m *DefaultSecurityGroupManager) UpdateSpaceWithDefaults(configDir, spaceName, orgName, ldapBindPassword string) error {
	defaults, err := m.Cfg.GetSpaceDefaults()
	if err != nil || defaults == nil {
		return nil
	}

	var ldapCfg *ldap.Config
	if ldapBindPassword == "" {
		ldapCfg = &ldap.Config{
			Enabled: false,
		}
	} else {
		if ldapCfg, err = m.LdapMgr.GetConfig(configDir, ldapBindPassword); err != nil {
			lo.G.Error(err)
			return err
		}
	}

	uaacUsers, err := m.UAACMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	defaults.Org = orgName
	defaults.Space = spaceName
	return m.updateSpaceUsers(ldapCfg, defaults, uaacUsers)
}

func (m *DefaultSecurityGroupManager) doesSpaceExist(spaces []*cloudcontroller.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Entity.Name == spaceName {
			return true
		}
	}
	return false
}

func (m *DefaultSecurityGroupManager) DeleteSpaces(configDir string, peekDeletion bool) error {
	configSpaceList, err := m.Cfg.Spaces()
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {

		if !input.EnableDeleteSpaces {
			lo.G.Info(fmt.Sprintf("Space deletion is not enabled for %s.  Set enable-delete-spaces: true in spaces.yml", input.Org))
			continue //Skip all orgs that have not opted-in
		}

		configuredSpaces := make(map[string]bool)
		for _, spaceName := range input.Spaces {
			configuredSpaces[spaceName] = true
		}

		org, err := m.OrgMgr.FindOrg(input.Org)
		if err != nil {
			return err
		}
		spaces, err := m.CloudController.ListSpaces(org.MetaData.GUID)
		if err != nil {
			return err
		}

		spacesToDelete := make([]*cloudcontroller.Space, 0)
		for _, space := range spaces {
			if _, exists := configuredSpaces[space.Entity.Name]; !exists {
				spacesToDelete = append(spacesToDelete, space)
			}
		}

		if peekDeletion {
			for _, space := range spacesToDelete {
				lo.G.Info(fmt.Sprintf("Peek - Would Delete [%s] space in org %s", space.Entity.Name, input.Org))
			}
		} else {
			for _, space := range spacesToDelete {
				lo.G.Info(fmt.Sprintf("Deleting [%s] space in org %s", space.Entity.Name, input.Org))
				if err := m.CloudController.DeleteSpace(space.MetaData.GUID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}*/
