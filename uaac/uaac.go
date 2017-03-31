package uaac

import (
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
		fmt.Println(fmt.Sprintf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID))
		return nil
	} else {
		url := fmt.Sprintf("%s/Users", m.Host)
		payload := fmt.Sprintf(`{"userName":"%s","emails":[{"value":"%s"}],"origin":"%s","externalId":"%s"}`, userName, userEmail, origin, strings.Replace(externalID, "\\,", ",", 1))
		if _, err := http.NewManager().Post(url, m.UUACToken, payload); err != nil {
			return err
		}
		fmt.Println("successfully added user", userName)
		return nil
	}
}

//ListUsers -
func (m *DefaultUAACManager) ListUsers() (map[string]string, error) {
	lo.G.Info("Getting UAAC users from Cloud Foundry")
	users := make(map[string]string)
	url := fmt.Sprintf("%s/Users?count=5000", m.Host)
	userList := new(UserList)
	if err := http.NewManager().Get(url, m.UUACToken, userList); err != nil {
		return nil, err
	}
	for _, user := range userList.Users {
		users[strings.ToLower(user.Name)] = user.ID
	}
	return users, nil
}

//DefaultUAACManager -
type DefaultUAACManager struct {
	Host      string
	UUACToken string
}
