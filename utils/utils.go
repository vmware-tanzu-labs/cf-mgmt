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
func (m *DefaultManager) FindFiles(configDir, pattern string) ([]string, error) {
	var foundFiles = make([]string, 0)
	err := filepath.Walk(configDir,
		func(path string, info os.FileInfo, e error) error {
			if strings.Contains(path, pattern) {
				foundFiles = append(foundFiles, path)
			}
			return e
		})
	return foundFiles, err
}

//DeleteDirectory - deletes a directory
func (m *DefaultManager) DeleteDirectory(path string) error {
	err := os.RemoveAll(path)
	return err
}

//FileOrDirectoryExists - checks if file exists
func (m *DefaultManager) FileOrDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

//LoadFileBytes - Load a file and return the bytes
func (m *DefaultManager) LoadFileBytes(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
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
	WriteFile(configFile string, dataType interface{}) error
	WriteFileBytes(configFile string, data []byte) error
	FileOrDirectoryExists(path string) bool
	LoadFileBytes(path string) ([]byte, error)
	DeleteDirectory(path string) error
}

//DefaultManager -
type DefaultManager struct {
}
