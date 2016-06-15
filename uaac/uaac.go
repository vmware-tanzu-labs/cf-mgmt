package uaac

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

//NewManager -
func NewManager(systemDomain, uuacToken string) (mgr Manager) {
	return &DefaultUAACManager{
		SystemDomain: systemDomain,
		UUACToken:    uuacToken,
	}
}

//CreateUser -
func (m *DefaultUAACManager) CreateUser(userName, userEmail, userDN string) (err error) {
	var res *http.Response
	var body string
	var errs []error
	url := fmt.Sprintf("https://uaa.%s/Users", m.SystemDomain)
	request := gorequest.New()
	post := request.Post(url)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.Set("Authorization", "BEARER "+m.UUACToken)
	post.Set("Content-Type", "application/json")
	sendString := fmt.Sprintf(`{"userName":"%s","emails":[{"value":"%s"}],"origin":"ldap","externalId":"%s"}`, userName, userEmail, strings.Replace(userDN, "\\,", ",", 1))
	post.Send(sendString)
	if res, body, errs = post.End(); len(errs) == 0 && res.StatusCode == http.StatusCreated {
		fmt.Println("successfully added user", userName)
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf(body)
	}
	return
}

//ListUsers -
func (m *DefaultUAACManager) ListUsers() (users map[string]string, err error) {
	var res *http.Response
	var body string
	var errs []error
	users = make(map[string]string)
	url := fmt.Sprintf("https://uaa.%s/Users", m.SystemDomain)
	request := gorequest.New()
	get := request.Get(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+m.UUACToken)

	if res, body, errs = get.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
		userList := new(UserList)
		if err = json.Unmarshal([]byte(body), &userList); err == nil {
			userList := userList.Users
			for _, user := range userList {
				users[user.Name] = user.ID
			}
		}
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf(body)
	}

	return
}

//DefaultUAACManager -
type DefaultUAACManager struct {
	SystemDomain string
	UUACToken    string
}
