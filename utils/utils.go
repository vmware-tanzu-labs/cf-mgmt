package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/parnurzeal/gorequest"
)

//NewDefaultManager -
func NewDefaultManager() (mgr Manager) {
	return &DefaultManager{}
}

//FindFiles -
func (m *DefaultManager) FindFiles(configDir, pattern string) (files []string, err error) {
	m.filePattern = pattern
	filepath.Walk(configDir, m.walkDirectories)
	files = m.filePaths
	return
}

//HTTPPut -
func (m *DefaultManager) HTTPPut(url, token, payload string) error {
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

//HTTPPost -
func (m *DefaultManager) HTTPPost(url, token, payload string) (string, error) {
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

//HTTPGet - return struct marshalled into target
func (m *DefaultManager) HTTPGet(url, token string, target interface{}) error {
	request := gorequest.New()
	get := request.Get(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	res, body, errs := get.End()
	if len(errs) > 0 {
		return errs[0]
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(body)
	}
	return json.Unmarshal([]byte(body), &target)
}

func (m *DefaultManager) walkDirectories(path string, info os.FileInfo, e error) (err error) {
	if strings.Contains(path, m.filePattern) {
		m.filePaths = append(m.filePaths, path)
	}
	return
}

//LoadFile -
func (m *DefaultManager) LoadFile(configFile string, dataType interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err == nil {
		err = yaml.Unmarshal(data, dataType)
	}
	return
}

//WriteFile -
func (m *DefaultManager) WriteFile(configFile string, dataType interface{}) (err error) {
	var data []byte
	if data, err = yaml.Marshal(dataType); err == nil {
		err = ioutil.WriteFile(configFile, data, 0755)
	}
	return
}

//Manager -
type Manager interface {
	FindFiles(directoryName, pattern string) (files []string, err error)
	HTTPPut(url, token, payload string) (err error)
	HTTPPost(url, token, payload string) (body string, err error)
	HTTPGet(url, token string, target interface{}) (err error)
	LoadFile(configFile string, dataType interface{}) (err error)
	WriteFile(configFile string, dataType interface{}) (err error)
}

//DefaultManager -
type DefaultManager struct {
	filePattern string
	filePaths   []string
}
