package uaa

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
	"github.com/xchapter7x/lo"
)

//NewHttpManager -
func NewHttpManager() (mgr HttpManager) {
	return &DefaultHttpManager{}
}

func ShallowDefaultTransport() *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	return &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
	}
}

//Put -
func (m *DefaultHttpManager) Put(url, token, payload string) error {
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	put := request.Put(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+token)
	put.Send(payload)
	res, body, errs := put.End()
	if len(errs) > 0 {
		return errs[0]
	}
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}

	return nil
}

//Post -
func (m *DefaultHttpManager) Post(url, token, payload string) (string, error) {
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	post := request.Post(url)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.Set("Authorization", "BEARER "+token)
	post.Send(payload)
	res, body, errs := post.End()
	if len(errs) > 0 {
		return "", errs[0]
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return body, nil
}

//Get - return struct marshalled into target
func (m *DefaultHttpManager) Get(url, token string, target interface{}) error {
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	get := request.Get(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	res, body, errs := get.End()
	if len(errs) > 0 {
		lo.G.Error(errs)
		return errs[0]
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("get: status %d, body %s", res.StatusCode, body)
	}
	return json.Unmarshal([]byte(body), &target)
}

// Delete deletes a given resource on the server.
func (m *DefaultHttpManager) Delete(url, token string) error {
	request := gorequest.New()
	request.Transport = ShallowDefaultTransport()
	get := request.Delete(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	res, _, errs := get.End()
	if len(errs) > 0 {
		lo.G.Error(errs)
		return errs[0]
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete: call failed with status %d", res.StatusCode)
	}
	return nil
}

//HttpManager -
type HttpManager interface {
	Put(url, token, payload string) (err error)
	Post(url, token, payload string) (body string, err error)
	Get(url, token string, target interface{}) (err error)
	Delete(url, token string) error
}

//DefaultHttpManager -
type DefaultHttpManager struct {
}
