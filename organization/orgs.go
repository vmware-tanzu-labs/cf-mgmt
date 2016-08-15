package organization

import (
	"encoding/json"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {
	return &DefaultOrgManager{
		SysDomain: sysDomain,
		Token:     token,
		UAACToken: uaacToken,
	}
}

//CreateQuotas -
func (m *DefaultOrgManager) CreateQuotas(configDir string) (err error) {
	var quotas map[string]string
	var org Resource
	var targetQuotaGUID string
	if quotas, err = m.listQuotas(); err == nil {
		files, _ := utils.NewDefaultManager().FindFiles(configDir, "orgConfig.yml")
		for _, f := range files {
			input := &InputUpdateOrgs{}
			if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
				if input.EnableOrgQuota {
					if org, err = m.FindOrg(input.Org); err == nil {
						quotaName := org.Entity.Name
						if quotaGUID, ok := quotas[quotaName]; ok {
							lo.G.Info("Updating quota", quotaName)
							if err = m.updateQuota(quotaGUID, quotaName, input); err == nil {
								lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
								m.updateOrgQuota(org.MetaData.GUID, quotaGUID)
							}
						} else {
							lo.G.Info("Creating quota", quotaName)
							if targetQuotaGUID, err = m.createQuota(quotaName, input); err == nil {
								lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
								m.updateOrgQuota(org.MetaData.GUID, targetQuotaGUID)
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

func (m *DefaultOrgManager) createQuota(quotaName string, quota *InputUpdateOrgs) (quotaGUID string, err error) {
	var body string
	url := fmt.Sprintf("https://api.%s/v2/quota_definitions", m.SysDomain)
	sendString := fmt.Sprintf(`{"name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t}`, quotaName, quota.MemoryLimit, quota.InstanceMemoryLimit, quota.TotalRoutes, quota.TotalServices, quota.PaidServicePlansAllowed)
	if body, err = utils.NewDefaultManager().HTTPPost(url, m.Token, sendString); err == nil {
		quotaResource := new(Resource)
		if err = json.Unmarshal([]byte(body), &quotaResource); err == nil {
			quotaGUID = quotaResource.MetaData.GUID
		}
	}
	return
}
func (m *DefaultOrgManager) updateQuota(quotaGUID, quotaName string, quota *InputUpdateOrgs) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/quota_definitions/%s", m.SysDomain, quotaGUID)
	sendString := fmt.Sprintf(`{"guid":"%s","name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t}`, quotaGUID, quotaName, quota.MemoryLimit, quota.InstanceMemoryLimit, quota.TotalRoutes, quota.TotalServices, quota.PaidServicePlansAllowed)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultOrgManager) updateOrgQuota(orgGUID, quotaGUID string) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/organizations/%s", m.SysDomain, orgGUID)
	sendString := fmt.Sprintf(`{"quota_definition_guid":"%s"}`, quotaGUID)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultOrgManager) listQuotas() (quotas map[string]string, err error) {
	quotas = make(map[string]string)
	var body string
	url := fmt.Sprintf("https://api.%s/v2/quota_definitions", m.SysDomain)
	if body, err = utils.NewDefaultManager().HTTPGet(url, m.Token); err == nil {
		quotaResources := new(Resources)
		if err = json.Unmarshal([]byte(body), &quotaResources); err == nil {
			for _, quota := range quotaResources.Resource {
				quotas[quota.Entity.Name] = quota.MetaData.GUID
			}
		}
	}
	return
}

//AddUser -
func (m *DefaultOrgManager) AddUser(orgName, userName string) (err error) {
	lo.G.Info("Adding", userName, "to", orgName)
	var org Resource
	if org, err = m.FindOrg(orgName); err == nil {
		orgGUID := org.MetaData.GUID
		url := fmt.Sprintf("https://api.%s/v2/organizations/%s/users", m.SysDomain, orgGUID)
		sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
		err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	}
	return
}

//CreateOrgs -
func (m *DefaultOrgManager) CreateOrgs(configDir string) (err error) {
	var configFile = configDir + "/orgs.yml"
	lo.G.Info("Processing org file", configFile)
	input := &InputOrgs{}
	if err = utils.NewDefaultManager().LoadFile(configFile, input); err == nil {
		if len(input.Orgs) == 0 {
			lo.G.Info("No orgs in config file")
		}
		if err = m.fetchOrgs(); err == nil {
			for _, orgName := range input.Orgs {
				if m.doesOrgExist(orgName) {
					lo.G.Info(fmt.Sprintf("[%s] org already exists", orgName))
				} else {
					lo.G.Info(fmt.Sprintf("Creating [%s] org", orgName))
					m.CreateOrg(orgName)
				}
			}
		}
	}
	return
}

func (m *DefaultOrgManager) doesOrgExist(orgName string) (result bool) {
	result = false
	for _, org := range m.Orgs {
		if org.Entity.Name == orgName {
			result = true
			return
		}
	}
	return

}

//CreateOrg -
func (m *DefaultOrgManager) CreateOrg(orgName string) (org Resource, err error) {
	url := fmt.Sprintf("https://api.%s/v2/organizations", m.SysDomain)
	sendString := fmt.Sprintf(`{"name":"%s"}`, orgName)
	if _, err = utils.NewDefaultManager().HTTPPost(url, m.Token, sendString); err == nil {
		org, err = m.FindOrg(orgName)
	}
	return
}

//FindOrg -
func (m *DefaultOrgManager) FindOrg(orgName string) (org Resource, err error) {
	if err = m.fetchOrgs(); err == nil {
		for _, theOrg := range m.Orgs {
			if theOrg.Entity.Name == orgName {
				org = theOrg
				return
			}
		}
	}
	return
}

//UpdateOrgUsers -
func (m *DefaultOrgManager) UpdateOrgUsers(configDir, ldapBindPassword string) (err error) {
	var org Resource
	var ldapMgr ldap.Manager
	if ldapMgr, err = ldap.NewDefaultManager(configDir, ldapBindPassword); err == nil {
		if ldapMgr.IsEnabled() {
			files, _ := utils.NewDefaultManager().FindFiles(configDir, "orgConfig.yml")
			for _, f := range files {
				lo.G.Info("Processing org file", f)
				input := &InputUpdateOrgs{}
				if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
					if org, err = m.FindOrg(input.Org); err == nil {
						uaacMgr := uaac.NewManager(m.SysDomain, m.UAACToken)
						lo.G.Info("User sync for org", input.Org)
						if err = m.updateUsers(ldapMgr, uaacMgr, org, "managers", input.ManagerGroup); err != nil {
							return
						}
						if err = m.updateUsers(ldapMgr, uaacMgr, org, "auditors", input.AuditorGroup); err != nil {
							return
						}
						if err = m.updateUsers(ldapMgr, uaacMgr, org, "billing_managers", input.BillingManagerGroup); err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}

func (m *DefaultOrgManager) updateUsers(ldapMgr ldap.Manager, uaacMgr uaac.Manager, org Resource, role, groupName string) (err error) {
	var groupUsers []ldap.User
	var uaacUsers map[string]string
	if groupName != "" {
		lo.G.Info("Getting users for group", groupName)
		if groupUsers, err = ldapMgr.GetUserIDs(groupName); err == nil {
			if uaacUsers, err = uaacMgr.ListUsers(); err == nil {
				for _, groupUser := range groupUsers {
					if _, userExists := uaacUsers[groupUser.UserID]; userExists {
						lo.G.Info("User", groupUser.UserID, "already exists")
					} else {
						lo.G.Info("User", groupUser.UserID, "doesn't exist so creating in UAA")
						if err = uaacMgr.CreateUser(groupUser.UserID, groupUser.Email, groupUser.UserDN); err != nil {
							return
						}
					}
					lo.G.Info("Adding user to groups")
					if err = m.addRole(groupUser.UserID, role, org); err != nil {
						lo.G.Error(err)
						return
					}
				}
			}
		}
	}
	return
}

func (m *DefaultOrgManager) addRole(userName, role string, org Resource) (err error) {
	orgName := org.Entity.Name
	if err = m.AddUser(orgName, userName); err != nil {
		return
	}
	lo.G.Info("Adding", userName, "to", org.Entity.Name, "with role", role)

	url := fmt.Sprintf("https://api.%s/v2/organizations/%s/%s", m.SysDomain, org.MetaData.GUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultOrgManager) fetchOrgs() (err error) {
	var body string
	url := fmt.Sprintf("https://api.%s/v2/organizations", m.SysDomain)
	if body, err = utils.NewDefaultManager().HTTPGet(url, m.Token); err == nil {
		orgResources := new(Resources)
		if err = json.Unmarshal([]byte(body), &orgResources); err == nil {
			m.Orgs = orgResources.Resource
		}
	}
	return
}
