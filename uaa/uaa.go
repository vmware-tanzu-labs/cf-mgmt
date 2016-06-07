package uaa

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"
)

//NewDefaultUAAManager -
func NewDefaultUAAManager(sysDomain, userID, pwd string) (mgr Manager) {
	return &DefaultUAAManager{
		SysDomain: sysDomain,
		UserID:    userID,
		Password:  pwd,
	}
}

//GetToken -
func (m *DefaultUAAManager) GetToken() (token string, err error) {
	var res *http.Response
	var body string
	var errs []error
	request := gorequest.New()
	tokenURL := fmt.Sprintf("https://login.%s/oauth/token", m.SysDomain)
	params := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {m.UserID},
		"password":      {m.Password},
	}
	header, _ := "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")), strings.NewReader(params.Encode())
	sendString := params.Encode()

	request.TargetType = "form"
	post := request.Post(tokenURL)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.Set("Authorization", header)
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
