package uaa

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"

	"github.com/xchapter7x/lo"
)

//Manager -
type Manager interface {
	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (map[string]*User, error)
	CreateExternalUser(userName, userEmail, externalID, origin string) (err error)
}

//UserList -
type UserList struct {
	Users        []User `json:"resources"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
	TotalResults int    `json:"totalResults"`
}

//User -
type User struct {
	ID         string      `json:"id"`
	UserName   string      `json:"userName"`
	ExternalID string      `json:"externalId"`
	Origin     string      `json:"origin"`
	Emails     []UserEmail `json:"emails"`
}

func (u *User) Email() string {
	for _, email := range u.Emails {
		if email.Primary {
			return email.Value
		}
	}
	if len(u.Emails) > 0 {
		return u.Emails[0].Value
	}
	return ""
}

type UserEmail struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	Host  string
	Token string
	Http  HttpManager
	Peek  bool
}

type Pagination interface {
	GetNextURL(url string) string
	AddInstances(Pagination)
}

//NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, token string, peek bool) Manager {
	return &DefaultUAAManager{
		Host:  fmt.Sprintf("https://uaa.%s", sysDomain),
		Token: token,
		Http:  NewHttpManager(),
		Peek:  peek,
	}
}

//GetCFToken -
func GetCFToken(host, userID, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/oauth/token", host)
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	post := request.Post(tokenURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = "cf"
	post.BasicAuth.Password = ""

	params := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {userID},
		"password":      {password},
	}
	post.Send(params.Encode())
	res, body, errs := post.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot get CF token, error %v: %s", res.StatusCode, body)
	}

	t := Token{}
	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return "", err
	}

	return t.AccessToken, nil
}

//GetUAACToken -
func GetUAACToken(host, userID, secret string) (string, error) {
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	request.TargetType = "form"
	post := request.Post(fmt.Sprintf("%s/oauth/token", host))
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = userID
	post.BasicAuth.Password = secret

	params := url.Values{
		"grant_type":    {"client_credentials"},
		"response_type": {"token"},
	}
	post.Send(params.Encode())

	res, body, errs := post.End()
	if len(errs) > 0 {
		return "", errs[0]
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot get UAAC token, error %v: %s", res.StatusCode, body)
	}
	t := Token{}
	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return "", fmt.Errorf("cannot read token: %v", err)
	}
	return t.AccessToken, nil
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
	url := fmt.Sprintf("%s/Users", m.Host)
	payload := fmt.Sprintf(`{"userName":"%s","emails":[{"value":"%s"}],"origin":"%s","externalId":"%s"}`, userName, userEmail, origin, strings.Replace(externalID, "\\,", ",", -1))
	if _, err := m.Http.Post(url, m.Token, payload); err != nil {
		return err
	}
	lo.G.Infof("successfully added user [%s]", userName)
	return nil
}

//ListUsers - Returns a map containing username as key and user guid as value
func (m *DefaultUAAManager) ListUsers() (map[string]*User, error) {
	userMap := make(map[string]*User)
	usersList, err := m.getUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range usersList.Users {
		userMap[strings.ToLower(user.UserName)] = &user
		if user.ExternalID != "" {
			userMap[strings.ToLower(user.ExternalID)] = &user
		}
	}
	return userMap, nil
}

//TODO Anwar - Make this API use pagination
func (m *DefaultUAAManager) getUsers() (*UserList, error) {
	lo.G.Debug("Getting users from Cloud Foundry")
	url := fmt.Sprintf("%s/Users?count=5000", m.Host)
	userList := &UserList{}
	err := m.listResources(url, userList, NewUserListResources)
	if err != nil {
		return nil, err
	}
	lo.G.Debugf("Found %d users in the CF instance", len(userList.Users))
	return userList, nil
}

func (m *DefaultUAAManager) listResources(url string, target Pagination, createInstance func() Pagination) error {
	var err = m.Http.Get(url, m.Token, target)
	if err != nil {
		return err
	}
	if target.GetNextURL(url) == "" {
		return nil
	}

	nextURL := target.GetNextURL(url)
	for nextURL != "" {
		lo.G.Debugf("NextURL: %s", nextURL)
		tempTarget := createInstance()
		err = m.Http.Get(nextURL, m.Token, tempTarget)
		if err != nil {
			return err
		}
		target.AddInstances(tempTarget)
		nextURL = tempTarget.GetNextURL(url)
	}
	return nil
}
