package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {
	return &DefaultSpaceManager{
		SysDomain: sysDomain,
		Token:     token,
		UAACToken: uaacToken,
	}
}

//CreateApplicationSecurityGroups -
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) (err error) {
	var contents string
	var targetSGGUID string
	var sgs map[string]string
	var space Resource
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		input := &InputUpdateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			if input.EnableSecurityGroup {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
					if contents, err = m.getSecurityFileContents(securityGroupFile); err == nil {
						sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
						if sgs, err = m.listSecurityGroups(); err == nil {
							if sgGUID, ok := sgs[sgName]; ok {
								lo.G.Info("Updating security group", sgName)
								if err = m.updateSecurityGroup(sgGUID, sgName, contents); err == nil {
									lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
									m.updateSpaceSecurityGroup(space.MetaData.GUID, sgGUID)
								}
							} else {
								lo.G.Info("Creating security group", sgName)
								if targetSGGUID, err = m.createSecurityGroup(sgName, contents); err == nil {
									lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
									m.updateSpaceSecurityGroup(space.MetaData.GUID, targetSGGUID)
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) updateSecurityGroup(sgGUID, sgName, contents string) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/security_groups/%s", m.SysDomain, sgGUID)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultSpaceManager) createSecurityGroup(sgName, contents string) (sgGUID string, err error) {
	var body string
	url := fmt.Sprintf("https://api.%s/v2/security_groups", m.SysDomain)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	if body, err = utils.NewDefaultManager().HTTPPost(url, m.Token, sendString); err == nil {
		sgResource := new(Resource)
		if err = json.Unmarshal([]byte(body), &sgResource); err == nil {
			sgGUID = sgResource.MetaData.GUID
		}
	}
	return
}

func (m *DefaultSpaceManager) updateSpaceSecurityGroup(spaceGUID, sgGUID string) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/security_groups/%s/spaces/%s", m.SysDomain, sgGUID, spaceGUID)
	sendString := ""
	//err = utils.NewDefaultManager().HTTPDelete(url, m.Token, sendString)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultSpaceManager) getSecurityFileContents(securityGroupFile string) (contents string, err error) {
	var f *os.File
	buf := bytes.NewBuffer(nil)

	if f, err = os.Open(securityGroupFile); err == nil {
		io.Copy(buf, f)
		f.Close()
		contents = string(buf.Bytes())
	}
	return
}
func (m *DefaultSpaceManager) listSecurityGroups() (securityGroups map[string]string, err error) {
	securityGroups = make(map[string]string)
	var body string
	url := fmt.Sprintf("https://api.%s/v2/security_groups", m.SysDomain)
	if body, err = utils.NewDefaultManager().HTTPGet(url, m.Token); err == nil {
		sgResources := new(Resources)
		if err = json.Unmarshal([]byte(body), &sgResources); err == nil {
			for _, sg := range sgResources.Resource {
				securityGroups[sg.Entity.Name] = sg.MetaData.GUID
			}
		}
	}
	return
}

//CreateQuotas -
func (m *DefaultSpaceManager) CreateQuotas(configDir string) (err error) {
	var quotas map[string]string
	var space Resource
	var targetQuotaGUID string

	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		input := &InputUpdateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			lo.G.Info("Processing file", f)
			if input.EnableSpaceQuota {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					quotaName := space.Entity.Name
					if quotas, err = m.listQuotas(space.Entity.OrgGUID); err == nil {
						if quotaGUID, ok := quotas[quotaName]; ok {
							lo.G.Info("Updating quota", quotaName)
							if err = m.updateQuota(space.Entity.OrgGUID, quotaGUID, quotaName, input); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								m.updateSpaceQuota(space.MetaData.GUID, quotaGUID)
							}
						} else {
							lo.G.Info("Creating quota", quotaName)
							if targetQuotaGUID, err = m.createQuota(space.Entity.OrgGUID, quotaName, input); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								m.updateSpaceQuota(space.MetaData.GUID, targetQuotaGUID)
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

func (m *DefaultSpaceManager) createQuota(orgGUID, quotaName string, quota *InputUpdateSpaces) (quotaGUID string, err error) {
	var body string
	url := fmt.Sprintf("https://api.%s/v2/space_quota_definitions", m.SysDomain)
	sendString := fmt.Sprintf(`{"name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t,"organization_guid":"%s"}`, quotaName, quota.MemoryLimit, quota.InstanceMemoryLimit, quota.TotalRoutes, quota.TotalServices, quota.PaidServicePlansAllowed, orgGUID)
	if body, err = utils.NewDefaultManager().HTTPPost(url, m.Token, sendString); err == nil {
		quotaResource := new(Resource)
		if err = json.Unmarshal([]byte(body), &quotaResource); err == nil {
			quotaGUID = quotaResource.MetaData.GUID
		}
	}
	return
}
func (m *DefaultSpaceManager) updateQuota(orgGUID, quotaGUID, quotaName string, quota *InputUpdateSpaces) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/space_quota_definitions/%s", m.SysDomain, quotaGUID)
	sendString := fmt.Sprintf(`{"guid":"%s","name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t,"organization_guid":"%s"}`, quotaGUID, quotaName, quota.MemoryLimit, quota.InstanceMemoryLimit, quota.TotalRoutes, quota.TotalServices, quota.PaidServicePlansAllowed, orgGUID)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultSpaceManager) updateSpaceQuota(spaceGUID, quotaGUID string) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/space_quota_definitions/%s/spaces/%s", m.SysDomain, quotaGUID, spaceGUID)
	sendString := ""
	//err = utils.NewDefaultManager().HTTPDelete(url, m.Token, sendString)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultSpaceManager) listQuotas(orgGUID string) (quotas map[string]string, err error) {
	quotas = make(map[string]string)
	var body string
	url := fmt.Sprintf("https://api.%s/v2/organizations/%s/space_quota_definitions", m.SysDomain, orgGUID)
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

//UpdateSpaces -
func (m *DefaultSpaceManager) UpdateSpaces(configDir string) (err error) {
	var space Resource
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputUpdateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			//if input, err = m.loadUpdateFile(f); err == nil {
			if space, err = m.FindSpace(input.Org, input.Space); err == nil {
				lo.G.Info("Processing space", space.Entity.Name)
				if input.AllowSSH != space.Entity.AllowSSH {
					if err = m.updateSpaceSSH(input.AllowSSH, space); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

//UpdateSpaceUsers -
func (m *DefaultSpaceManager) UpdateSpaceUsers(configDir string) (err error) {
	var space Resource
	var ldapMgr ldap.Manager
	if ldapMgr, err = ldap.NewDefaultManager(configDir); err == nil {
		if ldapMgr.IsEnabled() {
			files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
			for _, f := range files {
				lo.G.Info("Processing space file", f)
				input := &InputUpdateSpaces{}
				if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
					if space, err = m.FindSpace(input.Org, input.Space); err == nil {
						uaacMgr := uaac.NewManager(m.SysDomain, m.UAACToken)
						lo.G.Info("User sync for space", space.Entity.Name)
						if err = m.updateUsers(ldapMgr, uaacMgr, space, "developers", input.DeveloperGroup); err != nil {
							return
						}
						if err = m.updateUsers(ldapMgr, uaacMgr, space, "managers", input.ManagerGroup); err != nil {
							return
						}
						if err = m.updateUsers(ldapMgr, uaacMgr, space, "auditors", input.AuditorGroup); err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}
func (m *DefaultSpaceManager) updateUsers(ldapMgr ldap.Manager, uaacMgr uaac.Manager, space Resource, role, groupName string) (err error) {
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
					if err = m.addRole(groupUser.UserID, role, space); err != nil {
						lo.G.Error(err)
						return
					}
				}
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) addRole(userName, role string, space Resource) (err error) {
	orgName := space.Entity.Org.OrgEntity.Name
	orgMgr := organization.NewManager(m.SysDomain, m.Token, m.UAACToken)
	if err = orgMgr.AddUser(orgName, userName); err != nil {
		return
	}
	lo.G.Info("Adding", userName, "to", space.Entity.Name, "with role", role)

	url := fmt.Sprintf("https://api.%s/v2/spaces/%s/%s", m.SysDomain, space.MetaData.GUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

func (m *DefaultSpaceManager) updateSpaceSSH(value bool, space Resource) (err error) {
	url := fmt.Sprintf("https://api.%s/v2/spaces/%s", m.SysDomain, space.MetaData.GUID)
	sendString := fmt.Sprintf(`{"allow_ssh":%t}`, value)
	err = utils.NewDefaultManager().HTTPPut(url, m.Token, sendString)
	return
}

//CreateSpace -
func (m *DefaultSpaceManager) CreateSpace(orgName, spaceName string) (space Resource, err error) {
	var orgGUID string
	if orgGUID, err = m.getOrgGUID(orgName); err == nil {
		url := fmt.Sprintf("https://api.%s/v2/spaces", m.SysDomain)
		sendString := fmt.Sprintf(`{"name":"%s", "organization_guid":"%s"}`, spaceName, orgGUID)
		if _, err = utils.NewDefaultManager().HTTPPost(url, m.Token, sendString); err == nil {
			space, err = m.FindSpace(orgGUID, spaceName)
		}
	}
	return
}

//FindSpace -
func (m *DefaultSpaceManager) FindSpace(orgName, spaceName string) (space Resource, err error) {

	if err = m.fetchSpaces(orgName); err == nil {
		for _, theSpace := range m.Spaces {
			if theSpace.Entity.Name == spaceName {
				space = theSpace
				return
			}
		}
	}
	return
}

//CreateSpaces -
func (m *DefaultSpaceManager) CreateSpaces(configDir string) (err error) {
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaces.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputCreateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			if len(input.Spaces) == 0 {
				lo.G.Info("No spaces in config file", f)
			}
			if err = m.fetchSpaces(input.Org); err == nil {
				for _, spaceName := range input.Spaces {
					if m.doesSpaceExist(spaceName) {
						lo.G.Info(fmt.Sprintf("[%s] space already exists", spaceName))
					} else {
						lo.G.Info(fmt.Sprintf("Creating [%s] space in [%s] org", spaceName, input.Org))
						m.CreateSpace(input.Org, spaceName)
					}
				}
			} else {
				return
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) doesSpaceExist(spaceName string) (result bool) {
	result = false
	for _, space := range m.Spaces {
		if space.Entity.Name == spaceName {
			result = true
			return
		}
	}
	return

}

func (m *DefaultSpaceManager) getOrgGUID(orgName string) (orgGUID string, err error) {
	var org organization.Resource
	orgMgr := organization.NewManager(m.SysDomain, m.Token, m.UAACToken)
	if org, err = orgMgr.FindOrg(orgName); err == nil {
		if org == (organization.Resource{}) {
			err = fmt.Errorf("Org [%s] does not exist", orgName)
			return
		}
		orgGUID = org.MetaData.GUID
	}
	return
}

func (m *DefaultSpaceManager) fetchSpaces(orgName string) (err error) {
	var body string
	var orgGUID string
	if orgGUID, err = m.getOrgGUID(orgName); err == nil {
		url := fmt.Sprintf("https://api.%s/v2/organizations/%s/spaces?inline-relations-depth=1", m.SysDomain, orgGUID)
		if body, err = utils.NewDefaultManager().HTTPGet(url, m.Token); err == nil {
			spaceResources := new(Resources)
			if err = json.Unmarshal([]byte(body), &spaceResources); err == nil {
				m.Spaces = spaceResources.Resource
			}
		}
	}
	return
}
