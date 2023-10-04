package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func NewManager(orgClient CFOrgClient, orgReader organizationreader.Reader, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:       cfg,
		OrgReader: orgReader,
		OrgClient: orgClient,
		Peek:      peek,
	}
}

// DefaultManager -
type DefaultManager struct {
	Cfg       config.Reader
	OrgReader organizationreader.Reader
	OrgClient CFOrgClient
	SpaceMgr  space.Manager
	Peek      bool
}

// CreateOrgs -
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

	orgsYamlList, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}
	orgsSet := map[string]struct{}{}
	for _, org := range orgsYamlList.Orgs {
		orgsSet[org] = struct{}{}
	}

	for _, org := range desiredOrgs {
		if _, ok := orgsSet[org.Org]; !ok {
			return fmt.Errorf("[%s] found in an orgConfig but not in orgs.yml", org.Org)
		}
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

// DeleteOrgs -
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

	orgsToDelete := make([]*resource.Organization, 0)
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
		// if err := m.SpaceMgr.DeleteSpacesForOrg(org.GUID, org.Name); err != nil {
		// 	return err
		// }
		if err := m.DeleteOrg(org); err != nil {
			return err
		}
	}
	m.OrgReader.ClearOrgList()
	return nil
}

func doesOrgExist(orgName string, orgs []*resource.Organization) bool {
	for _, org := range orgs {
		if strings.EqualFold(org.Name, orgName) {
			return true
		}
	}
	return false
}
func doesOrgExistFromRename(orgName string, orgs []*resource.Organization) bool {
	for _, org := range orgs {
		if strings.EqualFold(org.Name, orgName) {
			return true
		}
	}
	return false
}

func (m *DefaultManager) orgNames(orgs []*resource.Organization) []string {
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
	org, err := m.OrgClient.Create(context.Background(), &resource.OrganizationCreate{
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
	_, err = m.updateOrg(org.GUID, &resource.OrganizationUpdate{
		Name: newOrgName,
	})
	org.Name = newOrgName
	return err
}

func (m *DefaultManager) DeleteOrg(org *resource.Organization) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org %s", org.Name)
		return nil
	}
	lo.G.Infof("Deleting [%s] org", org.Name)
	_, err := m.OrgClient.Delete(context.Background(), org.GUID)
	return err
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

func (m *DefaultManager) updateOrg(orgGUID string, orgRequest *resource.OrganizationUpdate) (*resource.Organization, error) {
	return m.OrgClient.Update(context.Background(), orgGUID, orgRequest)
}

func (m *DefaultManager) UpdateOrgsMetadata() error {
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
			if org.Metadata == nil {
				org.Metadata = &resource.Metadata{}
			}
			//clear any labels that start with the prefix
			for key, _ := range org.Metadata.Labels {
				if strings.Contains(key, globalCfg.MetadataPrefix) {
					org.Metadata.Labels[key] = nil
				}
			}
			if orgConfig.Metadata.Labels != nil {
				for key, value := range orgConfig.Metadata.Labels {
					if len(value) > 0 {
						org.Metadata.SetLabel(globalCfg.MetadataPrefix, key, value)
					} else {
						org.Metadata.RemoveLabel(globalCfg.MetadataPrefix, key)
					}
				}
			}
			//clear any Annotations that start with the prefix
			for key, _ := range org.Metadata.Annotations {
				if strings.Contains(key, globalCfg.MetadataPrefix) {
					org.Metadata.Annotations[key] = nil
				}
			}
			if orgConfig.Metadata.Annotations != nil {
				for key, value := range orgConfig.Metadata.Annotations {
					if len(value) > 0 {
						org.Metadata.SetAnnotation(globalCfg.MetadataPrefix, key, value)
					} else {
						org.Metadata.RemoveAnnotation(globalCfg.MetadataPrefix, key)
					}
				}
			}
			_, err = m.updateOrg(org.GUID, &resource.OrganizationUpdate{
				Name:     org.Name,
				Metadata: org.Metadata,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
