package space

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {
	return &DefaultSpaceManager{
		UAACMgr:         uaac.NewManager(sysDomain, uaacToken),
		CloudController: cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token),
		OrgMgr:          organization.NewManager(sysDomain, token, uaacToken),
		LdapMgr:         ldap.NewManager(),
		UtilsMgr:        utils.NewDefaultManager(),
	}
}

func (m *DefaultSpaceManager) GetSpaceConfigs(configDir string) ([]*InputUpdateSpaces, error) {
	if files, err := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml"); err == nil {
		spaceConfigs := []*InputUpdateSpaces{}
		for _, f := range files {
			lo.G.Info("Processing space file", f)
			input := &InputUpdateSpaces{}
			if err = m.UtilsMgr.LoadFile(f, input); err == nil {
				spaceConfigs = append(spaceConfigs, input)
				if input.EnableSecurityGroup {
					securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
					lo.G.Debug("Loading security group contents", securityGroupFile)
					var bytes []byte
					if bytes, err = ioutil.ReadFile(securityGroupFile); err == nil {
						lo.G.Debug("setting security group contents", string(bytes))
						input.SecurityGroupContents = string(bytes)
					} else {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}

		return spaceConfigs, nil
	} else {
		return nil, err
	}
}

//CreateApplicationSecurityGroups -
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) error {
	var targetSGGUID string
	var sgs map[string]string
	var space *cloudcontroller.Space

	if spaceConfigs, err := m.GetSpaceConfigs(configDir); err != nil {
		return err
	} else {
		for _, input := range spaceConfigs {
			if input.EnableSecurityGroup {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
					if sgs, err = m.CloudController.ListSecurityGroups(); err == nil {
						if sgGUID, ok := sgs[sgName]; ok {
							lo.G.Info("Updating security group", sgName)
							if err = m.CloudController.UpdateSecurityGroup(sgGUID, sgName, input.SecurityGroupContents); err == nil {
								lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
								m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgGUID)
							}
						} else {
							lo.G.Info("Creating security group", sgName)
							if targetSGGUID, err = m.CloudController.CreateSecurityGroup(sgName, input.SecurityGroupContents); err == nil {
								lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
								m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, targetSGGUID)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

//CreateQuotas -
func (m *DefaultSpaceManager) CreateQuotas(configDir string) error {
	var quotas map[string]string
	var space *cloudcontroller.Space
	var targetQuotaGUID string

	if spaceConfigs, err := m.GetSpaceConfigs(configDir); err != nil {
		return err
	} else {
		for _, input := range spaceConfigs {
			if input.EnableSpaceQuota {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					quotaName := space.Entity.Name
					if quotas, err = m.CloudController.ListSpaceQuotas(space.Entity.OrgGUID); err == nil {
						if quotaGUID, ok := quotas[quotaName]; ok {
							lo.G.Info("Updating quota", quotaName)
							if err = m.CloudController.UpdateSpaceQuota(space.Entity.OrgGUID, quotaGUID,
								quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								err = m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, quotaGUID)
							}
						} else {
							lo.G.Info("Creating quota", quotaName)
							if targetQuotaGUID, err = m.CloudController.CreateSpaceQuota(space.Entity.OrgGUID,
								quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								err = m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, targetQuotaGUID)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

//UpdateSpaces -
func (m *DefaultSpaceManager) UpdateSpaces(configDir string) error {
	var space *cloudcontroller.Space

	if spaceConfigs, err := m.GetSpaceConfigs(configDir); err != nil {
		return err
	} else {
		for _, input := range spaceConfigs {
			if space, err = m.FindSpace(input.Org, input.Space); err == nil {
				lo.G.Info("Processing space", space.Entity.Name)
				if input.AllowSSH != space.Entity.AllowSSH {
					if err = m.CloudController.UpdateSpaceSSH(input.AllowSSH, space.MetaData.GUID); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

//UpdateSpaceUsers -
func (m *DefaultSpaceManager) UpdateSpaceUsers(configDir, ldapBindPassword string) error {
	var config *ldap.Config
	var uaacUsers map[string]string
	var err error

	config, err = m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	uaacUsers, err = m.UAACMgr.ListUsers()

	if err != nil {
		lo.G.Error(err)
		return err
	}

	var spaceConfigs []*InputUpdateSpaces

	if spaceConfigs, err = m.GetSpaceConfigs(configDir); err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range spaceConfigs {
		if err = m.updateSpaceUsers(config, input, uaacUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultSpaceManager) updateSpaceUsers(config *ldap.Config, input *InputUpdateSpaces, uaacUsers map[string]string) error {
	if space, err := m.FindSpace(input.Org, input.Space); err == nil {
		lo.G.Info("User sync for space", space.Entity.Name)
		if err = m.UpdateSpaceDevelopers(config, space, input, uaacUsers); err != nil {
			lo.G.Error(err)
			return err
		}
		if err = m.UpdateSpaceManagers(config, space, input, uaacUsers); err != nil {
			lo.G.Error(err)
			return err
		}
		if err = m.UpdateSpaceAuditors(config, space, input, uaacUsers); err != nil {
			lo.G.Error(err)
			return err
		}
		return nil
	} else {
		return err
	}
}

//UpdateSpaceDevelopers Retrieves users from the Given `LDAP and adds them as space developers
func (m *DefaultSpaceManager) UpdateSpaceDevelopers(config *ldap.Config, space *cloudcontroller.Space, input *InputUpdateSpaces, uaacUsers map[string]string) error {
	spaceUsers, err := m.CloudController.GetSpaceUsers(space.MetaData.GUID, "developers")
	if err != nil {
		return err
	}
	if config.Enabled {
		var users []ldap.User
		users, err = m.getLdapUsers(config, input.GetDeveloperGroup(), input.Developer.LdapUser)
		if err != nil {
			return err
		}
		for _, user := range users {
			if _, ok := spaceUsers[user.UserID]; !ok {
				err = m.updateLdapUser(config, space, "developers", uaacUsers, user)
				if err != nil {
					return err
				}
			} else {
				delete(spaceUsers, user.UserID)
			}
		}
	} else {
		lo.G.Info("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range input.Developer.Users {
		if _, ok := spaceUsers[userID]; !ok {
			if err = m.addUserToOrgAndRole(userID, space.Entity.OrgGUID, space.MetaData.GUID, "developers"); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, userID)
		}
	}
	for spaceUser, spaceUserGUID := range spaceUsers {
		lo.G.Info(fmt.Sprintf("removing %s from %s", spaceUser, space.Entity.Name))
		err = m.CloudController.RemoveSpaceUser(space.MetaData.GUID, spaceUserGUID, "developers")
		if err != nil {
			lo.G.Error(fmt.Sprintf("Unable to remove user : %s from space developer role in spce : %s", spaceUser, space.Entity.Name))
			lo.G.Error(fmt.Errorf("Cloud controller API error : %s", err))
			return err
		}
	}
	return nil
}

func (m *DefaultSpaceManager) UpdateSpaceManagers(config *ldap.Config, space *cloudcontroller.Space, input *InputUpdateSpaces, uaacUsers map[string]string) error {
	if config.Enabled {
		if users, err := m.getLdapUsers(config, input.GetManagerGroup(), input.Manager.LdapUser); err == nil {
			if _, err = m.updateLdapUsers(config, space, "managers", uaacUsers, users); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			lo.G.Error(err)
			return err
		}
	} else {
		lo.G.Info("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range input.Manager.Users {
		if err := m.addUserToOrgAndRole(userID, space.Entity.OrgGUID, space.MetaData.GUID, "managers"); err != nil {
			lo.G.Error(err)
			return err
		}
	}

	return nil
}

func (m *DefaultSpaceManager) UpdateSpaceAuditors(config *ldap.Config, space *cloudcontroller.Space, input *InputUpdateSpaces, uaacUsers map[string]string) error {
	if config.Enabled {
		if users, err := m.getLdapUsers(config, input.GetAuditorGroup(), input.Auditor.LdapUser); err == nil {
			if _, err = m.updateLdapUsers(config, space, "auditors", uaacUsers, users); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			lo.G.Error(err)
			return err
		}
	} else {
		lo.G.Info("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range input.Auditor.Users {
		if err := m.addUserToOrgAndRole(userID, space.Entity.OrgGUID, space.MetaData.GUID, "auditors"); err != nil {
			lo.G.Error(err)
			return err
		}
	}
	return nil
}

func (m *DefaultSpaceManager) getLdapUsers(config *ldap.Config, groupName string, userList []string) ([]ldap.User, error) {
	users := []ldap.User{}
	if groupName != "" {
		lo.G.Info("Finding LDAP user for group : ", groupName)
		if groupUsers, err := m.LdapMgr.GetUserIDs(config, groupName); err == nil {
			users = append(users, groupUsers...)
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	for _, user := range userList {
		if ldapUser, err := m.LdapMgr.GetUser(config, user); err == nil {
			if ldapUser != nil {
				users = append(users, *ldapUser)
			}
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	return users, nil
}

func (m *DefaultSpaceManager) updateLdapUser(config *ldap.Config,
	space *cloudcontroller.Space,
	role string, uaacUsers map[string]string,
	user ldap.User) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	}
	userID = strings.ToLower(userID)

	if _, userExists := uaacUsers[userID]; userExists {
		lo.G.Info("User : ", userID, "already exists in cloud foundry")
	} else {
		if userID != "" {
			lo.G.Info("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAACMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Error(err)
				return err
			}
			uaacUsers[userID] = userID
		}
	}
	if userID != "" {
		if err := m.addUserToOrgAndRole(userID, space.Entity.OrgGUID, space.MetaData.GUID, role); err != nil {
			lo.G.Error(err)
			return err
		}
	}
	return nil
}

func (m *DefaultSpaceManager) updateLdapUsers(config *ldap.Config,
	space *cloudcontroller.Space,
	role string, uaacUsers map[string]string,
	users []ldap.User) (map[string]string, error) {

	//Keep tracks for current LDAP users that are being iterated in the loop below
	ldpaUserIDMap := make(map[string]string)

	for _, user := range users {
		userID := user.UserID
		externalID := user.UserDN
		if config.Origin != "ldap" {
			userID = user.Email
			externalID = user.Email
		}
		userID = strings.ToLower(userID)
		// Add the current user id to the map, where key and value are the same
		// Using a Map helps us not having to iterate through the same for loop`
		ldpaUserIDMap[userID] = userID

		if _, userExists := uaacUsers[userID]; userExists {
			lo.G.Info("User : ", userID, "already exists in PCF")
		} else {
			if userID != "" {
				lo.G.Info("User", userID, "doesn't exist in PCF, so creating user")
				if err := m.UAACMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
					lo.G.Error(err)
					return nil, err
				}
				uaacUsers[userID] = userID
			}
		}
		if userID != "" {
			if err := m.addUserToOrgAndRole(userID, space.Entity.OrgGUID, space.MetaData.GUID, role); err != nil {
				lo.G.Error(err)
				return nil, err
			}
		}
	}
	return ldpaUserIDMap, nil
}

func (m *DefaultSpaceManager) addUserToOrgAndRole(userID, orgGUID, spaceGUID, role string) error {
	lo.G.Info(fmt.Sprintf("Adding user to org :  %s and space: %s ", orgGUID, spaceGUID))
	if err := m.CloudController.AddUserToOrg(userID, orgGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	lo.G.Info(fmt.Sprintf("Adding user to space: %s  with role: %s", spaceGUID, role))
	if err := m.CloudController.AddUserToSpaceRole(userID, role, spaceGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	return nil
}

//FindSpace -
func (m *DefaultSpaceManager) FindSpace(orgName, spaceName string) (*cloudcontroller.Space, error) {
	if orgGUID, err := m.OrgMgr.GetOrgGUID(orgName); err != nil {
		return nil, err
	} else {
		if spaces, err := m.CloudController.ListSpaces(orgGUID); err == nil {
			for _, theSpace := range spaces {
				if theSpace.Entity.Name == spaceName {
					return &theSpace, nil
				}
			}
			return nil, fmt.Errorf("Space [%s] not found in org [%s]", spaceName, orgName)
		} else {
			return nil, err
		}
	}
}

func (m *DefaultSpaceManager) GetSpaceConfigList(configDir string) ([]InputCreateSpaces, error) {

	if files, err := m.UtilsMgr.FindFiles(configDir, "spaces.yml"); err != nil {
		return nil, err
	} else {
		spaceList := []InputCreateSpaces{}
		for _, f := range files {
			lo.G.Info("Processing space file", f)
			input := InputCreateSpaces{}
			if err := m.UtilsMgr.LoadFile(f, &input); err == nil {
				spaceList = append(spaceList, input)
			}
		}
		return spaceList, nil
	}
}

//CreateSpaces -
func (m *DefaultSpaceManager) CreateSpaces(configDir, ldapBindPassword string) error {

	if configSpaceList, err := m.GetSpaceConfigList(configDir); err != nil {
		return err
	} else {
		for _, input := range configSpaceList {
			if len(input.Spaces) >= 0 {
				var orgGUID string
				if orgGUID, err = m.OrgMgr.GetOrgGUID(input.Org); err != nil {
					return err
				}
				var spaces []cloudcontroller.Space
				if spaces, err = m.CloudController.ListSpaces(orgGUID); err == nil {
					for _, spaceName := range input.Spaces {
						if m.doesSpaceExist(spaces, spaceName) {
							lo.G.Info(fmt.Sprintf("[%s] space already exists", spaceName))
						} else {
							lo.G.Info(fmt.Sprintf("Creating [%s] space in [%s] org", spaceName, input.Org))
							if err = m.CloudController.CreateSpace(spaceName, orgGUID); err == nil {
								if err = m.UpdateSpaceWithDefaults(configDir, spaceName, input.Org, ldapBindPassword); err != nil {
									lo.G.Error(err)
									return err
								}
							} else {
								lo.G.Error(err)
								return err
							}
						}
					}
				}
			}
		}
		return nil
	}
}

func (m *DefaultSpaceManager) UpdateSpaceWithDefaults(configDir, spaceName, orgName, ldapBindPassword string) error {
	defaultSpaceConfigFile := configDir + "/spaceDefaults.yml"
	if m.UtilsMgr.DoesFileExist(defaultSpaceConfigFile) {
		var config *ldap.Config
		var uaacUsers map[string]string
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

		if uaacUsers, err = m.UAACMgr.ListUsers(); err != nil {
			lo.G.Error(err)
			return err
		}

		var defaultSpaceConfig *InputUpdateSpaces

		if err = m.UtilsMgr.LoadFile(defaultSpaceConfigFile, &defaultSpaceConfig); err == nil {
			defaultSpaceConfig.Org = orgName
			defaultSpaceConfig.Space = spaceName
			if err = m.updateSpaceUsers(config, defaultSpaceConfig, uaacUsers); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			lo.G.Error(err)
			return err
		}
	} else {
		lo.G.Info(defaultSpaceConfigFile, "doesn't exist")
		return nil
	}
}

func (m *DefaultSpaceManager) doesSpaceExist(spaces []cloudcontroller.Space, spaceName string) (result bool) {
	result = false
	for _, space := range spaces {
		if space.Entity.Name == spaceName {
			result = true
			return
		}
	}
	return

}
