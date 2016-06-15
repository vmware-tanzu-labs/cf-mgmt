package space

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/parnurzeal/gorequest"
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
	var input InputUpdateSpaces
	var space Resource
	var ldapMgr ldap.Manager
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		if input, err = m.loadUpdateFile(f); err == nil {
			if space, err = m.FindSpace(input.Org, input.Space); err == nil {
				lo.G.Info("Processing space", space.Entity.Name)
				if input.AllowSSH != space.Entity.AllowSSH {
					if err = m.updateSpaceSSH(input.AllowSSH, space); err != nil {
						return
					}
				}

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
						lo.G.Info("User", groupUser.UserID, "doesn't exist so creating")
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
	var res *http.Response
	url := fmt.Sprintf("https://api.%s/v2/spaces/%s/%s", m.SysDomain, space.MetaData.GUID, role)
	var body string
	var errs []error
	request := gorequest.New()
	put := request.Put(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+m.Token)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	put.Send(sendString)
	if res, _, errs = put.End(); len(errs) == 0 && res.StatusCode == http.StatusCreated {
		return
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return
}

func (m *DefaultSpaceManager) updateSpaceSSH(value bool, space Resource) (err error) {
	var res *http.Response
	url := fmt.Sprintf("https://api.%s/v2/spaces/%s", m.SysDomain, space.MetaData.GUID)
	var body string
	var errs []error
	request := gorequest.New()
	put := request.Put(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+m.Token)
	sendString := fmt.Sprintf(`{"allow_ssh":%t}`, value)
	put.Send(sendString)
	if res, _, errs = put.End(); len(errs) == 0 && res.StatusCode == http.StatusCreated {
		return
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return
}

//CreateSpace -
func (m *DefaultSpaceManager) CreateSpace(orgName, spaceName string) (space Resource, err error) {
	var res *http.Response
	url := fmt.Sprintf("https://api.%s/v2/spaces", m.SysDomain)
	var body string
	var errs []error
	var orgGUID string
	if orgGUID, err = m.getOrgGUID(orgName); err == nil {
		request := gorequest.New()
		post := request.Post(url)
		post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		post.Set("Authorization", "BEARER "+m.Token)
		sendString := fmt.Sprintf(`{"name":"%s", "organization_guid":"%s"}`, spaceName, orgGUID)
		post.Send(sendString)
		if res, body, errs = post.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
			space, err = m.FindSpace(orgGUID, spaceName)
		} else if len(errs) > 0 {
			err = errs[0]
		} else {
			err = fmt.Errorf(body)
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
	var spaceNames []string
	var orgName string
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaces.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		if orgName, spaceNames, err = m.loadCreateFile(f); err == nil {
			if len(spaceNames) == 0 {
				lo.G.Info("No spaces in config file", f)
			}
			if err = m.fetchSpaces(orgName); err == nil {
				for _, spaceName := range spaceNames {
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
func (m *DefaultSpaceManager) loadCreateFile(configFile string) (orgName string, spaces []string, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err == nil {
		c := &InputCreateSpaces{}
		if err = yaml.Unmarshal(data, c); err == nil {
			orgName = c.Org
			spaces = c.Spaces
		}
	}
	return
}

func (m *DefaultSpaceManager) loadUpdateFile(configFile string) (input InputUpdateSpaces, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err == nil {
		c := &InputUpdateSpaces{}
		if err = yaml.Unmarshal(data, c); err == nil {
			input = *c
		}
	}
	return
}

func (m *DefaultSpaceManager) fetchSpaces(orgName string) (err error) {
	var res *http.Response
	var body string
	var errs []error
	var orgGUID string
	if orgGUID, err = m.getOrgGUID(orgName); err == nil {
		spacesURL := fmt.Sprintf("https://api.%s/v2/organizations/%s/spaces?inline-relations-depth=1", m.SysDomain, orgGUID)
		request := gorequest.New()
		get := request.Get(spacesURL)
		get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		get.Set("Authorization", "BEARER "+m.Token)

		if res, body, errs = get.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
			spaceResources := new(Resources)
			if err = json.Unmarshal([]byte(body), &spaceResources); err == nil {
				m.Spaces = spaceResources.Resource
			}
		} else if len(errs) > 0 {
			err = errs[0]
		} else {
			err = fmt.Errorf(body)
		}
	}
	return
}
