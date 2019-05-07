package organization

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pkg/errors"
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
	orgs   []cfclient.Org
}

func (m *DefaultManager) init() error {
	if m.orgs == nil {
		orgs, err := m.Client.ListOrgs()
		if err != nil {
			return err
		}
		m.orgs = orgs
	}
	return nil
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
	m.orgs = nil
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
		} else if doesOrgExistFromRename(org.OriginalOrg, currentOrgs) {
			lo.G.Debugf("renamed org [%s] already exists as [%s]", org.Org, org.OriginalOrg)
			if err := m.RenameOrg(org.OriginalOrg, org.Org); err != nil {
				return err
			}
			continue
		} else {
			lo.G.Debugf("[%s] org doesn't exist in list [%v]", org.Org, desiredOrgs)
		}
		if err := m.CreateOrg(org.Org, m.orgNames(currentOrgs)); err != nil {
			return err
		}
	}
	m.orgs = nil
	return nil
}

//DeleteOrgs -
func (m *DefaultManager) DeleteOrgs() error {
	m.orgs = nil
	orgsConfig, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}

	if !orgsConfig.EnableDeleteOrgs {
		lo.G.Debug("Org deletion is not enabled.  Set enable-delete-orgs: true")
		return nil
	}

	renamedOrgs := make(map[string]string)
	configuredOrgs := make(map[string]bool)
	for _, orgName := range orgsConfig.Orgs {
		orgConfig, err := m.Cfg.GetOrgConfig(orgName)
		if err != nil {
			return err
		}
		if orgConfig.OriginalOrg != "" {
			renamedOrgs[orgConfig.OriginalOrg] = orgName
		}
		configuredOrgs[orgName] = true
	}

	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}

	orgsToDelete := make([]cfclient.Org, 0)
	for _, org := range orgs {
		if _, exists := configuredOrgs[org.Name]; !exists {
			if !Matches(org.Name, orgsConfig.ProtectedOrgList()) {
				if _, renamed := renamedOrgs[org.Name]; !renamed {
					orgsToDelete = append(orgsToDelete, org)
				}
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
	m.orgs = nil
	return nil
}

func doesOrgExist(orgName string, orgs []cfclient.Org) bool {
	for _, org := range orgs {
		if strings.EqualFold(org.Name, orgName) {
			return true
		}
	}
	return false
}
func doesOrgExistFromRename(orgName string, orgs []cfclient.Org) bool {
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
	err := m.init()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(m.orgs))
	return m.orgs, nil
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
	org, err := m.Client.CreateOrg(cfclient.OrgRequest{
		Name: orgName,
	})
	if m.orgs == nil {
		m.orgs = []cfclient.Org{}
	}
	m.orgs = append(m.orgs, org)
	return err
}

func (m *DefaultManager) RenameOrg(originalOrgName, newOrgName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: renaming org %s to %s", originalOrgName, newOrgName)
		return nil
	}
	lo.G.Infof("renaming org %s to %s", originalOrgName, newOrgName)
	org, err := m.FindOrg(originalOrgName)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateOrg(org.Guid, cfclient.OrgRequest{
		Name: newOrgName,
	})
	org.Name = newOrgName
	return err
}

func (m *DefaultManager) DeleteOrg(org cfclient.Org) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org %s", org.Name)
		return nil
	}
	lo.G.Infof("Deleting [%s] org", org.Name)
	return m.Client.DeleteOrg(org.Guid, true, false)
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

func (m *DefaultManager) UpdateOrgsMetadata() error {
	supports, err := m.Client.SupportsMetadataAPI()
	if err != nil {
		return errors.Wrap(err, "checking if supports v3 metadata api")
	}
	if !supports {
		lo.G.Infof("Your deployment does not yet support v3 metadata api")
		return nil
	}

	orgConfigList, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	globalCfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	for _, orgConfig := range orgConfigList {
		if orgConfig.Metadata != nil {
			org, err := m.FindOrg(orgConfig.Org)
			if err != nil {
				return err
			}
			metadata := cfclient.Metadata{}
			if orgConfig.Metadata.Labels != nil {
				for key, value := range orgConfig.Metadata.Labels {
					if len(value) > 0 {
						metadata.AddLabel(globalCfg.MetadataPrefix, key, value)
					} else {
						metadata.RemoveLabel(globalCfg.MetadataPrefix, key)
					}
				}
			}
			if orgConfig.Metadata.Annotations != nil {
				for key, value := range orgConfig.Metadata.Annotations {
					if len(value) > 0 {
						metadata.AddAnnotation(fmt.Sprintf("%s/%s", globalCfg.MetadataPrefix, key), value)
					} else {
						metadata.RemoveAnnotation(fmt.Sprintf("%s/%s", globalCfg.MetadataPrefix, key))
					}
				}
			}
			err = m.UpdateOrgMetadata(org, metadata)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) UpdateOrgMetadata(org cfclient.Org, metadata cfclient.Metadata) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update org %s metadata", org.Name)
		return nil
	}
	lo.G.Infof("update org [%s] metadata", org.Name)
	return m.Client.UpdateOrgMetadata(org.Guid, metadata)
}
