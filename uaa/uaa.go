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
	GetCFToken(password string) (string, error)
	GetUAACToken(secret string) (string, error)
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
func NewDefaultUAAManager(sysDomain, userID string) Manager {
	return &DefaultUAAManager{
		Host:   fmt.Sprintf("https://uaa.%s", sysDomain),
		UserID: userID,
	}
}

//GetCFToken -
func (m *DefaultUAAManager) GetCFToken(password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/oauth/token", m.Host)
	post := gorequest.New().Post(tokenURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = "cf"
	post.BasicAuth.Password = ""

	params := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {m.UserID},
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
func (m *DefaultUAAManager) GetUAACToken(secret string) (string, error) {
	request := gorequest.New()
	request.TargetType = "form"
	post := request.Post(fmt.Sprintf("%s/oauth/token", m.Host))
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.BasicAuth.Username = m.UserID
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
