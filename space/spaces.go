package space

import (
	"fmt"
	"net/url"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(client CFClient, uaaMgr uaa.Manager,
	orgMgr organization.Manager,
	securityGroupMgr securitygroup.Manager,
	cfg config.Reader, peek bool) Manager {
	ldapMgr := ldap.NewManager()
	return &DefaultManager{
		Cfg:              cfg,
		UAAMgr:           uaaMgr,
		Client:           client,
		OrgMgr:           orgMgr,
		LdapMgr:          ldapMgr,
		SecurityGroupMgr: securityGroupMgr,
		Peek:             peek,
		UserMgr:          NewUserManager(client, peek),
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg              config.Reader
	FilePattern      string
	FilePaths        []string
	Client           CFClient
	UAAMgr           uaa.Manager
	OrgMgr           organization.Manager
	SecurityGroupMgr securitygroup.Manager
	LdapMgr          ldap.Manager
	UserMgr          UserMgr
	Peek             bool
}

//CreateApplicationSecurityGroups -
func (m *DefaultManager) CreateApplicationSecurityGroups(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	sgs, err := m.SecurityGroupMgr.ListNonDefaultSecurityGroups()
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
				lo.G.Infof("Binding NAMED security group %s to space %s", securityGroupName, space.Name)
				m.SecurityGroupMgr.AssignSecurityGroupToSpace(space.Guid, sgInfo.Guid)
			} else {
				return fmt.Errorf("Security group [%s] does not exist", securityGroupName)
			}
		}

		if input.EnableSecurityGroup {
			sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
			var sgGUID string
			if sgInfo, ok := sgs[sgName]; ok {
				lo.G.Debug("Updating security group", sgName)
				if err := m.SecurityGroupMgr.UpdateSecurityGroup(sgInfo.Guid, sgName, input.SecurityGroupContents); err != nil {
					return err
				}
				sgGUID = sgInfo.Guid
			} else {
				lo.G.Debug("Creating security group", sgName)
				securityGroup, err := m.SecurityGroupMgr.CreateSecurityGroup(sgName, input.SecurityGroupContents)
				sgs[sgName] = *securityGroup
				if err != nil {
					return err
				}
				sgGUID = securityGroup.Guid
			}
			lo.G.Infof("Binding security group %s to space %s", sgName, space.Name)
			return m.SecurityGroupMgr.AssignSecurityGroupToSpace(space.Guid, sgGUID)
		}
	}
	return nil
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	spaceQuotas, err := m.Client.ListOrgSpaceQuotas(orgGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	for _, quota := range spaceQuotas {
		quotas[quota.Name] = quota.Guid
	}
	return quotas, nil
}

func (m *DefaultManager) UpdateSpaceQuota(quotaGUID string, quota cfclient.SpaceQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %s with %+v", quotaGUID, quota)
		return nil
	}
	_, err := m.Client.UpdateSpaceQuota(quotaGUID, cfclient.SpaceQuotaRequest{
		Name:                    quota.Name,
		OrganizationGuid:        quota.OrganizationGuid,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	return err
}

func (m *DefaultManager) AssignQuotaToSpace(spaceGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quotaGUID %s to spaceGUID %s", quotaGUID, spaceGUID)
		return nil
	}
	return m.Client.AssignSpaceQuota(quotaGUID, spaceGUID)
}

func (m *DefaultManager) CreateSpaceQuota(quota cfclient.SpaceQuota) (*cfclient.SpaceQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %+v", quota)
		return nil, nil
	}
	spaceQuota, err := m.Client.CreateSpaceQuota(cfclient.SpaceQuotaRequest{
		Name:                    quota.Name,
		OrganizationGuid:        quota.OrganizationGuid,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	if err != nil {
		return nil, err
	}
	return spaceQuota, nil
}

//CreateQuotas -
func (m *DefaultManager) CreateQuotas(configDir string) error {
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
		quotaName := space.Name
		quotas, err := m.ListAllSpaceQuotasForOrg(space.OrganizationGuid)
		if err != nil {
			continue
		}

		quota := cfclient.SpaceQuota{
			OrganizationGuid: space.OrganizationGuid, Name: quotaName,
			MemoryLimit:             input.MemoryLimit,
			InstanceMemoryLimit:     input.InstanceMemoryLimit,
			TotalRoutes:             input.TotalRoutes,
			TotalServices:           input.TotalServices,
			NonBasicServicesAllowed: input.PaidServicePlansAllowed,
			TotalReservedRoutePorts: input.TotalReservedRoutePorts,
			TotalServiceKeys:        input.TotalServiceKeys,
			AppInstanceLimit:        input.AppInstanceLimit,
		}
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Debug("Updating quota", quotaName)
			if err := m.UpdateSpaceQuota(quotaGUID, quota); err != nil {
				continue
			}
			lo.G.Infof("Assigning %s to %s", quotaName, space.Name)
			return m.AssignQuotaToSpace(space.Guid, quotaGUID)
		} else {
			lo.G.Debug("Creating quota", quotaName)
			spaceQuota, err := m.CreateSpaceQuota(quota)
			if err != nil {
				continue
			}
			lo.G.Infof("Assigning %s to %s", quotaName, space.Name)
			return m.AssignQuotaToSpace(space.Guid, spaceQuota.Guid)
		}
	}
	return nil
}

func (m *DefaultManager) UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: setting sshAllowed to %v for spaceGUID %s", sshAllowed, spaceGUID)
		return nil
	}
	space, err := m.Client.GetSpaceByGuid(spaceGUID)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateSpace(spaceGUID, cfclient.SpaceRequest{
		Name:             space.Name,
		AllowSSH:         sshAllowed,
		OrganizationGuid: space.OrganizationGuid,
	})
	return err
}

//UpdateSpaces -
func (m *DefaultManager) UpdateSpaces(configDir string) error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			continue
		}
		lo.G.Debug("Processing space", space.Name)
		if input.AllowSSH != space.AllowSSH {
			if err := m.UpdateSpaceSSH(input.AllowSSH, space.Guid); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) ListSpaces(orgGUID string) ([]cfclient.Space, error) {
	spaces, err := m.Client.ListSpacesByQuery(url.Values{
		"organization_guid": []string{orgGUID},
	})
	if err != nil {
		return nil, err
	}
	return spaces, err

}

//FindSpace -
func (m *DefaultManager) FindSpace(orgName, spaceName string) (cfclient.Space, error) {
	orgGUID, err := m.OrgMgr.GetOrgGUID(orgName)
	if err != nil {
		return cfclient.Space{}, err
	}
	spaces, err := m.ListSpaces(orgGUID)
	if err != nil {
		return cfclient.Space{}, err
	}
	for _, theSpace := range spaces {
		if theSpace.Name == spaceName {
			return theSpace, nil
		}
	}
	return cfclient.Space{}, fmt.Errorf("space [%s] not found in org [%s]", spaceName, orgName)
}

func (m *DefaultManager) CreateSpace(spaceName, orgGUID string) error {
	_, err := m.Client.CreateSpace(cfclient.SpaceRequest{
		Name:             spaceName,
		OrganizationGuid: orgGUID,
	})
	return err
}

//CreateSpaces -
func (m *DefaultManager) CreateSpaces(configDir, ldapBindPassword string) error {
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
		spaces, err := m.ListSpaces(orgGUID)
		if err != nil {
			continue
		}
		for _, spaceName := range input.Spaces {
			if m.doesSpaceExist(spaces, spaceName) {
				lo.G.Debugf("[%s] space already exists", spaceName)
				continue
			}
			lo.G.Infof("Creating [%s] space in [%s] org", spaceName, input.Org)
			if err = m.CreateSpace(spaceName, orgGUID); err != nil {
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

func (m *DefaultManager) UpdateSpaceWithDefaults(configDir, spaceName, orgName, ldapBindPassword string) error {
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

	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	defaults.Org = orgName
	defaults.Space = spaceName
	return m.updateSpaceUsers(ldapCfg, defaults, uaaUsers)
}

func (m *DefaultManager) doesSpaceExist(spaces []cfclient.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Name == spaceName {
			return true
		}
	}
	return false
}

func (m *DefaultManager) DeleteSpaces(configDir string) error {
	configSpaceList, err := m.Cfg.Spaces()
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {

		if !input.EnableDeleteSpaces {
			lo.G.Debugf("Space deletion is not enabled for %s.  Set enable-delete-spaces: true in spaces.yml", input.Org)
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
		spaces, err := m.ListSpaces(org.Guid)
		if err != nil {
			return err
		}

		spacesToDelete := make([]cfclient.Space, 0)
		for _, space := range spaces {
			if _, exists := configuredSpaces[space.Name]; !exists {
				spacesToDelete = append(spacesToDelete, space)
			}
		}

		for _, space := range spacesToDelete {
			lo.G.Infof("Deleting [%s] space in org %s", space.Name, input.Org)
			if err := m.DeleteSpace(space.Guid); err != nil {
				return err
			}
		}

	}

	return nil
}

//DeleteSpace - deletes a space based on GUID
func (m *DefaultManager) DeleteSpace(spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete space with GUID %s", spaceGUID)
		return nil
	}
	return m.Client.DeleteSpace(spaceGUID, true, true)
}

//UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers(configDir, ldapBindPassword string) error {
	config, err := m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	uaaUsers, err := m.UAAMgr.ListUsers()
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
		if err := m.updateSpaceUsers(config, &input, uaaUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultManager) updateSpaceUsers(config *ldap.Config, input *config.SpaceConfig, uaaUsers map[string]string) error {
	space, err := m.FindSpace(input.Org, input.Space)
	if err != nil {
		return err
	}
	if err = m.syncSpaceUsers(config, uaaUsers, UpdateUsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.Guid,
		OrgName:        input.Org,
		OrgGUID:        space.OrganizationGuid,
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		ListUsers:      m.UserMgr.ListSpaceDevelopers,
		RemoveUser:     m.UserMgr.RemoveSpaceDeveloperByUsername,
		AddUser:        m.UserMgr.AssociateSpaceDeveloperByUsername,
	}); err != nil {
		return err
	}

	if err = m.syncSpaceUsers(config, uaaUsers,
		UpdateUsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.Guid,
			OrgGUID:        space.OrganizationGuid,
			OrgName:        input.Org,
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			ListUsers:      m.UserMgr.ListSpaceManagers,
			RemoveUser:     m.UserMgr.RemoveSpaceManagerByUsername,
			AddUser:        m.UserMgr.AssociateSpaceManagerByUsername,
		}); err != nil {
		return err
	}
	if err = m.syncSpaceUsers(config, uaaUsers,
		UpdateUsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.Guid,
			OrgGUID:        space.OrganizationGuid,
			OrgName:        input.Org,
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			ListUsers:      m.UserMgr.ListSpaceAuditors,
			RemoveUser:     m.UserMgr.RemoveSpaceAuditorByUsername,
			AddUser:        m.UserMgr.AssociateSpaceAuditorByUsername,
		}); err != nil {
		return err
	}
	return nil
}

//UpdateSpaceUsers Update space users
func (m *DefaultManager) syncSpaceUsers(config *ldap.Config, uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	spaceUsers, err := updateUsersInput.ListUsers(updateUsersInput.SpaceGUID)
	if err != nil {
		return err
	}

	lo.G.Debugf("SpaceUsers before: %v", spaceUsers)
	if config.Enabled {
		var ldapUsers []ldap.User
		ldapUsers, err = m.LdapMgr.GetLdapUsers(config, updateUsersInput.LdapGroupNames, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %v", ldapUsers)
		for _, user := range ldapUsers {
			err = m.updateLdapUser(config, updateUsersInput, uaaUsers, user, spaceUsers)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		if _, userExists := uaaUsers[lowerUserID]; !userExists {
			return fmt.Errorf("user %s doesn't exist in cloud foundry, so must add internal user first", lowerUserID)
		}
		if _, ok := spaceUsers[lowerUserID]; !ok {
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, lowerUserID)
		}
	}

	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaaUsers[lowerUserEmail]; !userExists {
			lo.G.Debug("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err = m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				return err
			} else {
				uaaUsers[userEmail] = userEmail
			}
		}
		if _, ok := spaceUsers[lowerUserEmail]; !ok {
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userEmail); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, lowerUserEmail)
		}
	}
	if updateUsersInput.RemoveUsers {
		lo.G.Debugf("Deleting users for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
		for spaceUser, _ := range spaceUsers {
			err = updateUsersInput.RemoveUser(updateUsersInput.SpaceGUID, spaceUser)
			if err != nil {
				lo.G.Errorf("Cloud controller API error: %s", err)
				return err
			}
		}
	} else {
		lo.G.Debugf("Not removing users. Set enable-remove-users: true to spaceConfig for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
	}

	lo.G.Debugf("SpaceUsers after: %v", spaceUsers)
	return nil
}

func (m *DefaultManager) updateLdapUser(config *ldap.Config, updateUsersInput UpdateUsersInput,
	uaaUsers map[string]string,
	user ldap.User, spaceUsers map[string]string) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	} else {
		if user.Email == "" {
			user.Email = fmt.Sprintf("%s@user.from.ldap.cf", userID)
		}
	}
	userID = strings.ToLower(userID)

	if _, ok := spaceUsers[userID]; !ok {
		lo.G.Debugf("User[%s] not found in: %v", userID, spaceUsers)
		if _, userExists := uaaUsers[userID]; !userExists {
			lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userID)
				return nil
			} else {
				uaaUsers[userID] = userID
			}
		}
		if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
			return err
		}
	} else {
		delete(spaceUsers, userID)
	}
	return nil
}

func (m *DefaultManager) ListSpaceAuditors(spaceGUID string) (map[string]string, error) {
	return m.UserMgr.ListSpaceAuditors(spaceGUID)
}
func (m *DefaultManager) ListSpaceDevelopers(spaceGUID string) (map[string]string, error) {
	return m.UserMgr.ListSpaceDevelopers(spaceGUID)
}
func (m *DefaultManager) ListSpaceManagers(spaceGUID string) (map[string]string, error) {
	return m.UserMgr.ListSpaceManagers(spaceGUID)
}
func (m *DefaultManager) SpaceQuotaByName(name string) (cfclient.SpaceQuota, error) {
	return m.Client.GetSpaceQuotaByName(name)
}
