package uaa

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/xchapter7x/lo"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
)

type uaa interface {
	CreateUser(user uaaclient.User) (*uaaclient.User, error)
	ListUsers(filter string, sortBy string, attributes string, sortOrder uaaclient.SortOrder, startIndex int, itemsPerPage int) ([]uaaclient.User, uaaclient.Page, error)
}

// Manager -
type Manager interface {
	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (*Users, error)
	CreateExternalUser(userName, userEmail, externalID, origin string) (err error)
}

// DefaultUAAManager -
type DefaultUAAManager struct {
	Peek   bool
	Client uaa
	Users  *Users
}

type User struct {
	Username   string
	ExternalID string
	Email      string
	Origin     string
	GUID       string
}

// NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, clientID, clientSecret, userAgent string, httpClient *http.Client, peek bool) (Manager, error) {
	target := fmt.Sprintf("https://uaa.%s", sysDomain)

	client, err := uaaclient.New(
		target,
		uaaclient.WithClientCredentials(clientID, clientSecret, uaaclient.OpaqueToken),
		uaaclient.WithUserAgent(userAgent),
		uaaclient.WithClient(httpClient),
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

func (m *DefaultUAAManager) addUser(user User) {
	if m.Users == nil {
		m.Users = &Users{}
	}
	m.Users.Add(user)
}

// CreateExternalUser -
func (m *DefaultUAAManager) CreateExternalUser(userName, userEmail, externalID, origin string) error {
	if userName == "" || userEmail == "" || externalID == "" {
		return fmt.Errorf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID)
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: successfully added user [%s]", userName)
		m.addUser(User{
			Username:   userName,
			Email:      userEmail,
			ExternalID: externalID,
			Origin:     origin,
			GUID:       fmt.Sprintf("dry-run-%s-%s-guid", userName, origin),
		})
		return nil
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
		var requestError uaaclient.RequestError
		if errors.As(err, &requestError) {
			return fmt.Errorf("got an error calling %s with response %s", requestError.Url, requestError.ErrorResponse)
		}
		return err
	}
	m.addUser(User{
		Username:   userName,
		Email:      userEmail,
		ExternalID: externalID,
		Origin:     origin,
		GUID:       createdUser.ID,
	})
	lo.G.Infof("successfully added user [%s]", userName)
	return nil
}

// ListUsers - returns uaa.Users
func (m *DefaultUAAManager) ListUsers() (*Users, error) {
	if m.Users != nil {
		return m.Users, nil
	}

	users := &Users{}
	lo.G.Debug("Getting users from UAA")
	userList, err := m.ListAllUsers()
	if err != nil {
		var requestError uaaclient.RequestError
		if errors.As(err, &requestError) {
			return nil, fmt.Errorf("got an error calling %s with response %s", requestError.Url, requestError.ErrorResponse)
		}
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

	m.Users = users

	return users, nil
}

func (m *DefaultUAAManager) ListAllUsers() ([]uaaclient.User, error) {
	page := uaaclient.Page{
		StartIndex:   1,
		ItemsPerPage: 500,
	}
	var (
		results      []uaaclient.User
		currentPage  []uaaclient.User
		err          error
		totalResults int
	)
	currentPage, page, err = m.Client.ListUsers("", "id", "userName,id,externalId,emails,origin", "", page.StartIndex, page.ItemsPerPage)
	totalResults = page.TotalResults
	if err != nil {
		return nil, err
	}
	results = append(results, currentPage...)
	if (page.StartIndex + page.ItemsPerPage) <= page.TotalResults {
		page.StartIndex = page.StartIndex + page.ItemsPerPage
		for {
			currentPage, page, err = m.Client.ListUsers("", "id", "userName,id,externalId,emails,origin", "", page.StartIndex, page.ItemsPerPage)
			if err != nil {
				return nil, err
			}
			if totalResults != page.TotalResults {
				lo.G.Infof("Result size changed during pagination from %d to %d", totalResults, page.TotalResults)
				totalResults = page.TotalResults
			}
			results = append(results, currentPage...)

			if (page.StartIndex + page.ItemsPerPage) > page.TotalResults {
				break
			}
			page.StartIndex = page.StartIndex + page.ItemsPerPage
		}
	}
	if len(results) != totalResults {
		return nil, fmt.Errorf("results %d is not equal to expected results %d", len(results), totalResults)
	}
	uniqueMap := make(map[string]string)
	for _, user := range results {
		userKey := fmt.Sprintf("%s-%s", user.Username, user.Origin)
		if userId, ok := uniqueMap[userKey]; ok {
			return nil, fmt.Errorf("user with userName [%s], origin [%s], id [%s] already returned with id [%s]", user.Username, user.Origin, userId, user.ID)
		} else {
			uniqueMap[userKey] = user.ID
		}
	}
	return results, nil
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
