package cloudcontroller

import (
	"encoding/json"
	"fmt"

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

func (m *DefaultManager) ListSpaces(orgGUID string) ([]Space, error) {
	spaceResources := &SpaceResources{}
	url := fmt.Sprintf("%s/v2/organizations/%s/spaces?inline-relations-depth=1", m.Host, orgGUID)
	if err := m.HTTP.Get(url, m.Token, spaceResources); err == nil {
		spaces := spaceResources.Spaces
		return spaces, nil
	} else {
		return nil, err
	}
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
	var err error
	securityGroups := make(map[string]string)
	url := fmt.Sprintf("%s/v2/security_groups", m.Host)
	sgResources := &SecurityGroupResources{}
	if err = m.HTTP.Get(url, m.Token, sgResources); err == nil {
		for _, sg := range sgResources.SecurityGroups {
			securityGroups[sg.Entity.Name] = sg.MetaData.GUID
		}
		return securityGroups, nil
	}
	return nil, err
}

func (m *DefaultManager) UpdateSecurityGroup(sgGUID, sgName, contents string) error {
	url := fmt.Sprintf("%s/v2/security_groups/%s", m.Host, sgGUID)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (string, error) {
	url := fmt.Sprintf("%s/v2/security_groups", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s","rules":%s}`, sgName, contents)
	if body, err := m.HTTP.Post(url, m.Token, sendString); err != nil {
		return "", err
	} else {
		sgResource := &SecurityGroup{}
		if err := json.Unmarshal([]byte(body), &sgResource); err == nil {
			return sgResource.MetaData.GUID, nil
		} else {
			return "", err
		}
	}
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
	if body, err := m.HTTP.Post(url, m.Token, sendString); err == nil {
		quotaResource := &Quota{}
		if err = json.Unmarshal([]byte(body), &quotaResource); err == nil {
			return quotaResource.MetaData.GUID, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (m *DefaultManager) UpdateSpaceQuota(orgGUID, quotaGUID, quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) error {
	url := fmt.Sprintf("%s/v2/space_quota_definitions/%s", m.Host, quotaGUID)
	sendString := fmt.Sprintf(`{"guid":"%s","name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t,"organization_guid":"%s"}`, quotaGUID, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed, orgGUID)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListSpaceQuotas(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	url := fmt.Sprintf("%s/v2/organizations/%s/space_quota_definitions", m.Host, orgGUID)
	quotaResources := &Quotas{}
	if err := m.HTTP.Get(url, m.Token, quotaResources); err == nil {
		for _, quota := range quotaResources.Quotas {
			quotas[quota.Entity.Name] = quota.MetaData.GUID
		}
		return quotas, nil
	} else {
		return nil, err
	}
}

func (m *DefaultManager) CreateOrg(orgName string) error {
	url := fmt.Sprintf("%s/v2/organizations", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s"}`, orgName)
	_, err := m.HTTP.Post(url, m.Token, sendString)
	return err
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultManager) ListOrgs() ([]*Org, error) {
	url := fmt.Sprintf("%s/v2/organizations?results-per-page=100", m.Host)
	orgs := &Orgs{}
	var err = m.HTTP.Get(url, m.Token, orgs)
	if err != nil {
		return nil, err
	}
	if orgs.NextURL == "" {
		return orgs.Orgs, nil
	}
	nextURL := orgs.NextURL
	orgsTemp := &Orgs{}
	for nextURL != "" {
		url = fmt.Sprintf("%s%s", m.Host, nextURL)
		lo.G.Info("getOrgs() URL :", url)
		err = m.HTTP.Get(url, m.Token, orgsTemp)
		if err != nil {
			return nil, err
		}
		orgs.Orgs = append(orgs.Orgs, orgsTemp.Orgs...)
		nextURL = orgsTemp.NextURL
	}
	lo.G.Info("Total orgs returned :", len(orgs.Orgs))
	return orgs.Orgs, nil
}

func (m *DefaultManager) AddUserToOrgRole(userName, role, orgGUID string) error {
	url := fmt.Sprintf("%s/v2/organizations/%s/%s", m.Host, orgGUID, role)
	sendString := fmt.Sprintf(`{"username": "%s"}`, userName)
	return m.HTTP.Put(url, m.Token, sendString)
}

func (m *DefaultManager) ListQuotas() (map[string]string, error) {
	quotas := make(map[string]string)
	url := fmt.Sprintf("%s/v2/quota_definitions", m.Host)
	quotaResources := &Quotas{}
	if err := m.HTTP.Get(url, m.Token, quotaResources); err == nil {
		for _, quota := range quotaResources.Quotas {
			quotas[quota.Entity.Name] = quota.MetaData.GUID
		}
		return quotas, nil
	} else {
		return nil, err
	}
}

func (m *DefaultManager) CreateQuota(quotaName string,
	memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
	paidServicePlansAllowed bool) (string, error) {
	url := fmt.Sprintf("%s/v2/quota_definitions", m.Host)
	sendString := fmt.Sprintf(`{"name":"%s","memory_limit":%d,"instance_memory_limit":%d,"total_routes":%d,"total_services":%d,"non_basic_services_allowed":%t}`, quotaName, memoryLimit, instanceMemoryLimit, totalRoutes, totalServices, paidServicePlansAllowed)
	if body, err := m.HTTP.Post(url, m.Token, sendString); err == nil {
		quotaResource := &Quota{}
		if err = json.Unmarshal([]byte(body), &quotaResource); err == nil {
			return quotaResource.MetaData.GUID, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
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
	var err = m.HTTP.Get(url, m.Token, users)
	if err != nil {
		return nil, err
	}
	nextURL := users.NextURL
	usersTemp := &OrgSpaceUsers{}
	for nextURL != "" {
		url = fmt.Sprintf("%s%s", m.Host, nextURL)
		err = m.HTTP.Get(url, m.Token, usersTemp)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, usersTemp.Users...)
		nextURL = usersTemp.NextURL
	}
	lo.G.Info(fmt.Sprintf("Total %d users with %s role returned for %s  %s", len(users.Users), role, entityType, entityGUID))
	for _, user := range users.Users {
		userMap[user.Entity.UserName] = user.MetaData.GUID
	}
	return userMap, nil
}

//RemoveCFUser - Un assigns a given from the given user for a given org and space
func (m *DefaultManager) RemoveCFUser(entityGUID, entityType, userGUID, role string) error {
	url := fmt.Sprintf("%s/v2/%s/%s/%s/%s", m.Host, entityType, entityGUID, role, userGUID)
	return m.HTTP.Delete(url, m.Token)
}
