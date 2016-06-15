package utils

import (
	"os"
	"path/filepath"
	"strings"
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

func (m *DefaultManager) walkDirectories(path string, info os.FileInfo, e error) (err error) {
	if strings.Contains(path, m.filePattern) {
		m.filePaths = append(m.filePaths, path)
	}
	return
}

//Manager -
type Manager interface {
	FindFiles(directoryName, pattern string) (files []string, err error)
}

//DefaultManager -
type DefaultManager struct {
	filePattern string
	filePaths   []string
}
