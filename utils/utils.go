package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//NewDefaultManager -
func NewDefaultManager() (mgr Manager) {
	return &DefaultManager{}
}

//FindFiles -
func (m *DefaultManager) FindFiles(configDir, pattern string) (files []string, err error) {
	m.filePattern = pattern
	err = filepath.Walk(configDir, m.walkDirectories)
	files = m.filePaths
	return
}

func (m *DefaultManager) walkDirectories(path string, info os.FileInfo, e error) (err error) {
	if strings.Contains(path, m.filePattern) {
		m.filePaths = append(m.filePaths, path)
	}
	return e
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
	LoadFile(configFile string, dataType interface{}) (err error)
	WriteFile(configFile string, dataType interface{}) (err error)
}

//DefaultManager -
type DefaultManager struct {
	filePattern string
	filePaths   []string
}
