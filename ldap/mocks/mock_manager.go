// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pivotalservices/cf-mgmt/ldap (interfaces: Manager)

// Package mock_ldap is a generated GoMock package.
package mock_ldap

import (
	gomock "github.com/golang/mock/gomock"
	ldap "github.com/pivotalservices/cf-mgmt/ldap"
	ldap0 "github.com/go-ldap/ldap"
	reflect "reflect"
)

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GetConfig mocks base method
func (m *MockManager) GetConfig(arg0, arg1 string) (*ldap.Config, error) {
	ret := m.ctrl.Call(m, "GetConfig", arg0, arg1)
	ret0, _ := ret[0].(*ldap.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig
func (mr *MockManagerMockRecorder) GetConfig(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockManager)(nil).GetConfig), arg0, arg1)
}

// GetLdapUser mocks base method
func (m *MockManager) GetLdapUser(arg0 *ldap.Config, arg1, arg2 string) (*ldap.User, error) {
	ret := m.ctrl.Call(m, "GetLdapUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ldap.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLdapUser indicates an expected call of GetLdapUser
func (mr *MockManagerMockRecorder) GetLdapUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLdapUser", reflect.TypeOf((*MockManager)(nil).GetLdapUser), arg0, arg1, arg2)
}

// GetUser mocks base method
func (m *MockManager) GetUser(arg0 *ldap.Config, arg1 string) (*ldap.User, error) {
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*ldap.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockManagerMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockManager)(nil).GetUser), arg0, arg1)
}

// GetUserIDs mocks base method
func (m *MockManager) GetUserIDs(arg0 *ldap.Config, arg1 string) ([]ldap.User, error) {
	ret := m.ctrl.Call(m, "GetUserIDs", arg0, arg1)
	ret0, _ := ret[0].([]ldap.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDs indicates an expected call of GetUserIDs
func (mr *MockManagerMockRecorder) GetUserIDs(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDs", reflect.TypeOf((*MockManager)(nil).GetUserIDs), arg0, arg1)
}

// LdapConnection mocks base method
func (m *MockManager) LdapConnection(arg0 *ldap.Config) (*ldap0.Conn, error) {
	ret := m.ctrl.Call(m, "LdapConnection", arg0)
	ret0, _ := ret[0].(*ldap0.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LdapConnection indicates an expected call of LdapConnection
func (mr *MockManagerMockRecorder) LdapConnection(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LdapConnection", reflect.TypeOf((*MockManager)(nil).LdapConnection), arg0)
}
