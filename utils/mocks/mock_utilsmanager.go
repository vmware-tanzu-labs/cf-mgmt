package mock_utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	reflect "reflect"
	"strings"
)

// MockUtilsManager is a mock of Manager interface
type MockUtilsManager struct {
	MockFileData         map[string]interface{}
	MockFileDataHasError bool
}

// NewMockUtilsManagercreates a new mock utils instance
func NewMockUtilsManager() *MockUtilsManager {
	mapinst := make(map[string]interface{}, 0)
	mock := &MockUtilsManager{
		MockFileData:         mapinst,
		MockFileDataHasError: false,
	}
	return mock
}

//LoadFileBytes - Load a file and return the bytes
func (m *MockUtilsManager) LoadFileBytes(arg0 string) ([]byte, error) {
	data, exists := m.MockFileData[arg0]
	if exists {
		// Convert interface {} to bytes
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("%s not created in test data", arg0)
}

// DeleteDirectory mocks base method
func (m *MockUtilsManager) DeleteDirectory(arg0 string) error {
	var deleteKeys []string
	for k := range m.MockFileData {
		if strings.HasPrefix(k, arg0) {
			deleteKeys = append(deleteKeys, k)
		}
	}

	for _, keyToDelete := range deleteKeys {
		delete(m.MockFileData, keyToDelete)
	}
	return nil
}

// FileOrDirectoryExists mocks base method
func (m *MockUtilsManager) FileOrDirectoryExists(arg0 string) bool {
	_, exists := m.MockFileData[arg0]
	return exists
}

// FindFiles mocks base method
func (m *MockUtilsManager) FindFiles(arg0, arg1 string) ([]string, error) {
	var err error = nil
	foundFiles := make([]string, 0)
	for k := range m.MockFileData {
		if strings.HasPrefix(k, arg0+"/") && strings.Contains(k, arg1) {
			foundFiles = append(foundFiles, k)
		}
	}

	if len(foundFiles) == 0 {
		err = fmt.Errorf("Nothing found in test data")
	}

	return foundFiles, err
}

// LoadFile mocks base method
func (m *MockUtilsManager) LoadFile(arg0 string, arg1 interface{}) error {
	var err error = nil

	// Make dir/file into ./dir/file
	if !strings.HasPrefix(arg0, "./") && !strings.HasPrefix(arg0, "/") {
		arg0 = "./" + arg0
	}

	if !m.FileOrDirectoryExists(arg0) {
		return nil
	}

	// Get the value type of arg1 and then perform a set of valuetype of the MockFileData[arg0]
	v := reflect.ValueOf(arg1)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	v.Set(reflect.ValueOf(m.MockFileData[arg0]))
	if m.MockFileDataHasError {
		err = fmt.Errorf("Error Injected via MockFileDataHasError flag")
	}
	return err
}

// WriteFile mocks base method
func (m *MockUtilsManager) WriteFile(arg0 string, arg1 interface{}) error {
	m.MockFileData[arg0] = arg1
	return nil
}

// WriteFileBytes mocks base method
func (m *MockUtilsManager) WriteFileBytes(arg0 string, arg1 []byte) error {
	m.MockFileData[arg0] = arg1
	return nil
}
