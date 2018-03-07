package cloudcontroller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/http"
	"github.com/xchapter7x/lo"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

func NewManager(host, token, version string, peek bool) (Manager, error) {
	c := &cfclient.Config{
		ApiAddress:        host,
		SkipSslValidation: true,
		Token:             token,
		UserAgent:         fmt.Sprintf("cf-mgmt/%s", version),
	}

	client, err := cfclient.NewClient(c)
	if err != nil {
		return nil, err
	}
	return &DefaultManager{
		Client: *client,
		Token:  token,
		Peek:   peek,
		HTTP:   http.NewManager(),
	}, nil
}

func (m *DefaultManager) AddUserToOrg(userName, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to orgGUID %s", userName, orgGUID)
		return nil
	}
	_, err := m.Client.AssociateOrgUserByUsername(orgGUID, userName)
	return err
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
			m.DeleteOrg(org.Guid)
		}
	}
	return errors.New(fmt.Sprintf("org[%s] not found", orgName))
}

//DeleteSpace - deletes a space based on GUID
func (m *DefaultManager) DeleteSpace(spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete space with GUID %s", spaceGUID)
		return nil
	}
	return m.Client.DeleteSpace(spaceGUID, true, true)
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

//ListIsolationSegments : Returns all isolation segments
func (m *DefaultManager) ListIsolationSegments() ([]cfclient.IsolationSegment, error) {
	isolationSegments, err := m.Client.ListIsolationSegments()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total isolation segments returned :", len(isolationSegments))
	return isolationSegments, nil
}

func (m *DefaultManager) AddUserToOrgRole(userName, role, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org GUID %s", userName, role, orgGUID)
		return nil
	}
	roleToLower := strings.ToLower(role)
	if roleToLower == "auditors" {
		_, err := m.Client.AssociateOrgAuditorByUsername(orgGUID, userName)
		return err
	} else if roleToLower == "billing_managers" {
		_, err := m.Client.AssociateOrgManagerByUsername(orgGUID, userName)
		return err
	} else if roleToLower == "managers" {
		_, err := m.Client.AssociateOrgManagerByUsername(orgGUID, userName)
		return err
	} else {
		return fmt.Errorf("Role %s is not valid", role)
	}
	return nil
}

func (m *DefaultManager) ListAllOrgQuotas() (map[string]string, error) {
	quotas := make(map[string]string)
	orgQutotas, err := m.Client.ListOrgQuotas()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total org quotas returned :", len(orgQutotas))
	for _, quota := range orgQutotas {
		quotas[quota.Name] = quota.Guid
	}
	return quotas, nil
}

func (m *DefaultManager) CreateQuota(quota QuotaEntity) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create quota %+v", quota)
		return "dry-run-quota-guid", nil
	}

	orgQuota, err := m.Client.CreateOrgQuota(cfclient.OrgQuotaRequest{
		Name: quota.GetName(),
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		TotalPrivateDomains:     quota.TotalPrivateDomains,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	if err != nil {
		return "", err
	}
	return orgQuota.Guid, nil
}

func (m *DefaultManager) UpdateQuota(quotaGUID string, quota QuotaEntity) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %+v with GUID %s", quota, quotaGUID)
		return nil
	}
	_, err := m.Client.UpdateOrgQuota(quotaGUID, cfclient.OrgQuotaRequest{
		Name: quota.GetName(),
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		TotalPrivateDomains:     quota.TotalPrivateDomains,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	return err
}

func (m *DefaultManager) AssignQuotaToOrg(orgGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota GUID %s to org GUID %s", quotaGUID, orgGUID)
		return nil
	}
	org, err := m.Client.GetOrgByGuid(orgGUID)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateOrg(orgGUID, cfclient.OrgRequest{
		Name:                org.Name,
		QuotaDefinitionGuid: quotaGUID,
	})
	return err
}

//GetCFUsers Returns a list of space users who has a given role
func (m *DefaultManager) GetCFUsers(entityGUID, entityType, role string) (map[string]string, error) {
	userMap := make(map[string]string)
	roleLowerCase := strings.ToLower(role)
	entityTypeLowerCase := strings.ToLower(entityType)

	var users []cfclient.User
	var err error
	if entityTypeLowerCase == "spaces" {
		if roleLowerCase == "auditors" {
			users, err = m.Client.ListSpaceAuditors(entityGUID)
			if err != nil {
				return nil, err
			}
		} else if roleLowerCase == "developers" {
			users, err = m.Client.ListSpaceDevelopers(entityGUID)
			if err != nil {
				return nil, err
			}
		} else if roleLowerCase == "managers" {
			users, err = m.Client.ListSpaceManagers(entityGUID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Role %s is not valid", role)
		}
	} else if entityTypeLowerCase == "organizations" {
		if roleLowerCase == "auditors" {
			users, err = m.Client.ListOrgAuditors(entityGUID)
			if err != nil {
				return nil, err
			}
		} else if roleLowerCase == "billing_managers" {
			users, err = m.Client.ListOrgBillingManagers(entityGUID)
			if err != nil {
				return nil, err
			}
		} else if roleLowerCase == "managers" {
			users, err = m.Client.ListOrgManagers(entityGUID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Role %s is not valid", role)
		}
	} else {
		return nil, fmt.Errorf("EntityType %s is not valid", entityType)
	}

	lo.G.Debug("Total users returned :", len(users))

	for _, user := range users {
		userMap[strings.ToLower(user.Username)] = user.Guid
	}
	return userMap, nil
}

//RemoveCFUser - Un assigns a given from the given user for a given org and space
func (m *DefaultManager) RemoveCFUserByUserName(entityGUID, entityType, userName, role string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s for type %s with role %s", userName, entityGUID, entityType, role)
		return nil
	}
	roleLowerCase := strings.ToLower(role)
	entityTypeLowerCase := strings.ToLower(entityType)

	if entityTypeLowerCase == "spaces" {
		if roleLowerCase == "auditors" {
			return m.Client.RemoveSpaceAuditorByUsername(entityGUID, userName)
		} else if roleLowerCase == "developers" {
			return m.Client.RemoveSpaceDeveloperByUsername(entityGUID, userName)
		} else if roleLowerCase == "managers" {
			return m.Client.RemoveSpaceManagerByUsername(entityGUID, userName)
		} else {
			return fmt.Errorf("Role %s is not valid", role)
		}
	} else if entityTypeLowerCase == "organizations" {
		if roleLowerCase == "auditors" {
			return m.Client.RemoveOrgAuditorByUsername(entityGUID, userName)
		} else if roleLowerCase == "billing_managers" {
			return m.Client.RemoveOrgBillingManagerByUsername(entityGUID, userName)
		} else if roleLowerCase == "managers" {
			return m.Client.RemoveOrgManagerByUsername(entityGUID, userName)
		} else {
			return fmt.Errorf("Role %s is not valid", role)
		}
	} else {
		return fmt.Errorf("EntityType %s is not valid", entityType)
	}
	return nil
}

func (m *DefaultManager) OrgQuotaByName(name string) (cfclient.OrgQuota, error) {
	return m.Client.GetOrgQuotaByName(name)
}
func (m *DefaultManager) SpaceQuotaByName(name string) (cfclient.SpaceQuota, error) {
	return m.Client.GetSpaceQuotaByName(name)
}

func (m *DefaultManager) ListAllPrivateDomains() (map[string]PrivateDomainInfo, error) {
	domains, err := m.Client.ListDomains()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(domains))
	privateDomainMap := make(map[string]PrivateDomainInfo)
	for _, privateDomain := range domains {
		privateDomainMap[privateDomain.Name] = PrivateDomainInfo{
			OrgGUID:           privateDomain.OwningOrganizationGuid,
			PrivateDomainGUID: privateDomain.Guid,
		}
	}
	return privateDomainMap, nil
}

func (m *DefaultManager) ListOrgOwnedPrivateDomains(orgGUID string) (map[string]string, error) {
	orgOwnedPrivateDomainMap := make(map[string]string)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID == privateDomain.OwningOrganizationGuid {
			orgOwnedPrivateDomainMap[privateDomain.Name] = privateDomain.Guid
		}
	}
	return orgOwnedPrivateDomainMap, nil
}

func (m *DefaultManager) ListOrgSharedPrivateDomains(orgGUID string) (map[string]string, error) {
	orgSharedPrivateDomainMap := make(map[string]string)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID != privateDomain.OwningOrganizationGuid {
			orgSharedPrivateDomainMap[privateDomain.Name] = privateDomain.Guid
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

func (m *DefaultManager) DeletePrivateDomain(guid string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Delete private domain %s", guid)
		return nil
	}
	return m.Client.DeleteDomain(guid)
}
func (m *DefaultManager) CreatePrivateDomain(orgGUID, privateDomain string) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create private domain %s for org GUID %s", privateDomain, orgGUID)
		return "dry-run-private-domain-guid", nil
	}
	domain, err := m.Client.CreateDomain(privateDomain, orgGUID)
	if err != nil {
		return "", err
	}
	return domain.Guid, nil
}
func (m *DefaultManager) SharePrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	_, err := m.Client.ShareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
	return err
}
func (m *DefaultManager) RemoveSharedPrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: remove share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	return m.Client.UnshareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
}
