package azureAD

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

func NewManager(config *config.AzureADConfig) (*Manager, error) {
	token, err := getAccessToken(config)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Config:   config,
		Token:    token,
		userMap:  make(map[string]*UserType),
		groupMap: make(map[string][]string),
	}, nil
}

func (m *Manager) isGroupInCache(groupName string) bool {
	if m.groupMap == nil {
		return false
	}
	if _, ok := m.groupMap[groupName]; ok {
		return true
	}

	return false
}

func (m *Manager) addGroupToCache(groupName string, result []string) {
	if m.groupMap == nil {
		m.groupMap = make(map[string][]string)
	}
	m.groupMap[groupName] = result
}

func (m *Manager) GetADToken() string {
	return m.Token.AccessToken
}

func (m *Manager) GraphGetGroupMembers(token, groupName string) ([]string, error) {
	if m.isGroupInCache(groupName) {
		lo.G.Debugf("Group %s found in cache", groupName)
		return m.groupMap[groupName], nil
	}

	groupId, err := graphGetIdFromName(token, groupName)
	if err != nil {
		lo.G.Criticalf("Error converting group Name to group ID: %s", err)
		return nil, err
	}

	requestURL := GraphURL + "groups/" + groupId + "/members"

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token

	body, _ := doHttpCall(requestURL, headers)

	result := AADGroupMemberListType{}
	if err = json.Unmarshal(body, &result); err != nil {
		lo.G.Debugf("Reading body into result struct failed: %s", err)
		return nil, err
	}

	var userList []string
	for _, u := range result.Value {
		userList = append(userList, u.Upn)
	}

	m.addGroupToCache(groupName, userList)
	return userList, nil
}

func graphGetIdFromName(token, name string) (string, error) {
	// Some magic because the filter uses spaces and single quotes.
	// Spaces MUST be encode to %20 (not +, as is the default for go)
	spaceEncodedString := "$select=id&$filter=displayName eq '" + name + "'"
	t := &url.URL{Path: spaceEncodedString}
	requestURL := GraphURL + "groups?" + t.String()

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["ConsistencyLevel"] = "eventual"

	body, _ := doHttpCall(requestURL, headers)

	result := AADGroupType{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		lo.G.Errorf("Reading group Name body into result struct failed: %s", err)
		return "", err
	}

	if len(result.Value) != 1 {
		lo.G.Errorf("Number of Id's returned for groupname search should be exactly one!")
		return "", err
	}

	return result.Value[0].Id, nil
}

func doHttpCall(requestUrl string, Headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	lo.G.Debugf("Request URI: %v", requestUrl)
	if err != nil {
		lo.G.Criticalf("Something went wrong with http.NewRequest: %v ", err)
		return nil, err
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		lo.G.Criticalf("Something went wrong with the request: %v ", err)
		return nil, err
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return body, nil
}
