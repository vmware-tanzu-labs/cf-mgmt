package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager() (mgr Manager) {
	return &DefaultManager{}
}

//Put -
func (m *DefaultManager) Put(url, token, payload string) error {
	request := gorequest.New()
	put := request.Put(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+token)
	put.Send(payload)
	res, body, errs := put.End()
	if len(errs) > 0 {
		return errs[0]
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}

	return nil
}

//Post -
func (m *DefaultManager) Post(url, token, payload string) (string, error) {
	request := gorequest.New()
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
func (m *DefaultManager) Get(url, token string, target interface{}) error {
	request := gorequest.New()
	get := request.Get(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	res, body, errs := get.End()
	if len(errs) > 0 {
		lo.G.Error(errs)
		return errs[0]
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return json.Unmarshal([]byte(body), &target)
}

//Delete- Deletes a given resource on the server
func (m *DefaultManager) Delete(url, token string) error {
	request := gorequest.New()
	get := request.Delete(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	res, _, errs := get.End()
	if len(errs) > 0 {
		lo.G.Error(errs)
		return errs[0]
	}
	if res.StatusCode != http.StatusOK || res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Delete call failed with Status : %d", res.StatusCode)
	}
	return nil
}

//Manager -
type Manager interface {
	Put(url, token, payload string) (err error)
	Post(url, token, payload string) (body string, err error)
	Get(url, token string, target interface{}) (err error)
	Delete(url, token string) error
}

//DefaultManager -
type DefaultManager struct {
}
