package cloudcontroller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/http"
	"github.com/xchapter7x/lo"
)

func NewManager(host, token string) Manager {
	return &DefaultManager{
		Host:  host,
		Token: token,
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
	lo.G.Info("Total spaces returned :", len(spaceResources.Spaces))
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
	url := fmt.Sprintf("%s/v2/spaces/%s/%s", m.Host, spaceGUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err := m.HTTP.Put(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) AddUserToOrg(userName, orgGUID string) error {
	url := fmt.Sprintf("%s/v2/organizations/%s/users", m.Host, orgGUID)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	err := m.HTTP.Put(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error {
	url := fmt.Sprintf("%s/v2/spaces/%s", m.Host, spaceGUID)
	sendString := fmt.Sprintf(`{"allow_ssh":%t}`, sshAllowed)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListSecurityGroups() (map[string]string, error) {
	securityGroups := make(map[string]string)
	url := fmt.Sprintf("%s/v2/security_groups", m.Host)
	sgResources := &SecurityGroupResources{}
	err := m.listResources(url, sgResources, NewSecurityGroupResources)
	if err != nil {
		return nil, err
	}
	lo.G.Info("Total security groups returned :", len(sgResources.SecurityGroups))
	for _, sg := range sgResources.SecurityGroups {
		securityGroups[sg.Entity.Name] = sg.MetaData.GUID
	}
	return securityGroups, nil
}

func (m *DefaultManager) UpdateSecurityGroup(sgGUID, sgName, contents string) error {
	url := fmt.Sprintf("%s/v2/security_groups/%s", m.Host, sgGUID)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (string, error) {
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

func (m *DefaultManager) AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error {
	url := fmt.Sprintf("%s/v2/security_groups/%s/spaces/%s", m.Host, sgGUID, spaceGUID)
	err := m.HTTP.Put(url, m.Token, "")
	return err
}

func (m *DefaultManager) AssignQuotaToSpace(spaceGUID, quotaGUID string) error {
	url := fmt.Sprintf("%s/v2/space_quota_definitions/%s/spaces/%s", m.Host, quotaGUID, spaceGUID)
	err := m.HTTP.Put(url, m.Token, "")
	return err
}

func (m *DefaultManager) CreateSpaceQuota(orgGUID, quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) (string, error) {
	url := fmt.Sprintf("%s/v2/space_quota_definitions", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t,"organization_guid":"%s"}`, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed, orgGUID)
	body, err := m.HTTP.Post(url, m.Token, sendString)
	if err != nil {
		return "", err
	}
	quotaResource := &Quota{}
	if err = json.Unmarshal([]byte(body), &quotaResource); err != nil {
		return "", err
	}
	return quotaResource.MetaData.GUID, nil
}

func (m *DefaultManager) UpdateSpaceQuota(orgGUID, quotaGUID, quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) error {
	url := fmt.Sprintf("%s/v2/space_quota_definitions/%s", m.Host, quotaGUID)
	sendString := fmt.Sprintf(`{"guid":"%s","name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t,"organization_guid":"%s"}`, quotaGUID, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed, orgGUID)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	url := fmt.Sprintf("%s/v2/organizations/%s/space_quota_definitions", m.Host, orgGUID)
	quotaResources := &Quotas{}
	err := m.listResources(url, quotaResources, NewQuotasResources)
	if err != nil {
		return nil, err
	}
	lo.G.Info("Total space quotas returned :", len(quotaResources.Quotas))
	for _, quota := range quotaResources.Quotas {
		quotas[quota.Entity.Name] = quota.MetaData.GUID
	}
	return quotas, nil
}

func (m *DefaultManager) CreateOrg(orgName string) error {
	url := fmt.Sprintf("%s/v2/organizations", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s"}`, orgName)
	_, err := m.HTTP.Post(url, m.Token, sendString)
	return err
}

func (m *DefaultManager) DeleteOrg(orgName string) error {
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			url := fmt.Sprintf("%s/v2/organizations/%s?recursive=true", m.Host, org.MetaData.GUID)
			err = m.HTTP.Delete(url, m.Token)
		}
	}

	return err
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultManager) ListOrgs() ([]*Org, error) {
	url := fmt.Sprintf("%s/v2/organizations?results-per-page=100", m.Host)
	orgs := &Orgs{}
	err := m.listResources(url, orgs, NewOrgResources)
	if err != nil {
		return nil, err
	}
	lo.G.Info("Total orgs returned :", len(orgs.Orgs))
	return orgs.Orgs, nil
}

func (m *DefaultManager) AddUserToOrgRole(userName, role, orgGUID string) error {
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
	lo.G.Info("Total org quotas returned :", len(quotaResources.Quotas))
	for _, quota := range quotaResources.Quotas {
		quotas[quota.Entity.Name] = quota.MetaData.GUID
	}
	return quotas, nil
}

func (m *DefaultManager) CreateQuota(quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) (string, error) {
	url := fmt.Sprintf("%s/v2/quota_definitions", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t}`, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed)
	body, err := m.HTTP.Post(url, m.Token, sendString)
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

func (m *DefaultManager) UpdateQuota(quotaGUID, quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) error {

	url := fmt.Sprintf("%s/v2/quota_definitions/%s", m.Host, quotaGUID)
	sendString := fmt.Sprintf(`{"guid":"%s","name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t}`, quotaGUID, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) AssignQuotaToOrg(orgGUID, quotaGUID string) error {
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
	lo.G.Info("Total users returned :", len(users.Users))

	for _, user := range users.Users {
		userMap[strings.ToLower(user.Entity.UserName)] = user.MetaData.GUID
	}
	return userMap, nil
}

//RemoveCFUser - Un assigns a given from the given user for a given org and space
func (m *DefaultManager) RemoveCFUser(entityGUID, entityType, userGUID, role string) error {
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
