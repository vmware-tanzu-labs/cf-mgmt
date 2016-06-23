package space

import (
	"encoding/json"
	"fmt"

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
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputUpdateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			if space, err = m.FindSpace(input.Org, input.Space); err == nil {
				if ldapMgr, err = ldap.NewDefaultManager(configDir); err == nil {
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
	var orgName string
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
						lo.G.Info(fmt.Sprintf("Creating [%s] space in [%s] org", spaceName, orgName))
						m.CreateSpace(orgName, spaceName)
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
