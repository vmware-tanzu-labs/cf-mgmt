package uaa

import (
	"fmt"
	"strings"

	"github.com/xchapter7x/lo"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
)

//go:generate counterfeiter -o fakes/uaa_client.go uaa.go uaa
type uaa interface {
	CreateUser(user uaaclient.User) (*uaaclient.User, error)
	ListAllUsers(filter string, sortBy string, attributes string, sortOrder uaaclient.SortOrder) ([]uaaclient.User, error)
}

//Manager -
type Manager interface {
	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (map[string]*uaaclient.User, error)
	CreateExternalUser(userName, userEmail, externalID, origin string) (err error)
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	Peek   bool
	Client uaa
}

//NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, clientID, clientSecret string, peek bool) (Manager, error) {
	target := fmt.Sprintf("https://uaa.%s", sysDomain)
	client, err := uaaclient.NewWithClientCredentials(target, "", clientID, clientSecret, uaaclient.OpaqueToken, true)
	if err != nil {
		return nil, err
	}
	return &DefaultUAAManager{
		Client: client,
		Peek:   peek,
	}, nil
}

//CreateExternalUser -
func (m *DefaultUAAManager) CreateExternalUser(userName, userEmail, externalID, origin string) error {
	if userName == "" || userEmail == "" || externalID == "" {
		return fmt.Errorf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID)
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: successfully added user [%s]", userName)
		return nil
	}

	m.Client.CreateUser(uaaclient.User{
		Username:   userName,
		ExternalID: externalID,
		Origin:     origin,
		Emails: []uaaclient.Email{
			uaaclient.Email{
				Value: userEmail,
			},
		},
	})
	lo.G.Infof("successfully added user [%s]", userName)
	return nil
}

//ListUsers - Returns a map containing username as key and user guid as value
func (m *DefaultUAAManager) ListUsers() (map[string]*uaaclient.User, error) {
	userMap := make(map[string]*uaaclient.User)
	lo.G.Debug("Getting users from Cloud Foundry")
	users, err := m.Client.ListAllUsers("", "", "", "")
	if err != nil {
		return nil, err
	}
	lo.G.Debugf("Found %d users in the CF instance", len(users))
	for _, user := range users {
		userMap[strings.ToLower(user.Username)] = &user
		if user.ExternalID != "" {
			userMap[strings.ToLower(user.ExternalID)] = &user
		}
	}
	return userMap, nil
}
