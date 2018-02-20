package cloudcontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/http"
	"github.com/xchapter7x/lo"
)

func NewManager(host, token string, peek bool) Manager {
	return &DefaultManager{
		Host:  host,
		Token: token,
		Peek:  peek,
		HTTP:  http.NewManager(),
	}
}

func (m *DefaultManager) CreateSpace(spaceName, orgGUID string) error {
	url := fmt.Sprintf("%s/v2/spaces", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s", "organization_guid":"%s"}`, spaceName, orgGUID)
	_, err := m.HTTP.Post(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) ListSpaces(orgGUID string) ([]*Space, error) {
	spaceResources := &SpaceResources{}
	url := fmt.Sprintf("%s/v2/organizations/%s/spaces", m.Host, orgGUID)
	err := m.listResources(url, spaceResources, NewSpaceResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total spaces returned :", len(spaceResources.Spaces))
	return spaceResources.Spaces, nil

}

func (m *DefaultManager) listResources(url string, target Pagination, createInstance func() Pagination) error {
	var err = m.HTTP.Get(url, m.Token, target)
	if err != nil {
		return err
	}
	if target.GetNextURL() == "" {
		return nil
	}
	nextURL := target.GetNextURL()
	for nextURL != "" {
		lo.G.Debugf("NextURL: %s", nextURL)
		tempTarget := createInstance()
		url = fmt.Sprintf("%s%s", m.Host, nextURL)
		err = m.HTTP.Get(url, m.Token, tempTarget)
		if err != nil {
			return err
		}
		target.AddInstances(tempTarget)
		nextURL = tempTarget.GetNextURL()
	}
	return nil
}

func (m *DefaultManager) AddUserToSpaceRole(userName, role, spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, role, spaceGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/spaces/%s/%s", m.Host, spaceGUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err := m.HTTP.Put(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) AddUserToOrg(userName, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to orgGUID %s", userName, orgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s/users", m.Host, orgGUID)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err := m.HTTP.Put(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: setting sshAllowed to %v for spaceGUID %s", sshAllowed, spaceGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/spaces/%s", m.Host, spaceGUID)
	sendString := fmt.Sprintf(`{"allow_ssh":%t}`, sshAllowed)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListNonDefaultSecurityGroups() (map[string]SecurityGroupInfo, error) {
	securityGroups := make(map[string]SecurityGroupInfo)
	groupMap, err := m.listSecurityGroups()
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

func (m *DefaultManager) listSecurityGroups() (map[string]SecurityGroupInfo, error) {
	securityGroups := make(map[string]SecurityGroupInfo)
	url := fmt.Sprintf("%s/v2/security_groups", m.Host)
	sgResources := &SecurityGroupResources{}
	err := m.listResources(url, sgResources, NewSecurityGroupResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total security groups returned :", len(sgResources.SecurityGroups))
	for _, sg := range sgResources.SecurityGroups {
		bytes, _ := json.Marshal(sg.Entity.Rules)
		securityGroups[sg.Entity.Name] = SecurityGroupInfo{
			GUID:           sg.MetaData.GUID,
			Rules:          string(bytes),
			DefaultRunning: sg.Entity.DefaultRunning,
			DefaultStaging: sg.Entity.DefaultStaging,
		}
	}
	return securityGroups, nil
}

//GetSecurityGroupRules - returns a array of rules based on sgGUID
func (m *DefaultManager) GetSecurityGroupRules(sgGUID string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/security_groups/%s", m.Host, sgGUID)
	sgRule := &SecurityGroupRule{}
	err := m.HTTP.Get(url, m.Token, sgRule)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(sgRule.Entity.Rules, "", "\t")
}

func (m *DefaultManager) UpdateSecurityGroup(sgGUID, sgName, contents string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: updating securityGroup %s with guid %s with contents %s", sgName, sgGUID, contents)
		return nil
	}
	url := fmt.Sprintf("%s/v2/security_groups/%s", m.Host, sgGUID)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating securityGroup %s with contents %s", sgName, contents)
		return "dry-run-security-group-guid", nil
	}
	url := fmt.Sprintf("%s/v2/security_groups", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	body, err := m.HTTP.Post(url, m.Token, sendString)
	if err != nil {
		return "", err
	}
	sgResource := &SecurityGroup{}
	err = json.Unmarshal([]byte(body), &sgResource)
	if err != nil {
		return "", err
	}
	return sgResource.MetaData.GUID, nil
}

func (m *DefaultManager) ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v2/spaces/%s/security_groups", m.Host, spaceGUID)
	sgResources := &SecurityGroupResources{}
	err := m.listResources(url, sgResources, NewSecurityGroupResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total security groups returned :", len(sgResources.SecurityGroups))
	names := make(map[string]string)
	for _, sg := range sgResources.SecurityGroups {
		if sg.Entity.DefaultRunning == false && sg.Entity.DefaultStaging == false {
			names[sg.Entity.Name] = sg.MetaData.GUID
		}
	}
	return names, nil
}

func (m *DefaultManager) AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s to spaceGUID %s", sgGUID, spaceGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/security_groups/%s/spaces/%s", m.Host, sgGUID, spaceGUID)
	err := m.HTTP.Put(url, m.Token, "")
	return err
}

func (m *DefaultManager) AssignQuotaToSpace(spaceGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quotaGUID %s to spaceGUID %s", quotaGUID, spaceGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/space_quota_definitions/%s/spaces/%s", m.Host, quotaGUID, spaceGUID)
	err := m.HTTP.Put(url, m.Token, "")
	return err
}

func (m *DefaultManager) CreateSpaceQuota(quota SpaceQuotaEntity) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %+v", quota)
		return "dry-run-space-quota-guid", nil
	}
	url := fmt.Sprintf("%s/v2/space_quota_definitions", m.Host)
	sendString, err := json.Marshal(quota)
	if err != nil {
		return "", err
	}

	body, err := m.HTTP.Post(url, m.Token, string(sendString))
	if err != nil {
		return "", err
	}
	quotaResource := &Quota{}
	if err = json.Unmarshal([]byte(body), &quotaResource); err != nil {
		return "", err
	}
	return quotaResource.MetaData.GUID, nil
}

func (m *DefaultManager) UpdateSpaceQuota(quotaGUID string, quota SpaceQuotaEntity) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %s with %+v", quotaGUID, quota)
		return nil
	}
	url := fmt.Sprintf("%s/v2/space_quota_definitions/%s", m.Host, quotaGUID)
	sendString, err := json.Marshal(quota)
	if err != nil {
		return err
	}
	return m.HTTP.Put(url, m.Token, string(sendString))
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	url := fmt.Sprintf("%s/v2/organizations/%s/space_quota_definitions", m.Host, orgGUID)
	quotaResources := &Quotas{}
	err := m.listResources(url, quotaResources, NewQuotasResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(quotaResources.Quotas))
	for _, quota := range quotaResources.Quotas {
		quotas[quota.Entity.Name] = quota.MetaData.GUID
	}
	return quotas, nil
}

func (m *DefaultManager) CreateOrg(orgName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org %s", orgName)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s"}`, orgName)
	_, err := m.HTTP.Post(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) DeleteOrg(orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org with GUID %s", orgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s?recursive=true", m.Host, orgGUID)
	return m.HTTP.Delete(url, m.Token)
}

func (m *DefaultManager) DeleteOrgByName(orgName string) error {
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			if m.Peek {
				lo.G.Infof("[dry-run]: delete org %s", orgName)
				return nil
			}
			url := fmt.Sprintf("%s/v2/organizations/%s?recursive=true", m.Host, org.MetaData.GUID)
			return m.HTTP.Delete(url, m.Token)
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
	url := fmt.Sprintf("%s/v2/spaces/%s?recursive=true", m.Host, spaceGUID)
	return m.HTTP.Delete(url, m.Token)
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultManager) ListOrgs() ([]*Org, error) {
	url := fmt.Sprintf("%s/v2/organizations?results-per-page=100", m.Host)
	orgs := &Orgs{}
	err := m.listResources(url, orgs, NewOrgResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(orgs.Orgs))
	return orgs.Orgs, nil
}

func (m *DefaultManager) AddUserToOrgRole(userName, role, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org GUID %s", userName, role, orgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s/%s", m.Host, orgGUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListAllOrgQuotas() (map[string]string, error) {
	quotas := make(map[string]string)
	url := fmt.Sprintf("%s/v2/quota_definitions", m.Host)
	quotaResources := &Quotas{}
	err := m.listResources(url, quotaResources, NewQuotasResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total org quotas returned :", len(quotaResources.Quotas))
	for _, quota := range quotaResources.Quotas {
		quotas[quota.Entity.Name] = quota.MetaData.GUID
	}
	return quotas, nil
}

func (m *DefaultManager) CreateQuota(quota QuotaEntity) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create quota %+v", quota)
		return "dry-run-quota-guid", nil
	}
	url := fmt.Sprintf("%s/v2/quota_definitions", m.Host)
	sendString, err := json.Marshal(quota)
	if err != nil {
		return "", err
	}
	body, err := m.HTTP.Post(url, m.Token, string(sendString))
	if err != nil {
		return "", err
	}
	quotaResource := &Quota{}
	err = json.Unmarshal([]byte(body), &quotaResource)
	if err != nil {
		return "", err
	}
	return quotaResource.MetaData.GUID, nil
}

func (m *DefaultManager) UpdateQuota(quotaGUID string, quota QuotaEntity) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %+v with GUID %s", quota, quotaGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/quota_definitions/%s", m.Host, quotaGUID)
	sendString, err := json.Marshal(quota)
	if err != nil {
		return err
	}
	return m.HTTP.Put(url, m.Token, string(sendString))
}

func (m *DefaultManager) AssignQuotaToOrg(orgGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota GUID %s to org GUID %s", quotaGUID, orgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s", m.Host, orgGUID)
	sendString := fmt.Sprintf(`{"quota_definition_guid":"%s"}`, quotaGUID)
	return m.HTTP.Put(url, m.Token, sendString)
}

//GetCFUsers Returns a list of space users who has a given role
func (m *DefaultManager) GetCFUsers(entityGUID, entityType, role string) (map[string]string, error) {
	userMap := make(map[string]string)
	url := fmt.Sprintf("%s/v2/%s/%s/%s?results-per-page=100", m.Host, entityType, entityGUID, role)
	users := &OrgSpaceUsers{}
	err := m.listResources(url, users, NewOrgSpaceUsers)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total users returned :", len(users.Users))

	for _, user := range users.Users {
		userMap[strings.ToLower(user.Entity.UserName)] = user.MetaData.GUID
	}
	return userMap, nil
}

//RemoveCFUser - Un assigns a given from the given user for a given org and space
func (m *DefaultManager) RemoveCFUser(entityGUID, entityType, userGUID, role string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user GUID %s from GUID %s for type %s with role %s", userGUID, entityGUID, entityType, role)
		return nil
	}
	url := fmt.Sprintf("%s/v2/%s/%s/%s/%s", m.Host, entityType, entityGUID, role, userGUID)
	return m.HTTP.Delete(url, m.Token)
}

//QuotaDef Returns quota definition for a given Quota
func (m *DefaultManager) QuotaDef(quotaDefGUID string, entityType string) (*Quota, error) {
	var apiPath string
	if "organizations" == entityType {
		apiPath = "quota_definitions"
	} else {
		apiPath = "space_quota_definitions"
	}
	url := fmt.Sprintf("%s/v2/%s/%s", m.Host, apiPath, quotaDefGUID)
	var err error
	quotaResource := &Quota{}
	if err = m.HTTP.Get(url, m.Token, quotaResource); err == nil {
		lo.G.Debugf("Quota returned : %v", quotaResource.Entity)
		return quotaResource, nil
	}
	lo.G.Errorf("Error from quota API call : %v", err)
	return nil, err
}

func (m *DefaultManager) ListAllPrivateDomains() (map[string]PrivateDomainInfo, error) {
	privateDomainResources := &PrivateDomainResources{}
	url := fmt.Sprintf("%s/v2/private_domains", m.Host)
	err := m.listResources(url, privateDomainResources, NewPrivateDomainResource)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(privateDomainResources.PrivateDomains))
	privateDomainMap := make(map[string]PrivateDomainInfo)
	for _, privateDomain := range privateDomainResources.PrivateDomains {
		privateDomainMap[privateDomain.Entity.Name] = PrivateDomainInfo{
			OrgGUID:           privateDomain.Entity.OrgGUID,
			PrivateDomainGUID: privateDomain.MetaData.GUID,
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
		if orgGUID == privateDomain.Entity.OrgGUID {
			orgOwnedPrivateDomainMap[privateDomain.Entity.Name] = privateDomain.MetaData.GUID
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
		if orgGUID != privateDomain.Entity.OrgGUID {
			orgSharedPrivateDomainMap[privateDomain.Entity.Name] = privateDomain.MetaData.GUID
		}
	}
	return orgSharedPrivateDomainMap, nil
}

func (m *DefaultManager) listOrgPrivateDomains(orgGUID string) ([]*PrivateDomain, error) {
	privateDomainResources := &PrivateDomainResources{}
	url := fmt.Sprintf("%s/v2/organizations/%s/private_domains", m.Host, orgGUID)
	err := m.listResources(url, privateDomainResources, NewPrivateDomainResource)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(privateDomainResources.PrivateDomains))
	return privateDomainResources.PrivateDomains, nil
}

func (m *DefaultManager) DeletePrivateDomain(guid string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Delete private domain %s", guid)
		return nil
	}
	url := fmt.Sprintf("%s/v2/private_domains/%s?async=false", m.Host, guid)
	return m.HTTP.Delete(url, m.Token)
}
func (m *DefaultManager) CreatePrivateDomain(orgGUID, privateDomain string) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create private domain %s for org GUID %s", privateDomain, orgGUID)
		return "dry-run-private-domain-guid", nil
	}
	url := fmt.Sprintf("%s/v2/private_domains", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s", "owning_organization_guid":"%s"}`, privateDomain, orgGUID)
	body, err := m.HTTP.Post(url, m.Token, sendString)
	if err != nil {
		return "", err
	}
	privateDomainResource := &PrivateDomain{}
	err = json.Unmarshal([]byte(body), &privateDomainResource)
	if err != nil {
		return "", err
	}
	return privateDomainResource.MetaData.GUID, nil
}
func (m *DefaultManager) SharePrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s/private_domains/%s", m.Host, sharedOrgGUID, privateDomainGUID)
	err := m.HTTP.Put(url, m.Token, "")
	return err
}
func (m *DefaultManager) RemoveSharedPrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: remove share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	url := fmt.Sprintf("%s/v2/organizations/%s/private_domains/%s", m.Host, sharedOrgGUID, privateDomainGUID)
	return m.HTTP.Delete(url, m.Token)
}
