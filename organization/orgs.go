package organization

import (
	"fmt"
	"regexp"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:    cfg,
		Client: client,
		Peek:   peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg    config.Reader
	Client CFClient
	Peek   bool
}

func (m *DefaultManager) GetOrgGUID(orgName string) (string, error) {
	org, err := m.FindOrg(orgName)
	if err != nil {
		return "", err
	}
	return org.Guid, nil
}

//CreateOrgs -
func (m *DefaultManager) CreateOrgs() error {
	desiredOrgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	currentOrgs, err := m.ListOrgs()
	if err != nil {
		return err
	}

	for _, org := range desiredOrgs {
		if doesOrgExist(org.Org, currentOrgs) {
			lo.G.Debugf("[%s] org already exists", org.Org)
			continue
		}
		lo.G.Infof("Creating [%s] org", org.Org)
		if err := m.CreateOrg(org.Org); err != nil {
			return err
		}
	}
	return nil
}

//DeleteOrgs -
func (m *DefaultManager) DeleteOrgs() error {
	orgsConfig, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}

	if !orgsConfig.EnableDeleteOrgs {
		lo.G.Debug("Org deletion is not enabled.  Set enable-delete-orgs: true")
		return nil
	}

	configuredOrgs := make(map[string]bool)
	for _, orgName := range orgsConfig.Orgs {
		configuredOrgs[orgName] = true
	}
	protectedOrgs := append(config.DefaultProtectedOrgs, orgsConfig.ProtectedOrgs...)

	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}

	orgsToDelete := make([]cfclient.Org, 0)
	for _, org := range orgs {
		if _, exists := configuredOrgs[org.Name]; !exists {
			if shouldDeleteOrg(org.Name, protectedOrgs) {
				orgsToDelete = append(orgsToDelete, org)
			} else {
				lo.G.Infof("Protected org [%s] - will not be deleted", org.Name)
			}
		}
	}

	for _, org := range orgsToDelete {
		lo.G.Infof("Deleting [%s] org", org.Name)
		if err := m.DeleteOrg(org.Guid); err != nil {
			return err
		}
	}

	return nil
}

func shouldDeleteOrg(orgName string, protectedOrgs []string) bool {
	for _, protectedOrgName := range protectedOrgs {
		match, _ := regexp.MatchString(protectedOrgName, orgName)
		if match {
			return false
		}
	}
	return true
}

func doesOrgExist(orgName string, orgs []cfclient.Org) bool {
	for _, org := range orgs {
		if org.Name == orgName {
			return true
		}
	}
	return false
}

//FindOrg -
func (m *DefaultManager) FindOrg(orgName string) (cfclient.Org, error) {
	orgs, err := m.ListOrgs()
	if err != nil {
		return cfclient.Org{}, err
	}
	for _, theOrg := range orgs {
		if theOrg.Name == orgName {
			return theOrg, nil
		}
	}
	return cfclient.Org{}, fmt.Errorf("org %q not found", orgName)
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultManager) ListOrgs() ([]cfclient.Org, error) {
	orgs, err := m.Client.ListOrgs()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(orgs))
	return orgs, nil
}

func (m *DefaultManager) CreateOrg(orgName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org %s", orgName)
		return nil
	}
	_, err := m.Client.CreateOrg(cfclient.OrgRequest{
		Name: orgName,
	})
	return err
}

func (m *DefaultManager) DeleteOrg(orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org with GUID %s", orgGUID)
		return nil
	}
	return m.Client.DeleteOrg(orgGUID, true, true)
}

func (m *DefaultManager) DeleteOrgByName(orgName string) error {
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if org.Name == orgName {
			if m.Peek {
				lo.G.Infof("[dry-run]: delete org %s", orgName)
				return nil
			}
			return m.DeleteOrg(org.Guid)
		}
	}
	return fmt.Errorf("org[%s] not found", orgName)
}

func (m *DefaultManager) UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error) {
	return m.Client.UpdateOrg(orgGUID, orgRequest)
}

func (m *DefaultManager) GetOrgByGUID(orgGUID string) (cfclient.Org, error) {
	return m.Client.GetOrgByGuid(orgGUID)
}
