package organization

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient, orgReader organizationreader.Reader, spaceMgr space.Manager, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:       cfg,
		Client:    client,
		OrgReader: orgReader,
		SpaceMgr:  spaceMgr,
		Peek:      peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg       config.Reader
	OrgReader organizationreader.Reader
	SpaceMgr  space.Manager
	Client    CFClient
	Peek      bool
}

//CreateOrgs -
func (m *DefaultManager) CreateOrgs() error {
	m.OrgReader.ClearOrgList()
	desiredOrgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	currentOrgs, err := m.OrgReader.ListOrgs()
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
	m.OrgReader.ClearOrgList()
	return nil
}

//DeleteOrgs -
func (m *DefaultManager) DeleteOrgs() error {
	m.OrgReader.ClearOrgList()
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

	orgs, err := m.OrgReader.ListOrgs()
	if err != nil {
		return err
	}

	orgsToDelete := make([]cfclient.Org, 0)
	for _, org := range orgs {
		if _, exists := configuredOrgs[org.Name]; !exists {
			if !util.Matches(org.Name, orgsConfig.ProtectedOrgList()) {
				if _, renamed := renamedOrgs[org.Name]; !renamed {
					orgsToDelete = append(orgsToDelete, org)
				}
			} else {
				lo.G.Infof("Protected org [%s] - will not be deleted", org.Name)
			}
		}
	}

	for _, org := range orgsToDelete {
		if err := m.ClearMetadata(org); err != nil {
			return err
		}
		if err := m.SpaceMgr.DeleteSpacesForOrg(org.Guid, org.Name); err != nil {
			return err
		}
		if err := m.DeleteOrg(org); err != nil {
			return err
		}
	}
	m.OrgReader.ClearOrgList()
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
	if err != nil {
		return err
	}
	m.OrgReader.AddOrgToList(org)
	return nil
}

func (m *DefaultManager) RenameOrg(originalOrgName, newOrgName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: renaming org %s to %s", originalOrgName, newOrgName)
		return nil
	}
	lo.G.Infof("renaming org %s to %s", originalOrgName, newOrgName)
	org, err := m.OrgReader.FindOrg(originalOrgName)
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
	orgs, err := m.OrgReader.ListOrgs()
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
			org, err := m.OrgReader.FindOrg(orgConfig.Org)
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
						// For bug in capi that removal doesn't include prefix
						metadata.RemoveAnnotation(key)
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

func (m *DefaultManager) ClearMetadata(org cfclient.Org) error {
	supports, err := m.Client.SupportsMetadataAPI()
	if err != nil {
		return err
	}
	if !supports {
		return nil
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: removing org metadata from org %s", org.Name)
		return nil
	}
	lo.G.Infof("removing org metadata from org %s", org.Name)
	return m.Client.RemoveOrgMetadata(org.Guid)
}
