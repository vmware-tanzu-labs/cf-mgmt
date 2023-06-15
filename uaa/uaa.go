package uaa

import (
	"fmt"
	"strings"

	"github.com/xchapter7x/lo"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/uaa_client.go uaa.go uaa
type uaa interface {
	CreateUser(user uaaclient.User) (*uaaclient.User, error)
	ListAllUsers(filter string, sortBy string, attributes string, sortOrder uaaclient.SortOrder) ([]uaaclient.User, error)
}

// Manager -
type Manager interface {
	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (*Users, error)
	CreateExternalUser(userName, userEmail, externalID, origin string) (GUID string, err error)
}

// DefaultUAAManager -
type DefaultUAAManager struct {
	Peek   bool
	Client uaa
}

type User struct {
	Username   string
	ExternalID string
	Email      string
	Origin     string
	GUID       string
}

// NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, clientID, clientSecret, userAgent string, peek bool) (Manager, error) {
	target := fmt.Sprintf("https://uaa.%s", sysDomain)
	client, err := uaaclient.New(
		target,
		uaaclient.WithClientCredentials(clientID, clientSecret, uaaclient.OpaqueToken),
		uaaclient.WithUserAgent(userAgent),
		uaaclient.WithSkipSSLValidation(true),
	)
	if err != nil {
		return nil, err
	}

	return &DefaultUAAManager{
		Client: client,
		Peek:   peek,
	}, nil
}

// CreateExternalUser -
func (m *DefaultUAAManager) CreateExternalUser(userName, userEmail, externalID, origin string) (string, error) {
	if userName == "" || userEmail == "" || externalID == "" {
		return "", fmt.Errorf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID)
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: successfully added user [%s]", userName)
		return fmt.Sprintf("dry-run-%s-%s-guid", userName, origin), nil
	}

	createdUser, err := m.Client.CreateUser(uaaclient.User{
		Username:   userName,
		ExternalID: externalID,
		Origin:     origin,
		Emails: []uaaclient.Email{
			{
				Value: userEmail,
			},
		},
	})
	if err != nil {
		return "", err
	}
	lo.G.Infof("successfully added user [%s]", userName)
	return createdUser.ID, nil
}

// ListUsers - returns uaa.Users
func (m *DefaultUAAManager) ListUsers() (*Users, error) {
	users := &Users{}
	lo.G.Debug("Getting users from Cloud Foundry")
	userList, err := m.Client.ListAllUsers("", "", "userName,id,externalId,emails,origin", "")
	if err != nil {
		return nil, err
	}

	lo.G.Debugf("Found %d users in the CF instance", len(userList))
	for _, user := range userList {
		userName := strings.Trim(user.Username, " ")
		externalID := strings.Trim(user.ExternalID, " ")
		lo.G.Debugf("Adding to users userID [%s], externalID [%s], origin [%s], email [%s], GUID [%s]", userName, externalID, user.Origin, Email(user), user.ID)
		users.Add(User{
			Username:   userName,
			ExternalID: user.ExternalID,
			Email:      Email(user),
			Origin:     user.Origin,
			GUID:       user.ID,
		})
	}
	return users, nil
}

func Email(u uaaclient.User) string {
	for _, email := range u.Emails {
		if email.Primary == nil {
			continue
		}
		if *email.Primary {
			return email.Value
		}
	}
	if len(u.Emails) > 0 {
		return u.Emails[0].Value
	}
	return ""
}
