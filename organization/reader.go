package organization

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

func NewReader(client CFClient, cfg config.Reader, peek bool) Reader {
	return &DefaultReader{
		Cfg:    cfg,
		Client: client,
		Peek:   peek,
	}
}

//DefaultManager -
type DefaultReader struct {
	Cfg    config.Reader
	Client CFClient
	Peek   bool
	orgs   []cfclient.Org
}

func (m *DefaultReader) init() error {
	if m.orgs == nil {
		orgs, err := m.Client.ListOrgs()
		if err != nil {
			return err
		}
		m.orgs = orgs
	}
	return nil
}

func (m *DefaultReader) ClearOrgList() {
	m.orgs = nil
}

func (m *DefaultReader) AddOrgToList(org cfclient.Org) {
	if m.orgs == nil {
		m.orgs = []cfclient.Org{}
	}
	m.orgs = append(m.orgs, org)
}

func (m *DefaultReader) GetOrgGUID(orgName string) (string, error) {
	org, err := m.FindOrg(orgName)
	if err != nil {
		return "", err
	}
	return org.Guid, nil
}

//FindOrg -
func (m *DefaultReader) FindOrg(orgName string) (cfclient.Org, error) {
	orgs, err := m.ListOrgs()
	if err != nil {
		return cfclient.Org{}, err
	}
	for _, theOrg := range orgs {
		if strings.EqualFold(theOrg.Name, orgName) {
			return theOrg, nil
		}
	}
	if m.Peek {
		return cfclient.Org{
			Name: orgName,
			Guid: fmt.Sprintf("%s-dry-run-org-guid", orgName),
		}, nil
	}
	return cfclient.Org{}, fmt.Errorf("org %q not found", orgName)
}

//FindOrgByGUID -
func (m *DefaultReader) FindOrgByGUID(orgGUID string) (cfclient.Org, error) {
	orgs, err := m.ListOrgs()
	if err != nil {
		return cfclient.Org{}, err
	}
	for _, theOrg := range orgs {
		if theOrg.Guid == orgGUID {
			return theOrg, nil
		}
	}
	if m.Peek {
		return cfclient.Org{
			Guid: orgGUID,
			Name: fmt.Sprintf("%s-dry-run-org-name", orgGUID),
		}, nil
	}
	return cfclient.Org{}, fmt.Errorf("org %q not found", orgGUID)
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultReader) ListOrgs() ([]cfclient.Org, error) {
	err := m.init()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(m.orgs))
	return m.orgs, nil
}

func (m *DefaultReader) GetOrgByGUID(orgGUID string) (cfclient.Org, error) {
	return m.Client.GetOrgByGuid(orgGUID)
}
