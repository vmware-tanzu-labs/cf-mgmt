package space

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) Manager {
	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)
	ldapMgr := ldap.NewManager()
	uaacMgr := uaac.NewManager(sysDomain, uaacToken)
	return &DefaultSpaceManager{
		UAACMgr:         uaacMgr,
		CloudController: cloudController,
		OrgMgr:          organization.NewManager(sysDomain, token, uaacToken),
		LdapMgr:         ldapMgr,
		UtilsMgr:        utils.NewDefaultManager(),
		UserMgr:         NewUserManager(cloudController, ldapMgr, uaacMgr),
	}
}

func (m *DefaultSpaceManager) GetSpaceConfigs(configDir string) ([]*InputSpaceConfig, error) {
	spaceDefaults := &InputSpaceConfig{}
	m.UtilsMgr.LoadFile(filepath.Join(configDir, "spaceDefaults.yml"), spaceDefaults)
	files, err := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	if err != nil {
		return nil, err
	}
	var spaceConfigs []*InputSpaceConfig
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputSpaceConfig{}
		if err = m.UtilsMgr.LoadFile(f, input); err != nil {
			return nil, err
		}
		input.Developer.LdapUsers = append(input.Developer.LdapUsers, spaceDefaults.Developer.LdapUsers...)
		input.Developer.Users = append(input.Developer.Users, spaceDefaults.Developer.Users...)
		input.Auditor.LdapUsers = append(input.Auditor.LdapUsers, spaceDefaults.Auditor.LdapUsers...)
		input.Auditor.Users = append(input.Auditor.Users, spaceDefaults.Auditor.Users...)
		input.Manager.LdapUsers = append(input.Manager.LdapUsers, spaceDefaults.Manager.LdapUsers...)
		input.Manager.Users = append(input.Manager.Users, spaceDefaults.Manager.Users...)

		input.Developer.LdapGroups = append(input.GetDeveloperGroups(), spaceDefaults.GetDeveloperGroups()...)
		input.Auditor.LdapGroups = append(input.GetAuditorGroups(), spaceDefaults.GetAuditorGroups()...)
		input.Manager.LdapGroups = append(input.GetManagerGroups(), spaceDefaults.GetManagerGroups()...)

		spaceConfigs = append(spaceConfigs, input)
		if input.EnableSecurityGroup {
			securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
			lo.G.Debug("Loading security group contents", securityGroupFile)
			var bytes []byte
			bytes, err := ioutil.ReadFile(securityGroupFile)
			if err != nil {
				return nil, err
			}
			lo.G.Debug("setting security group contents", string(bytes))
			input.SecurityGroupContents = string(bytes)
		}
	}
	return spaceConfigs, nil
}

//CreateApplicationSecurityGroups -
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) error {
	spaceConfigs, err := m.GetSpaceConfigs(configDir)
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
	}
	return nil
}

//CreateQuotas -
func (m *DefaultSpaceManager) CreateQuotas(configDir string) error {

	spaceConfigs, err := m.GetSpaceConfigs(configDir)
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
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Info("Updating quota", quotaName)
			if err := m.CloudController.UpdateSpaceQuota(space.Entity.OrgGUID, quotaGUID,
				quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err != nil {
				continue
			}
			lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
			m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, quotaGUID)
		} else {
			lo.G.Info("Creating quota", quotaName)
			targetQuotaGUID, err := m.CloudController.CreateSpaceQuota(space.Entity.OrgGUID,
				quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed)
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
	spaceConfigs, err := m.GetSpaceConfigs(configDir)
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

	spaceConfigs, err := m.GetSpaceConfigs(configDir)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(config, input, uaacUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultSpaceManager) updateSpaceUsers(config *ldap.Config, input *InputSpaceConfig, uaacUsers map[string]string) error {
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
		LdapUsers:      input.Developer.LdapUsers,
		Users:          input.Developer.Users,
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
			LdapUsers:      input.Manager.LdapUsers,
			Users:          input.Manager.Users,
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
			LdapUsers:      input.Auditor.LdapUsers,
			Users:          input.Auditor.Users,
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

func (m *DefaultSpaceManager) GetSpaceConfigList(configDir string) ([]InputSpaces, error) {
	files, err := m.UtilsMgr.FindFiles(configDir, "spaces.yml")
	if err != nil {
		return nil, err
	}
	spaceList := []InputSpaces{}
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := InputSpaces{}
		if err := m.UtilsMgr.LoadFile(f, &input); err == nil {
			spaceList = append(spaceList, input)
		}
	}
	return spaceList, nil
}

//CreateSpaces -
func (m *DefaultSpaceManager) CreateSpaces(configDir, ldapBindPassword string) error {
	configSpaceList, err := m.GetSpaceConfigList(configDir)
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
	defaultSpaceConfigFile := filepath.Join(configDir, "spaceDefaults.yml")
	if !m.UtilsMgr.FileOrDirectoryExists(defaultSpaceConfigFile) {
		return nil
	}
	var config *ldap.Config
	var err error
	if ldapBindPassword == "" {
		config = &ldap.Config{
			Enabled: false,
		}
	} else {
		if config, err = m.LdapMgr.GetConfig(configDir, ldapBindPassword); err != nil {
			lo.G.Error(err)
			return err
		}
	}

	uaacUsers, err := m.UAACMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	var defaultSpaceConfig *InputSpaceConfig
	if err = m.UtilsMgr.LoadFile(defaultSpaceConfigFile, &defaultSpaceConfig); err != nil {
		lo.G.Info(defaultSpaceConfigFile, "doesn't exist")
		return nil
	}
	defaultSpaceConfig.Org = orgName
	defaultSpaceConfig.Space = spaceName
	if err = m.updateSpaceUsers(config, defaultSpaceConfig, uaacUsers); err != nil {
		return err
	}

	return nil
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
	configSpaceList, err := m.GetSpaceConfigList(configDir)
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {

		if !input.EnableDeleteSpaces {
			lo.G.Info("Space deletion is not enabled.  Set enable-delete-space: true")
			return nil
		}

		configuredSpaces := make(map[string]bool)
		for _, spaceName := range input.Spaces {
			configuredSpaces[spaceName] = true
		}

		spaces, err := m.CloudController.ListSpaces(input.Org)
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
				lo.G.Info(fmt.Sprintf("Peek - Would Delete [%s] space", space.Entity.Name))
			}
		} else {
			for _, space := range spacesToDelete {
				lo.G.Info(fmt.Sprintf("Deleting [%s] space", space.Entity.Name))
				if err := m.CloudController.DeleteSpace(space.MetaData.GUID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
