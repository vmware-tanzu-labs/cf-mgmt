package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

//NewDefaultManager -
func NewDefaultManager() (mgr Manager) {
	return &DefaultManager{}
}

//FindFiles -
func (m *DefaultManager) FindFiles(configDir, pattern string) ([]string, error) {
	m.filePattern = pattern
	err := filepath.Walk(configDir, m.walkDirectories)
	return m.filePaths, err
}

func (m *DefaultManager) walkDirectories(path string, info os.FileInfo, e error) error {
	if strings.Contains(path, m.filePattern) {
		m.filePaths = append(m.filePaths, path)
	}
	return e
}

//FileOrDirectoryExists - checks if file exists
func (m *DefaultManager) FileOrDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

//LoadFile -
func (m *DefaultManager) LoadFile(configFile string, dataType interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, dataType)
}

//LoadJSONFile - this is a hack
func (m *DefaultManager) LoadJSONFile(configFile string, dataType interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dataType)
}

//WriteFileBytes -
func (m *DefaultManager) WriteFileBytes(configFile string, data []byte) error {
	return ioutil.WriteFile(configFile, data, 0755)
}

//WriteFile -
func (m *DefaultManager) WriteFile(configFile string, dataType interface{}) error {
	data, err := yaml.Marshal(dataType)
	if err != nil {
		return err
	}
	return m.WriteFileBytes(configFile, data)
}

//Manager -
type Manager interface {
	FindFiles(directoryName, pattern string) ([]string, error)
	LoadFile(configFile string, dataType interface{}) error
	LoadJSONFile(configFile string, dataType interface{}) error
	WriteFile(configFile string, dataType interface{}) error
	WriteFileBytes(configFile string, data []byte) error
	FileOrDirectoryExists(path string) bool
}

//DefaultManager -
type DefaultManager struct {
	filePattern string
	filePaths   []string
}
