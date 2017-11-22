package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

//FindFiles -
func FindFiles(configDir, pattern string) ([]string, error) {
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
func DeleteDirectory(path string) error {
	err := os.RemoveAll(path)
	return err
}

//FileOrDirectoryExists - checks if file exists
func FileOrDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

//LoadFileBytes - Load a file and return the bytes
func LoadFileBytes(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

//LoadFile -
func LoadFile(configFile string, dataType interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, dataType)
}

//WriteFileBytes -
func WriteFileBytes(configFile string, data []byte) error {
	return ioutil.WriteFile(configFile, data, 0755)
}

//WriteFile -
func WriteFile(configFile string, dataType interface{}) error {
	data, err := yaml.Marshal(dataType)
	if err != nil {
		return err
	}
	return WriteFileBytes(configFile, data)
}
