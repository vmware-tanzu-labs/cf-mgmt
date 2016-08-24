package uaa

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/parnurzeal/gorequest"
)

//Manager -
type Manager interface {
	GetCFToken(password string) (token string, err error)
	GetUAACToken(secret string) (token string, err error)
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	Host   string
	UserID string
}

//NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, userID string) (mgr Manager) {
	return &DefaultUAAManager{
		Host:   fmt.Sprintf("https://uaa.%s", sysDomain),
		UserID: userID,
	}
}

//GetCFToken -
func (m *DefaultUAAManager) GetCFToken(password string) (string, error) {
	request := gorequest.New()
	tokenURL := fmt.Sprintf("%s/oauth/token", m.Host)
	params := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {m.UserID},
		"password":      {password},
	}
	sendString := params.Encode()

	post := request.Post(tokenURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = "cf"
	post.BasicAuth.Password = ""
	post.Send(sendString)
	res, body, errs := post.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(body)
	}

	t := new(Token)
	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return "", err
	}

	return t.AccessToken, nil
}

//GetUAACToken -
func (m *DefaultUAAManager) GetUAACToken(secret string) (token string, err error) {
	var res *http.Response
	var body string
	var errs []error
	request := gorequest.New()
	tokenURL := fmt.Sprintf("%s/oauth/token", m.Host)
	params := url.Values{
		"grant_type":    {"client_credentials"},
		"response_type": {"token"},
	}
	sendString := params.Encode()

	request.TargetType = "form"
	post := request.Post(tokenURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = m.UserID
	post.BasicAuth.Password = secret
	post.Send(sendString)
	if res, body, errs = post.End(); len(errs) == 0 && res.StatusCode == http.StatusOK {
		t := new(Token)
		if err = json.Unmarshal([]byte(body), &t); err == nil {
			token = t.AccessToken
		}
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf(body)
	}
	return
}
