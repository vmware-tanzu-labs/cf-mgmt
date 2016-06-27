package utils

import (
	"crypto/tls"
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

//HTTPDelete -
func (m *DefaultManager) HTTPDelete(url, token, payload string) (err error) {
	var res *http.Response
	var body string
	var errs []error
	request := gorequest.New()
	put := request.Delete(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+token)
	put.Send(payload)
	if res, _, errs = put.End(); len(errs) == 0 && res.StatusCode == http.StatusCreated {
		return
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return
}

//HTTPPut -
func (m *DefaultManager) HTTPPut(url, token, payload string) (err error) {
	var res *http.Response
	var body string
	var errs []error
	request := gorequest.New()
	put := request.Put(url)
	put.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	put.Set("Authorization", "BEARER "+token)
	put.Send(payload)
	if res, _, errs = put.End(); len(errs) == 0 && res.StatusCode == http.StatusCreated {
		return
	} else if len(errs) > 0 {
		err = errs[0]
	} else {
		err = fmt.Errorf("Status %d, body %s", res.StatusCode, body)
	}
	return
}

//HTTPPost -
func (m *DefaultManager) HTTPPost(url, token, payload string) (body string, err error) {
	var res *http.Response
	var errs []error
	request := gorequest.New()
	post := request.Post(url)
	post.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	post.Set("Authorization", "BEARER "+token)
	post.Send(payload)
	if res, body, errs = post.End(); len(errs) != 0 || res.StatusCode != http.StatusOK {
		if len(errs) > 0 {
			err = errs[0]
		} else {
			err = fmt.Errorf(body)
		}
	}
	return
}

//HTTPGet -
func (m *DefaultManager) HTTPGet(url, token string) (body string, err error) {
	var res *http.Response
	var errs []error
	request := gorequest.New()
	get := request.Get(url)
	get.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	get.Set("Authorization", "BEARER "+token)

	if res, body, errs = get.End(); len(errs) != 0 || res.StatusCode != http.StatusOK {
		if len(errs) > 0 {
			err = errs[0]
		} else {
			err = fmt.Errorf(body)
		}
	}
	return
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
	HTTPDelete(url, token, payload string) (err error)
	HTTPPost(url, token, payload string) (body string, err error)
	HTTPGet(url, token string) (body string, err error)
	LoadFile(configFile string, dataType interface{}) (err error)
	WriteFile(configFile string, dataType interface{}) (err error)
}

//DefaultManager -
type DefaultManager struct {
	filePattern string
	filePaths   []string
}
