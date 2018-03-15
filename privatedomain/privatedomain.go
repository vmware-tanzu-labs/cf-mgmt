package privatedomain

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient, orgMgr organization.Manager, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:    cfg,
		OrgMgr: orgMgr,
		Client: client,
		Peek:   peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg    config.Reader
	OrgMgr organization.Manager
	Client CFClient
	Peek   bool
}

func (m *DefaultManager) CreatePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	orgs, err := m.OrgMgr.ListOrgs()
	if err != nil {
		return err
	}
	allPrivateDomains, err := m.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		orgGUID, err := m.getOrgGUID(orgs, orgConfig.Org)
		if err != nil {
			return err
		}
		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.PrivateDomains {
			if existingPrivateDomain, ok := allPrivateDomains[privateDomain]; ok {
				if orgGUID != existingPrivateDomain.OwningOrganizationGuid {
					existingOrgName, _ := m.getOrgName(orgs, existingPrivateDomain.OwningOrganizationGuid)
					msg := fmt.Sprintf("Private Domain %s already exists in org [%s]", privateDomain, existingOrgName)
					lo.G.Error(msg)
					return fmt.Errorf(msg)
				}
				lo.G.Debugf("Private Domain %s already exists for Org %s", privateDomain, orgConfig.Org)
			} else {
				lo.G.Infof("Creating Private Domain %s for Org %s", privateDomain, orgConfig.Org)
				privateDomain, err := m.CreatePrivateDomain(orgGUID, privateDomain)
				if err != nil {
					return err
				}
				allPrivateDomains[privateDomain.Name] = *privateDomain
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemovePrivateDomains {
			lo.G.Debugf("Looking for private domains to remove for org [%s]", orgConfig.Org)
			orgPrivateDomains, err := m.ListOrgOwnedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.DeletePrivateDomain(privateDomainGUID.Guid)
					if err != nil {
						return err
					}
				}
			}
		} else {
			lo.G.Debugf("Private domains will not be removed for org [%s], must set enable-remove-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultManager) SharePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	privateDomains, err := m.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	orgs, err := m.OrgMgr.ListOrgs()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		orgGUID, err := m.getOrgGUID(orgs, orgConfig.Org)
		if err != nil {
			return err
		}
		allSharedPrivateDomains, err := m.ListOrgSharedPrivateDomains(orgGUID)
		if err != nil {
			return err
		}

		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.SharedPrivateDomains {
			if _, ok := allSharedPrivateDomains[privateDomain]; !ok {
				if privateDomainGUID, ok := privateDomains[privateDomain]; ok {
					lo.G.Infof("Sharing Private Domain %s for Org %s", privateDomain, orgConfig.Org)
					err = m.SharePrivateDomain(orgGUID, privateDomainGUID.Guid)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("Private Domain [%s] is not defined", privateDomain)
				}
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemoveSharedPrivateDomains {
			lo.G.Debugf("Looking for shared private domains to remove for org [%s]", orgConfig.Org)
			orgSharedPrivateDomains, err := m.ListOrgSharedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgSharedPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Shared Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.RemoveSharedPrivateDomain(orgGUID, privateDomainGUID.Guid)
					if err != nil {
						return err
					}
				}
			}
		} else {
			lo.G.Debugf("Shared private domains will not be removed for org [%s], must set enable-remove-shared-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultManager) ListAllPrivateDomains() (map[string]cfclient.Domain, error) {
	domains, err := m.Client.ListDomains()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(domains))
	privateDomainMap := make(map[string]cfclient.Domain)
	for _, privateDomain := range domains {
		privateDomainMap[privateDomain.Name] = privateDomain
	}
	return privateDomainMap, nil
}

func (m *DefaultManager) CreatePrivateDomain(orgGUID, privateDomain string) (*cfclient.Domain, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create private domain %s for org GUID %s", privateDomain, orgGUID)
		return nil, nil
	}
	domain, err := m.Client.CreateDomain(privateDomain, orgGUID)
	if err != nil {
		return nil, err
	}
	return domain, nil
}
func (m *DefaultManager) SharePrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	_, err := m.Client.ShareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
	return err
}

func (m *DefaultManager) ListOrgSharedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error) {
	orgSharedPrivateDomainMap := make(map[string]cfclient.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID != privateDomain.OwningOrganizationGuid {
			orgSharedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgSharedPrivateDomainMap, nil
}

func (m *DefaultManager) listOrgPrivateDomains(orgGUID string) ([]cfclient.Domain, error) {
	privateDomains, err := m.Client.ListOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}

	lo.G.Debug("Total private domains returned :", len(privateDomains))
	return privateDomains, nil
}

func (m *DefaultManager) ListOrgOwnedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error) {
	orgOwnedPrivateDomainMap := make(map[string]cfclient.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID == privateDomain.OwningOrganizationGuid {
			orgOwnedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgOwnedPrivateDomainMap, nil
}

func (m *DefaultManager) DeletePrivateDomain(guid string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Delete private domain %s", guid)
		return nil
	}
	return m.Client.DeleteDomain(guid)
}

func (m *DefaultManager) RemoveSharedPrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: remove share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	return m.Client.UnshareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
}

func (m *DefaultManager) getOrgGUID(orgs []cfclient.Org, orgName string) (string, error) {
	for _, org := range orgs {
		if org.Name == orgName {
			return org.Guid, nil
		}
	}
	return "", fmt.Errorf("org %s does not exist", orgName)
}

func (m *DefaultManager) getOrgName(orgs []cfclient.Org, orgGUID string) (string, error) {
	for _, org := range orgs {
		if org.Guid == orgGUID {
			return org.Name, nil
		}
	}
	return "", fmt.Errorf("org for GUID %s does not exist", orgGUID)
}
