package space

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string, cfg config.Reader) Manager {
	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)
	ldapMgr := ldap.NewManager()
	uaacMgr := uaac.NewManager(sysDomain, uaacToken)
	return &DefaultSpaceManager{
		Cfg:             cfg,
		UAACMgr:         uaacMgr,
		CloudController: cloudController,
		OrgMgr:          organization.NewManager(sysDomain, token, uaacToken, cfg),
		LdapMgr:         ldapMgr,
		UserMgr:         NewUserManager(cloudController, ldapMgr, uaacMgr),
	}
}

//CreateApplicationSecurityGroups -
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	sgs, err := m.CloudController.ListSecurityGroups()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			return err
		}

		// iterate through and assign named security groups to the space - ensuring that they are up to date is
		// done elsewhere.
		for _, securityGroupName := range input.ASGs {
			lo.G.Debug("Security Group name: " + securityGroupName)
			if sgInfo, ok := sgs[securityGroupName]; ok {
				lo.G.Info("Binding NAMED security group", securityGroupName, "to space", space.Entity.Name)
				m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgInfo.GUID)
			} else {
				return fmt.Errorf("Security group [%s] does not exist", securityGroupName)
			}
		}

		if !input.EnableSecurityGroup {
			continue
		}
		sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
		var sgGUID string
		if sgInfo, ok := sgs[sgName]; ok {
			lo.G.Debug("Updating security group", sgName)
			if err := m.CloudController.UpdateSecurityGroup(sgInfo.GUID, sgName, input.SecurityGroupContents); err != nil {
				return err
			}
			sgGUID = sgInfo.GUID
		} else {
			lo.G.Debug("Creating security group", sgName)
			targetSGGUID, err := m.CloudController.CreateSecurityGroup(sgName, input.SecurityGroupContents)
			sgs[sgName] = cloudcontroller.SecurityGroupInfo{GUID: targetSGGUID, Rules: input.SecurityGroupContents}
			if err != nil {
				return err
			}
			sgGUID = targetSGGUID
		}
		lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
		m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgGUID)
	}
	return nil
}

//CreateQuotas -
func (m *DefaultSpaceManager) CreateQuotas(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		if !input.EnableSpaceQuota {
			continue
		}
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			continue
		}
		quotaName := space.Entity.Name
		quotas, err := m.CloudController.ListAllSpaceQuotasForOrg(space.Entity.OrgGUID)
		if err != nil {
			continue
		}

		quota := cloudcontroller.SpaceQuotaEntity{
			OrgGUID: space.Entity.OrgGUID,
			QuotaEntity: cloudcontroller.QuotaEntity{
				Name:                    quotaName,
				MemoryLimit:             input.MemoryLimit,
				InstanceMemoryLimit:     input.InstanceMemoryLimit,
				TotalRoutes:             input.TotalRoutes,
				TotalServices:           input.TotalServices,
				PaidServicePlansAllowed: input.PaidServicePlansAllowed,
				TotalPrivateDomains:     input.TotalPrivateDomains,
				TotalReservedRoutePorts: input.TotalReservedRoutePorts,
				TotalServiceKeys:        input.TotalServiceKeys,
				AppInstanceLimit:        input.AppInstanceLimit,
			},
		}
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Info("Updating quota", quotaName)
			if err := m.CloudController.UpdateSpaceQuota(quotaGUID, quota); err != nil {
				continue
			}
			lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
			m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, quotaGUID)
		} else {
			lo.G.Info("Creating quota", quotaName)
			targetQuotaGUID, err := m.CloudController.CreateSpaceQuota(quota)
			if err != nil {
				continue
			}
			lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
			m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, targetQuotaGUID)
		}
	}
	return nil
}

//UpdateSpaces -
func (m *DefaultSpaceManager) UpdateSpaces(configDir string) error {
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
func (m *DefaultSpaceManager) UpdateSpaceUsers(configDir, ldapBindPassword string) error {
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

func (m *DefaultSpaceManager) updateSpaceUsers(config *ldap.Config, input *config.SpaceConfig, uaacUsers map[string]string) error {
	space, err := m.FindSpace(input.Org, input.Space)
	if err != nil {
		return err
	}
	if err = m.UserMgr.UpdateSpaceUsers(config, uaacUsers, UpdateUsersInput{
		SpaceName:      space.Entity.Name,
		SpaceGUID:      space.MetaData.GUID,
		OrgName:        input.Org,
		OrgGUID:        space.Entity.OrgGUID,
		Role:           "developers",
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
	}); err != nil {
		return err
	}

	if err = m.UserMgr.UpdateSpaceUsers(config, uaacUsers,
		UpdateUsersInput{
			SpaceName:      space.Entity.Name,
			SpaceGUID:      space.MetaData.GUID,
			OrgGUID:        space.Entity.OrgGUID,
			OrgName:        input.Org,
			Role:           "managers",
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		}); err != nil {
		return err
	}
	if err = m.UserMgr.UpdateSpaceUsers(config, uaacUsers,
		UpdateUsersInput{
			SpaceName:      space.Entity.Name,
			SpaceGUID:      space.MetaData.GUID,
			OrgGUID:        space.Entity.OrgGUID,
			OrgName:        input.Org,
			Role:           "auditors",
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		}); err != nil {
		return err
	}
	return nil
}

//FindSpace -
func (m *DefaultSpaceManager) FindSpace(orgName, spaceName string) (*cloudcontroller.Space, error) {
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
func (m *DefaultSpaceManager) CreateSpaces(configDir, ldapBindPassword string) error {
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

func (m *DefaultSpaceManager) UpdateSpaceWithDefaults(configDir, spaceName, orgName, ldapBindPassword string) error {
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

func (m *DefaultSpaceManager) doesSpaceExist(spaces []*cloudcontroller.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Entity.Name == spaceName {
			return true
		}
	}
	return false
}

func (m *DefaultSpaceManager) DeleteSpaces(configDir string, peekDeletion bool) error {
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
}
