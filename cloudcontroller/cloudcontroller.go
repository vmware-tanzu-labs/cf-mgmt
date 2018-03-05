package cloudcontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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

func (m *DefaultManager) CreateSpace(spaceName, orgGUID string) error {
	_, err := m.Client.CreateSpace(cfclient.SpaceRequest{
		Name:             spaceName,
		OrganizationGuid: orgGUID,
	})
	return err
}

func (m *DefaultManager) ListSpaces(orgGUID string) ([]cfclient.Space, error) {
	spaces, err := m.Client.ListSpacesByQuery(url.Values{
		"organization_guid": []string{orgGUID},
	})
	if err != nil {
		return nil, err
	}
	return spaces, err

}

func (m *DefaultManager) AddUserToSpaceRole(userName, role, spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, role, spaceGUID)
		return nil
	}

	roleLowerCase := strings.ToLower(role)
	if roleLowerCase == "auditors" {
		_, err := m.Client.AssociateSpaceAuditorByUsername(spaceGUID, userName)
		return err
	} else if roleLowerCase == "developers" {
		_, err := m.Client.AssociateSpaceDeveloperByUsername(spaceGUID, userName)
		return err
	} else if roleLowerCase == "managers" {
		_, err := m.Client.AssociateSpaceManagerByUsername(spaceGUID, userName)
		return err
	} else {
		return fmt.Errorf("Role %s is not valid", role)
	}
	return nil
}

func (m *DefaultManager) AddUserToOrg(userName, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to orgGUID %s", userName, orgGUID)
		return nil
	}
	_, err := m.Client.AssociateOrgUserByUsername(orgGUID, userName)
	return err
}

func (m *DefaultManager) UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: setting sshAllowed to %v for spaceGUID %s", sshAllowed, spaceGUID)
		return nil
	}
	space, err := m.Client.GetSpaceByGuid(spaceGUID)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateSpace(spaceGUID, cfclient.SpaceRequest{
		Name:             space.Name,
		AllowSSH:         sshAllowed,
		OrganizationGuid: space.OrganizationGuid,
	})
	return err
}

func (m *DefaultManager) ListNonDefaultSecurityGroups() (map[string]SecurityGroupInfo, error) {
	securityGroups := make(map[string]SecurityGroupInfo)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.DefaultRunning == false && groupMap.DefaultStaging == false {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}

func (m *DefaultManager) ListDefaultSecurityGroups() (map[string]SecurityGroupInfo, error) {
	securityGroups := make(map[string]SecurityGroupInfo)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.DefaultRunning == true || groupMap.DefaultStaging == true {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}
func (m *DefaultManager) ListSecurityGroups() (map[string]SecurityGroupInfo, error) {
	securityGroups := make(map[string]SecurityGroupInfo)
	secGroups, err := m.Client.ListSecGroups()
	if err != nil {
		return securityGroups, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	for _, sg := range secGroups {
		bytes, _ := json.Marshal(sg.Rules)
		securityGroups[sg.Name] = SecurityGroupInfo{
			GUID:           sg.Guid,
			Rules:          string(bytes),
			DefaultRunning: sg.Running,
			DefaultStaging: sg.Staging,
		}
	}
	return securityGroups, nil
}

//GetSecurityGroupRules - returns a array of rules based on sgGUID
func (m *DefaultManager) GetSecurityGroupRules(sgGUID string) ([]byte, error) {
	secGroup, err := m.Client.GetSecGroup(sgGUID)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(secGroup.Rules, "", "\t")
}

func (m *DefaultManager) UpdateSecurityGroup(sgGUID, sgName, contents string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: updating securityGroup %s with guid %s with contents %s", sgName, sgGUID, contents)
		return nil
	}
	securityGroup := &cfclient.SecGroup{}
	err := json.Unmarshal([]byte(contents), &securityGroup)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateSecGroup(sgGUID, sgName, securityGroup.Rules, nil)
	return err
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating securityGroup %s with contents %s", sgName, contents)
		return "dry-run-security-group-guid", nil
	}
	securityGroup := &cfclient.SecGroup{}
	err := json.Unmarshal([]byte(contents), &securityGroup)
	if err != nil {
		return "", err
	}
	securityGroup, err = m.Client.CreateSecGroup(sgName, securityGroup.Rules, nil)
	return securityGroup.Guid, err
}

func (m *DefaultManager) ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error) {
	secGroups, err := m.Client.ListSpaceSecGroups(spaceGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	names := make(map[string]string)
	for _, sg := range secGroups {
		if sg.Running == false && sg.Staging == false {
			names[sg.Name] = sg.Guid
		}
	}
	return names, nil
}

func (m *DefaultManager) AssignRunningSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s as running security group", sgGUID)
		return nil
	}
	return m.Client.BindRunningSecGroup(sgGUID)
}
func (m *DefaultManager) AssignStagingSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s as staging security group", sgGUID)
		return nil
	}
	return m.Client.BindStagingSecGroup(sgGUID)
}
func (m *DefaultManager) UnassignRunningSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassinging sgGUID %s as running security group", sgGUID)
		return nil
	}
	return m.Client.UnbindRunningSecGroup(sgGUID)
}
func (m *DefaultManager) UnassignStagingSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassigning sgGUID %s as staging security group", sgGUID)
		return nil
	}
	return m.Client.UnbindStagingSecGroup(sgGUID)
}

func (m *DefaultManager) AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s to spaceGUID %s", sgGUID, spaceGUID)
		return nil
	}
	return m.Client.BindSecGroup(sgGUID, spaceGUID)
}

func (m *DefaultManager) AssignQuotaToSpace(spaceGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quotaGUID %s to spaceGUID %s", quotaGUID, spaceGUID)
		return nil
	}
	return m.Client.AssignSpaceQuota(quotaGUID, spaceGUID)
}

func (m *DefaultManager) CreateSpaceQuota(quota SpaceQuotaEntity) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %+v", quota)
		return "dry-run-space-quota-guid", nil
	}
	spaceQuota, err := m.Client.CreateSpaceQuota(cfclient.SpaceQuotaRequest{
		Name:                    quota.GetName(),
		OrganizationGuid:        quota.OrgGUID,
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	if err != nil {
		return "", err
	}
	return spaceQuota.Guid, nil
}

func (m *DefaultManager) UpdateSpaceQuota(quotaGUID string, quota SpaceQuotaEntity) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %s with %+v", quotaGUID, quota)
		return nil
	}
	_, err := m.Client.UpdateSpaceQuota(quotaGUID, cfclient.SpaceQuotaRequest{
		Name:                    quota.GetName(),
		OrganizationGuid:        quota.OrgGUID,
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	return err
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	spaceQuotas, err := m.Client.ListOrgSpaceQuotas(orgGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	for _, quota := range spaceQuotas {
		quotas[quota.Name] = quota.Guid
	}
	return quotas, nil
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
