package organization

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/parnurzeal/gorequest"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token string) (mgr Manager) {
	return &DefaultOrgManager{
		SysDomain: sysDomain,
		Token:     token,
	}
}

//CreateOrgs -
func (m *DefaultOrgManager) CreateOrgs(configDir string) (err error) {
	var orgNames []string
	var configFile = configDir + "/orgs.yml"
	lo.G.Info("Processing org file", configFile)
	if orgNames, err = m.loadInputFile(configFile); err == nil {
		if len(orgNames) == 0 {
			lo.G.Info("No orgs in config file")
		}
		if err = m.fetchOrgs(); err == nil {
			for _, orgName := range orgNames {
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

func (m *DefaultOrgManager) loadInputFile(configFile string) (orgs []string, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err == nil {
		c := &InputOrgs{}
		if err = yaml.Unmarshal(data, c); err == nil {
			orgs = c.Orgs
		}
	}
	return
}

//CreateOrg -
func (m *DefaultOrgManager) CreateOrg(orgName string) (org Resource, err error) {
	var res *http.Response
	orgsURL := fmt.Sprintf("https://api.%s/v2/organizations", m.SysDomain)

	var body string
	var errs []error
	request := gorequest.New()
	post := request.Post(orgsURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.Set("Authorization", "BEARER "+m.Token)
	sendString := fmt.Sprintf(`{"name":"%s"}`, orgName)
	post.Send(sendString)
	if res, body, errs = post.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
		org, err = m.FindOrg(orgName)
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf(body)
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

func (m *DefaultOrgManager) fetchOrgs() (err error) {
	var res *http.Response
	orgsURL := fmt.Sprintf("https://api.%s/v2/organizations", m.SysDomain)

	var body string
	var errs []error
	request := gorequest.New()
	get := request.Get(orgsURL)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+m.Token)

	if res, body, errs = get.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
		orgResources := new(Resources)
		if err = json.Unmarshal([]byte(body), &orgResources); err == nil {
			m.Orgs = orgResources.Resource
		}
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf(body)
	}
	return
}
