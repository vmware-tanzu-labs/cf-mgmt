package privatedomain

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/xchapter7x/lo"
)

func NewManager(domainClient CFDomainClient, jobClient CFJobClient, orgReader organizationreader.Reader, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:          cfg,
		OrgReader:    orgReader,
		DomainClient: domainClient,
		JobClient:    jobClient,
		Peek:         peek,
	}
}

// DefaultManager -
type DefaultManager struct {
	Cfg          config.Reader
	OrgReader    organizationreader.Reader
	DomainClient CFDomainClient
	JobClient    CFJobClient
	Peek         bool
}

func (m *DefaultManager) CreatePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	allPrivateDomains, err := m.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		org, err := m.OrgReader.FindOrg(orgConfig.Org)
		if err != nil {
			return err
		}
		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.PrivateDomains {
			if existingPrivateDomain, ok := allPrivateDomains[privateDomain]; ok {
				owningOrgGUID := existingPrivateDomain.Relationships.Organization.Data.GUID
				if org.GUID != owningOrgGUID {
					existingOrg, err := m.OrgReader.FindOrgByGUID(owningOrgGUID)
					if err != nil {
						return err
					}
					return fmt.Errorf("Private Domain %s already exists in org [%s]", privateDomain, existingOrg.Name)
				}
			} else {
				newPrivateDomain, err := m.CreatePrivateDomain(org, privateDomain)
				if err != nil {
					return err
				}
				allPrivateDomains[newPrivateDomain.Name] = newPrivateDomain
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemovePrivateDomains {
			orgPrivateDomains, err := m.ListOrgOwnedPrivateDomains(org.GUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomain := range orgPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					err = m.DeletePrivateDomain(privateDomain)
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
	for _, orgConfig := range orgConfigs {
		org, err := m.OrgReader.FindOrg(orgConfig.Org)
		if err != nil {
			return err
		}
		orgSharedPrivateDomains, err := m.ListOrgSharedPrivateDomains(org.GUID)
		if err != nil {
			return err
		}

		lo.G.Debugf("Org %s Shared Domains %+v", orgConfig.Org, reflect.ValueOf(orgSharedPrivateDomains).MapKeys())

		for _, privateDomainName := range orgConfig.SharedPrivateDomains {
			if _, ok := orgSharedPrivateDomains[privateDomainName]; !ok {
				if privateDomain, ok := privateDomains[privateDomainName]; ok {
					err = m.SharePrivateDomain(org, privateDomain)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("Private Domain [%s] is not defined", privateDomainName)
				}
			} else {
				lo.G.Debugf("Org %s already contains shared private domain %s", orgConfig.Org, privateDomainName)
				delete(orgSharedPrivateDomains, privateDomainName)
			}
		}

		if orgConfig.RemoveSharedPrivateDomains {
			lo.G.Debugf("Org %s Shared Domains to be removed %+v", orgConfig.Org, reflect.ValueOf(orgSharedPrivateDomains).MapKeys())
			for _, privateDomain := range orgSharedPrivateDomains {
				err = m.RemoveSharedPrivateDomain(org, privateDomain)
				if err != nil {
					return err
				}
			}
		} else {
			lo.G.Debugf("Shared private domains will not be removed for org [%s], must set enable-remove-shared-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultManager) ListAllPrivateDomains() (map[string]*resource.Domain, error) {
	domains, err := m.DomainClient.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(domains))
	privateDomainMap := make(map[string]*resource.Domain)
	for _, privateDomain := range domains {
		privateDomainMap[privateDomain.Name] = privateDomain
	}
	return privateDomainMap, nil
}

func (m *DefaultManager) CreatePrivateDomain(org *resource.Organization, privateDomain string) (*resource.Domain, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create private domain %s for org %s", privateDomain, org.Name)
		return &resource.Domain{
			GUID: "dry-run-guid",
			Name: privateDomain,
			Relationships: resource.DomainRelationships{
				Organization: resource.ToOneRelationship{
					Data: &resource.Relationship{
						GUID: org.GUID,
					},
				},
			},
		}, nil
	}
	lo.G.Infof("Creating Private Domain %s for Org %s", privateDomain, org.Name)
	r := resource.NewDomainCreate(privateDomain)
	r.Organization = &resource.ToOneRelationship{
		Data: &resource.Relationship{
			GUID: org.GUID,
		},
	}
	return m.DomainClient.Create(context.Background(), r)
}

func (m *DefaultManager) SharePrivateDomain(org *resource.Organization, domain *resource.Domain) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Share private domain %s for org %s", domain.Name, org.Name)
		return nil
	}
	lo.G.Infof("Share private domain %s for org %s", domain.Name, org.Name)
	_, err := m.DomainClient.Share(context.Background(), domain.GUID, org.GUID)
	return err
}

func (m *DefaultManager) ListOrgSharedPrivateDomains(orgGUID string) (map[string]*resource.Domain, error) {
	orgSharedPrivateDomainMap := make(map[string]*resource.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID != privateDomain.Relationships.Organization.Data.GUID {
			orgSharedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgSharedPrivateDomainMap, nil
}

func (m *DefaultManager) listOrgPrivateDomains(orgGUID string) ([]*resource.Domain, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return nil, nil
	}
	privateDomains, err := m.DomainClient.ListForOrganizationAll(context.Background(), orgGUID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "listOrgPrivateDomains")
	}

	lo.G.Debug("Total private domains returned :", len(privateDomains))
	return privateDomains, nil
}

func (m *DefaultManager) ListOrgOwnedPrivateDomains(orgGUID string) (map[string]*resource.Domain, error) {
	orgOwnedPrivateDomainMap := make(map[string]*resource.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID == privateDomain.Relationships.Organization.Data.GUID {
			orgOwnedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgOwnedPrivateDomainMap, nil
}

func (m *DefaultManager) DeletePrivateDomain(domain *resource.Domain) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Delete private domain %s", domain.Name)
		return nil
	}
	lo.G.Infof("Delete private domain %s", domain.Name)
	jobGUID, err := m.DomainClient.Delete(context.Background(), domain.GUID)
	if err != nil {
		return err
	}
	return m.JobClient.PollComplete(context.Background(), jobGUID, nil)
}

func (m *DefaultManager) RemoveSharedPrivateDomain(org *resource.Organization, domain *resource.Domain) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Unshare private domain %s for org %s", domain.Name, org.Name)
		return nil
	}
	lo.G.Infof("Unshare private domain %s for org %s", domain.Name, org.Name)
	return m.DomainClient.UnShare(context.Background(), domain.GUID, org.GUID)
}
