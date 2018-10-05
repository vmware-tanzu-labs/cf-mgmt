package organization

import (
	"fmt"
	"regexp"
	"strings"

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
		} else {
			lo.G.Debugf("[%s] org doesn't exist in list [%v]", org.Org, desiredOrgs)
		}
		if err := m.CreateOrg(org.Org, m.orgNames(currentOrgs)); err != nil {
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
		if err := m.DeleteOrg(org); err != nil {
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
		if strings.EqualFold(org.Name, orgName) {
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
	if m.Peek {
		return cfclient.Org{
			Name: orgName,
			Guid: fmt.Sprintf("%s-dry-run-org-guid", orgName),
		}, nil
	}
	return cfclient.Org{}, fmt.Errorf("org %q not found", orgName)
}

//FindOrgByGUID -
func (m *DefaultManager) FindOrgByGUID(orgGUID string) (cfclient.Org, error) {
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
func (m *DefaultManager) ListOrgs() ([]cfclient.Org, error) {
	orgs, err := m.Client.ListOrgs()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(orgs))
	return orgs, nil
}

func (m *DefaultManager) orgNames(orgs []cfclient.Org) []string {
	var orgNames []string
	for _, org := range orgs {
		orgNames = append(orgNames, org.Name)
	}

	return orgNames
}

func (m *DefaultManager) CreateOrg(orgName string, currentOrgs []string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org %s as it doesn't exist in %v", orgName, currentOrgs)
		return nil
	}
	lo.G.Infof("create org %s as it doesn't exist in %v", orgName, currentOrgs)
	_, err := m.Client.CreateOrg(cfclient.OrgRequest{
		Name: orgName,
	})
	return err
}

func (m *DefaultManager) DeleteOrg(org cfclient.Org) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org %s", org.Name)
		return nil
	}
	lo.G.Infof("Deleting [%s] org", org.Name)
	return m.Client.DeleteOrg(org.Guid, true, true)
}

func (m *DefaultManager) DeleteOrgByName(orgName string) error {
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if org.Name == orgName {
			return m.DeleteOrg(org)
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
