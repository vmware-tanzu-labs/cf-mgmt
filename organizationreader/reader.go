package organizationreader

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

func NewReader(client CFClient, orgClient CFOrgClient, cfg config.Reader, peek bool) Reader {
	return &DefaultReader{
		Cfg:       cfg,
		OrgClient: orgClient,
		Peek:      peek,
	}
}

// DefaultManager -
type DefaultReader struct {
	Cfg       config.Reader
	OrgClient CFOrgClient
	Peek      bool
	orgs      []*resource.Organization
}

func (m *DefaultReader) init() error {
	if m.orgs == nil {
		orgs, err := m.OrgClient.ListAll(context.Background(), &client.OrganizationListOptions{
			ListOptions: &client.ListOptions{
				PerPage: 5000,
			},
		})
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

func (m *DefaultReader) AddOrgToList(org *resource.Organization) {
	if m.orgs == nil {
		m.orgs = []*resource.Organization{}
	}
	m.orgs = append(m.orgs, org)
}

func (m *DefaultReader) GetOrgGUID(orgName string) (string, error) {
	org, err := m.FindOrg(orgName)
	if err != nil {
		return "", err
	}
	return org.GUID, nil
}

// FindOrg -
func (m *DefaultReader) FindOrg(orgName string) (*resource.Organization, error) {
	orgs, err := m.ListOrgs()
	if err != nil {
		return nil, err
	}
	for _, theOrg := range orgs {
		if strings.EqualFold(theOrg.Name, orgName) {
			return theOrg, nil
		}
	}
	if m.Peek {
		return &resource.Organization{
			Name: orgName,
			GUID: fmt.Sprintf("%s-dry-run-org-guid", orgName),
		}, nil
	}
	return nil, fmt.Errorf("org %q not found", orgName)
}

// FindOrgByGUID -
func (m *DefaultReader) FindOrgByGUID(orgGUID string) (*resource.Organization, error) {
	orgs, err := m.ListOrgs()
	if err != nil {
		return nil, err
	}
	for _, theOrg := range orgs {
		if theOrg.GUID == orgGUID {
			return theOrg, nil
		}
	}
	if m.Peek {
		return &resource.Organization{
			GUID: orgGUID,
			Name: fmt.Sprintf("%s-dry-run-org-name", orgGUID),
		}, nil
	}
	return nil, fmt.Errorf("org %q not found", orgGUID)
}

// ListOrgs : Returns all orgs in the given foundation
func (m *DefaultReader) ListOrgs() ([]*resource.Organization, error) {
	err := m.init()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(m.orgs))
	return m.orgs, nil
}

func (m *DefaultReader) GetDefaultIsolationSegment(org *resource.Organization) (string, error) {
	return m.OrgClient.GetDefaultIsolationSegment(context.Background(), org.GUID)
}
