package uaa

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"
	http2 "github.com/pivotalservices/cf-mgmt/http"
	"github.com/xchapter7x/lo"
)

//Manager -
type Manager interface {
	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (map[string]string, error)

	// Returns a map keyed by userid and value as User struct.
	// Return an empty map if an error occurs or if there are no users found
	UsersByID() (map[string]User, error)

	CreateExternalUser(userName, userEmail, externalID, origin string) (err error)
}

//UserList -
type UserList struct {
	Users []User `json:"resources"`
}

//User -
type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Origin   string `json:"origin"`
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	Host  string
	Token string
}

//NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, token string) Manager {
	return &DefaultUAAManager{
		Host:  fmt.Sprintf("https://uaa.%s", sysDomain),
		Token: token,
	}
}

//GetCFToken -
func GetCFToken(host, userID, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/oauth/token", host)
	post := gorequest.New().Post(tokenURL)
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
		msg := fmt.Sprintf("skipping user as missing name[%s], email[%s] or externalID[%s]", userName, userEmail, externalID)
		return errors.New(msg)
	}
	url := fmt.Sprintf("%s/Users", m.Host)
	payload := fmt.Sprintf(`{"userName":"%s","emails":[{"value":"%s"}],"origin":"%s","externalId":"%s"}`, userName, userEmail, origin, strings.Replace(externalID, "\\,", ",", 1))
	if _, err := http2.NewManager().Post(url, m.Token, payload); err != nil {
		return err
	}
	lo.G.Info("successfully added user", userName)
	return nil
}

//ListUsers - Returns a map containing username as key and user guid as value
func (m *DefaultUAAManager) ListUsers() (map[string]string, error) {
	userIDMap := make(map[string]string)
	usersList, err := getUsers(m.Host, m.Token)
	if err != nil {
		return nil, err
	}
	for _, user := range usersList.Users {
		userIDMap[strings.ToLower(user.UserName)] = user.ID
	}
	return userIDMap, nil
}

// UsersByID returns a map of Users keyed by ID.
func (m *DefaultUAAManager) UsersByID() (userIDMap map[string]User, err error) {
	userIDMap = make(map[string]User)
	userList, err := getUsers(m.Host, m.Token)
	if err != nil {
		return nil, err
	}
	for _, user := range userList.Users {
		userIDMap[strings.ToLower(user.UserName)] = user
	}
	return userIDMap, nil
}

//TODO Anwar - Make this API use pagination
func getUsers(host string, uaacToken string) (userList *UserList, err error) {
	lo.G.Debug("Getting users from Cloud Foundry")
	url := fmt.Sprintf("%s/Users?count=5000", host)
	userList = new(UserList)
	if err := http2.NewManager().Get(url, uaacToken, userList); err != nil {
		return nil, fmt.Errorf("couldn't retrieve users: %v", err)
	}
	lo.G.Debugf("Found %d users in the CF instance", len(userList.Users))
	return userList, nil
}
