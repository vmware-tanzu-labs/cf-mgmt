package organization

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {

	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)
	ldapMgr := ldap.NewManager()
	uaacMgr := uaac.NewManager(sysDomain, uaacToken)
	UserMgr := NewUserManager(cloudController, ldapMgr, uaacMgr)

	return &DefaultOrgManager{
		CloudController: cloudController,
		UAACMgr:         uaacMgr,
		UtilsMgr:        utils.NewDefaultManager(),
		LdapMgr:         ldapMgr,
		UserMgr:         UserMgr,
	}
}

func (m *DefaultOrgManager) GetOrgConfigs(configDir string) ([]*InputUpdateOrgs, error) {
	orgConfigs := []*InputUpdateOrgs{}
	if files, err := m.UtilsMgr.FindFiles(configDir, "orgConfig.yml"); err != nil {
		return nil, err
	} else {
		for _, f := range files {
			input := &InputUpdateOrgs{}
			if err = m.UtilsMgr.LoadFile(f, input); err == nil {
				orgConfigs = append(orgConfigs, input)
			} else {
				lo.G.Error(err)
				return nil, err
			}
		}
	}
	return orgConfigs, nil
}

//CreateQuotas -
func (m *DefaultOrgManager) CreateQuotas(configDir string) error {
	var quotas map[string]string
	var org *cloudcontroller.Org
	var targetQuotaGUID string
	var orgs []*InputUpdateOrgs
	var err error
	if orgs, err = m.GetOrgConfigs(configDir); err != nil {
		return err
	}
	if quotas, err = m.CloudController.ListQuotas(); err != nil {
		return err
	}
	for _, input := range orgs {
		if input.EnableOrgQuota {
			if org, err = m.FindOrg(input.Org); err != nil {
				return err
			} else {
				quotaName := org.Entity.Name
				if quotaGUID, ok := quotas[quotaName]; ok {
					lo.G.Info("Updating quota", quotaName)
					if err = m.CloudController.UpdateQuota(quotaGUID, quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
						lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
						if err := m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, quotaGUID); err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					lo.G.Info("Creating quota", quotaName)
					if targetQuotaGUID, err = m.CloudController.CreateQuota(quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
						lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
						if err := m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, targetQuotaGUID); err != nil {
							return err
						}
					} else {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (m *DefaultOrgManager) GetOrgGUID(orgName string) (string, error) {
	if org, err := m.FindOrg(orgName); err == nil {
		return org.MetaData.GUID, nil
	} else {
		return "", err
	}
}

//CreateOrgs -
func (m *DefaultOrgManager) CreateOrgs(configDir string) error {
	var configFile = configDir + "/orgs.yml"
	lo.G.Info("Processing org file", configFile)
	input := &InputOrgs{}
	if err := m.UtilsMgr.LoadFile(configFile, input); err != nil {
		return err
	}
	if orgs, err := m.CloudController.ListOrgs(); err == nil {
		for _, orgName := range input.Orgs {
			if m.DoesOrgExist(orgName, orgs) {
				lo.G.Info(fmt.Sprintf("[%s] org already exists", orgName))
			} else {
				lo.G.Info(fmt.Sprintf("Creating [%s] org", orgName))
				if err := m.CloudController.CreateOrg(orgName); err != nil {
					return err
				}
			}
		}
	} else {
		return err
	}

	return nil
}

func (m *DefaultOrgManager) DoesOrgExist(orgName string, orgs []*cloudcontroller.Org) bool {
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			return true
		}
	}
	return false

}

//FindOrg -
func (m *DefaultOrgManager) FindOrg(orgName string) (*cloudcontroller.Org, error) {
	if orgs, err := m.CloudController.ListOrgs(); err == nil {
		for _, theOrg := range orgs {
			if theOrg.Entity.Name == orgName {
				return theOrg, nil
			}
		}
		return nil, fmt.Errorf("Org named[%s] not found", orgName)
	} else {
		return nil, err
	}
}

//UpdateOrgUsers -
func (m *DefaultOrgManager) UpdateOrgUsers(configDir, ldapBindPassword string) error {

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

	var orgConfigs []*InputUpdateOrgs

	orgConfigs, err = m.GetOrgConfigs(configDir)

	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range orgConfigs {
		err = m.updateOrgUsers(config, input, uaacUsers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultOrgManager) updateOrgUsers(config *ldap.Config, input *InputUpdateOrgs, uaacUsers map[string]string) error {

	var err error
	var org *cloudcontroller.Org

	org, err = m.FindOrg(input.Org)

	if err != nil {
		return err
	}

	lo.G.Info("User sync for org : ", org.Entity.Name)

	err = m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:       org.Entity.Name,
			OrgGUID:       org.MetaData.GUID,
			Role:          "billing_managers",
			LdapGroupName: input.GetBillingManagerGroup(),
			LdapUsers:     input.BillingManager.LdapUser,
			Users:         input.BillingManager.Users,
			RemoveUsers:   input.RemoveUsers,
		})

	if err != nil {
		return err
	}

	err = m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:       org.Entity.Name,
			OrgGUID:       org.MetaData.GUID,
			Role:          "auditors",
			LdapGroupName: input.GetAuditorGroup(),
			LdapUsers:     input.Auditor.LdapUser,
			Users:         input.Auditor.Users,
			RemoveUsers:   input.RemoveUsers,
		})

	if err != nil {
		return err
	}

	err = m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:       org.Entity.Name,
			OrgGUID:       org.MetaData.GUID,
			Role:          "managers",
			LdapGroupName: input.GetManagerGroup(),
			LdapUsers:     input.Manager.LdapUser,
			Users:         input.Manager.Users,
			RemoveUsers:   input.RemoveUsers,
		})

	if err != nil {
		return err
	}
	return nil
}

func (m *DefaultOrgManager) UpdateBillingManagers(config *ldap.Config, org *cloudcontroller.Org, input *InputUpdateOrgs, uaacUsers map[string]string) error {
	if users, err := m.getLdapUsers(config, input.GetBillingManagerGroup(), input.BillingManager.LdapUser); err == nil {
		if err = m.updateLdapUsers(config, org, "billing_managers", uaacUsers, users); err != nil {
			return err
		}
		for _, userID := range input.BillingManager.Users {
			if err := m.addUserToOrgAndRole(userID, org.MetaData.GUID, "billing_managers"); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}

func (m *DefaultOrgManager) getLdapUsers(config *ldap.Config, groupName string, userList []string) ([]ldap.User, error) {
	users := []ldap.User{}
	if groupName != "" {
		if groupUsers, err := m.LdapMgr.GetUserIDs(config, groupName); err == nil {
			users = append(users, groupUsers...)
		} else {
			return nil, err
		}
	}
	for _, user := range userList {
		if ldapUser, err := m.LdapMgr.GetUser(config, user); err == nil {
			if ldapUser != nil {
				users = append(users, *ldapUser)
			}
		} else {
			return nil, err
		}
	}
	return users, nil
}

func (m *DefaultOrgManager) updateLdapUsers(config *ldap.Config, org *cloudcontroller.Org, role string, uaacUsers map[string]string, users []ldap.User) error {
	for _, user := range users {
		userID := user.UserID
		externalID := user.UserDN
		if config.Origin != "ldap" {
			userID = user.Email
			externalID = user.Email
		}
		if _, userExists := uaacUsers[strings.ToLower(userID)]; userExists {
			lo.G.Info("User", userID, "already exists")
		} else {
			if userID != "" {
				lo.G.Info("User", userID, "doesn't exist so creating in UAA")
				if err := m.UAACMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
					return err
				} else {
					uaacUsers[userID] = userID
				}
			}
		}
		if userID != "" {
			if err := m.addUserToOrgAndRole(userID, org.MetaData.GUID, role); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DefaultOrgManager) addUserToOrgAndRole(userID, orgGUID, role string) error {
	lo.G.Info("Adding user to groups")
	if err := m.CloudController.AddUserToOrg(userID, orgGUID); err != nil {
		return err
	}
	if err := m.CloudController.AddUserToOrgRole(userID, role, orgGUID); err != nil {
		return err
	}
	return nil
}
