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
	return &DefaultOrgManager{
		UAACMgr:         uaac.NewManager(sysDomain, uaacToken),
		CloudController: cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token),
		UtilsMgr:        utils.NewDefaultManager(),
		LdapMgr:         ldap.NewManager(),
	}
}

//CreateQuotas -
func (m *DefaultOrgManager) CreateQuotas(configDir string) (err error) {
	var quotas map[string]string
	var org *cloudcontroller.Org
	var targetQuotaGUID string
	if quotas, err = m.CloudController.ListQuotas(); err == nil {
		files, _ := m.UtilsMgr.FindFiles(configDir, "orgConfig.yml")
		for _, f := range files {
			input := &InputUpdateOrgs{}
			if err = m.UtilsMgr.LoadFile(f, input); err == nil {
				if input.EnableOrgQuota {
					if org, err = m.FindOrg(input.Org); err == nil {
						quotaName := org.Entity.Name
						if quotaGUID, ok := quotas[quotaName]; ok {
							lo.G.Info("Updating quota", quotaName)
							if err = m.CloudController.UpdateQuota(quotaGUID, quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
								m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, quotaGUID)
							}
						} else {
							lo.G.Info("Creating quota", quotaName)
							if targetQuotaGUID, err = m.CloudController.CreateQuota(quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
								m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, targetQuotaGUID)
							}
						}
					}
				}
			}
		}
	}

	if err != nil {
		lo.G.Error(err)
	}
	return
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
			if m.doesOrgExist(orgName, orgs) {
				lo.G.Info(fmt.Sprintf("[%s] org already exists", orgName))
			} else {
				lo.G.Info(fmt.Sprintf("Creating [%s] org", orgName))
				err = m.CloudController.CreateOrg(orgName)
			}
		}
	} else {
		return err
	}

	return nil
}

func (m *DefaultOrgManager) doesOrgExist(orgName string, orgs []*cloudcontroller.Org) (result bool) {
	result = false
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			result = true
			return
		}
	}
	return

}

//FindOrg -
func (m *DefaultOrgManager) FindOrg(orgName string) (*cloudcontroller.Org, error) {
	if orgs, err := m.CloudController.ListOrgs(); err == nil {
		for _, theOrg := range orgs {
			fmt.Println(theOrg.Entity.Name)
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
func (m *DefaultOrgManager) UpdateOrgUsers(configDir, ldapBindPassword string) (err error) {
	var org *cloudcontroller.Org
	var config *ldap.Config

	if config, err = m.LdapMgr.GetConfig(configDir, ldapBindPassword); err != nil {
		return
	}

	if config.Enabled {
		files, _ := m.UtilsMgr.FindFiles(configDir, "orgConfig.yml")
		for _, f := range files {
			lo.G.Info("Processing org file", f)
			input := &InputUpdateOrgs{}
			if err = m.UtilsMgr.LoadFile(f, input); err == nil {
				if org, err = m.FindOrg(input.Org); err == nil {
					lo.G.Info("User sync for org", input.Org)
					if err = m.updateUsers(config, org, "managers", input.ManagerGroup); err != nil {
						return
					}
					if err = m.updateUsers(config, org, "auditors", input.AuditorGroup); err != nil {
						return
					}
					if err = m.updateUsers(config, org, "billing_managers", input.BillingManagerGroup); err != nil {
						return
					}
				}
			}
		}

	}
	return
}

func (m *DefaultOrgManager) updateUsers(config *ldap.Config, org *cloudcontroller.Org, role, groupName string) (err error) {
	var groupUsers []ldap.User
	var uaacUsers map[string]string
	if uaacUsers, err = m.UAACMgr.ListUsers(); err != nil {
		return
	}
	if groupName != "" {
		lo.G.Info("Getting users for group", groupName)
		if groupUsers, err = m.LdapMgr.GetUserIDs(config, groupName); err == nil {
			for _, groupUser := range groupUsers {
				if _, userExists := uaacUsers[strings.ToLower(groupUser.UserID)]; userExists {
					lo.G.Info("User", groupUser.UserID, "already exists")
				} else {
					lo.G.Info("User", groupUser.UserID, "doesn't exist so creating in UAA")
					if err = m.UAACMgr.CreateLdapUser(groupUser.UserID, groupUser.Email, groupUser.UserDN); err != nil {
						return
					}
				}
				lo.G.Info("Adding user to groups")
				orgGUID := org.MetaData.GUID
				userName := groupUser.UserID
				if err = m.CloudController.AddUserToOrg(userName, orgGUID); err != nil {
					return
				}
				lo.G.Info("Adding", userName, "to", org.Entity.Name, "with role", role)
				err = m.CloudController.AddUserToOrgRole(userName, role, orgGUID)
				return
			}

		}
	}
	return
}
