package uaac

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/http"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(systemDomain, uuacToken string) (mgr Manager) {
	return &DefaultUAACManager{
		Host:      fmt.Sprintf("https://uaa.%s", systemDomain),
		UUACToken: uuacToken,
	}
}

//CreateExternalUser -
func (m *DefaultUAACManager) CreateExternalUser(userName, userEmail, externalID, origin string) error {
	if userName == "" || userEmail == "" || externalID == "" {
		msg := fmt.Sprintf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID)
		lo.G.Info(msg)
		return errors.New(msg)
	}
	url := fmt.Sprintf("%s/Users", m.Host)
	payload := fmt.Sprintf(`{"userName":"%s","emails":[{"value":"%s"}],"origin":"%s","externalId":"%s"}`, userName, userEmail, origin, strings.Replace(externalID, "\\,", ",", 1))
	if _, err := http.NewManager().Post(url, m.UUACToken, payload); err != nil {
		return err
	}
	lo.G.Info("successfully added user", userName)
	return nil
}

//ListUsers - Returns a map containing username as key and user guid as value
func (m *DefaultUAACManager) ListUsers() (map[string]string, error) {
	userIDMap := make(map[string]string)
	usersList, err := getUsers(m.Host, m.UUACToken)
	if err != nil {
		return nil, err
	}
	for _, user := range usersList.Users {
		userIDMap[strings.ToLower(user.UserName)] = user.ID
	}
	return userIDMap, nil
}

// UsersByID returns a map of Users keyed by ID.
func (m *DefaultUAACManager) UsersByID() (userIDMap map[string]User, err error) {
	userIDMap = make(map[string]User)
	userList, err := getUsers(m.Host, m.UUACToken)
	if err != nil {
		return nil, err
	}
	for _, user := range userList.Users {
		userIDMap[user.UserName] = user
	}
	return userIDMap, nil
}

//TODO Anwar - Make this API use pagination
func getUsers(host string, uaacToken string) (userList *UserList, err error) {
	lo.G.Info("Getting users from Cloud Foundry")
	url := fmt.Sprintf("%s/Users?count=5000", host)
	userList = new(UserList)
	if err := http.NewManager().Get(url, uaacToken, userList); err != nil {
		return nil, fmt.Errorf("couldn't retrieve users: %v", err)
	}
	lo.G.Infof("Found %d users in the CF instance", len(userList.Users))
	return userList, nil
}

//DefaultUAACManager -
type DefaultUAACManager struct {
	Host      string
	UUACToken string
}
