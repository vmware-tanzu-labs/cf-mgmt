package space

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/parnurzeal/gorequest"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token string) (mgr Manager) {
	return &DefaultSpaceManager{
		SysDomain: sysDomain,
		Token:     token,
	}
}

//CreateSpace -
func (m *DefaultSpaceManager) CreateSpace(orgGUID, spaceName string) (space Resource, err error) {
	var res *http.Response
	url := fmt.Sprintf("https://api.%s/v2/spaces", m.SysDomain)

	var body string
	var errs []error
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
	return
}

//FindSpace -
func (m *DefaultSpaceManager) FindSpace(orgGUID, spaceName string) (space Resource, err error) {
	if err = m.fetchSpaces(orgGUID); err == nil {
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
	var orgName, orgGUID string
	files, _ := ioutil.ReadDir(configDir)
	for _, f := range files {
		if strings.Contains(f.Name(), "-spaces.yml") {
			lo.G.Info("Processing space file", f.Name())
			if orgName, spaceNames, err = m.loadInputFile(configDir + "/" + f.Name()); err == nil {
				if len(spaceNames) == 0 {
					lo.G.Info("No spaces in config file")
				}
				if orgGUID, err = m.getOrgGUID(orgName); err == nil {
					if err = m.fetchSpaces(orgGUID); err == nil {
						for _, spaceName := range spaceNames {
							if m.doesSpaceExist(spaceName) {
								lo.G.Info(fmt.Sprintf("[%s] space already exists", spaceName))
							} else {
								lo.G.Info(fmt.Sprintf("Creating [%s] space in [%s] org", spaceName, orgName))
								m.CreateSpace(orgGUID, spaceName)
							}
						}
					}
				} else {
					return
				}
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
	orgMgr := organization.NewManager(m.SysDomain, m.Token)
	if org, err = orgMgr.FindOrg(orgName); err == nil {
		if org == (organization.Resource{}) {
			err = fmt.Errorf("Org [%s] does not exist", orgName)
			fmt.Println(err)
			return
		}
		orgGUID = org.MetaData.GUID
	}
	return
}
func (m *DefaultSpaceManager) loadInputFile(configFile string) (orgName string, spaces []string, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err == nil {
		c := &InputSpaces{}
		if err = yaml.Unmarshal(data, c); err == nil {
			orgName = c.Org
			spaces = c.Spaces
		}
	}
	return
}

func (m *DefaultSpaceManager) fetchSpaces(orgGUID string) (err error) {
	var res *http.Response
	spacesURL := fmt.Sprintf("https://api.%s/v2/organizations/%s/spaces", m.SysDomain, orgGUID)

	var body string
	var errs []error
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
	return
}
