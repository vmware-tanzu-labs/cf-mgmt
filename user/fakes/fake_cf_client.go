// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"net/url"
	"sync"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/user"
)

type FakeCFClient struct {
	CreateV3OrganizationRoleStub        func(string, string, string) (*cfclient.V3Role, error)
	createV3OrganizationRoleMutex       sync.RWMutex
	createV3OrganizationRoleArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	createV3OrganizationRoleReturns struct {
		result1 *cfclient.V3Role
		result2 error
	}
	createV3OrganizationRoleReturnsOnCall map[int]struct {
		result1 *cfclient.V3Role
		result2 error
	}
	CreateV3SpaceRoleStub        func(string, string, string) (*cfclient.V3Role, error)
	createV3SpaceRoleMutex       sync.RWMutex
	createV3SpaceRoleArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	createV3SpaceRoleReturns struct {
		result1 *cfclient.V3Role
		result2 error
	}
	createV3SpaceRoleReturnsOnCall map[int]struct {
		result1 *cfclient.V3Role
		result2 error
	}
	DeleteUserStub        func(string) error
	deleteUserMutex       sync.RWMutex
	deleteUserArgsForCall []struct {
		arg1 string
	}
	deleteUserReturns struct {
		result1 error
	}
	deleteUserReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteV3RoleStub        func(string) error
	deleteV3RoleMutex       sync.RWMutex
	deleteV3RoleArgsForCall []struct {
		arg1 string
	}
	deleteV3RoleReturns struct {
		result1 error
	}
	deleteV3RoleReturnsOnCall map[int]struct {
		result1 error
	}
	ListSpacesByQueryStub        func(url.Values) ([]cfclient.Space, error)
	listSpacesByQueryMutex       sync.RWMutex
	listSpacesByQueryArgsForCall []struct {
		arg1 url.Values
	}
	listSpacesByQueryReturns struct {
		result1 []cfclient.Space
		result2 error
	}
	listSpacesByQueryReturnsOnCall map[int]struct {
		result1 []cfclient.Space
		result2 error
	}
	ListV3OrganizationRolesByGUIDAndTypeStub        func(string, string) ([]cfclient.V3User, error)
	listV3OrganizationRolesByGUIDAndTypeMutex       sync.RWMutex
	listV3OrganizationRolesByGUIDAndTypeArgsForCall []struct {
		arg1 string
		arg2 string
	}
	listV3OrganizationRolesByGUIDAndTypeReturns struct {
		result1 []cfclient.V3User
		result2 error
	}
	listV3OrganizationRolesByGUIDAndTypeReturnsOnCall map[int]struct {
		result1 []cfclient.V3User
		result2 error
	}
	ListV3RolesByQueryStub        func(url.Values) ([]cfclient.V3Role, error)
	listV3RolesByQueryMutex       sync.RWMutex
	listV3RolesByQueryArgsForCall []struct {
		arg1 url.Values
	}
	listV3RolesByQueryReturns struct {
		result1 []cfclient.V3Role
		result2 error
	}
	listV3RolesByQueryReturnsOnCall map[int]struct {
		result1 []cfclient.V3Role
		result2 error
	}
	ListV3SpaceRolesByGUIDAndTypeStub        func(string, string) ([]cfclient.V3User, error)
	listV3SpaceRolesByGUIDAndTypeMutex       sync.RWMutex
	listV3SpaceRolesByGUIDAndTypeArgsForCall []struct {
		arg1 string
		arg2 string
	}
	listV3SpaceRolesByGUIDAndTypeReturns struct {
		result1 []cfclient.V3User
		result2 error
	}
	listV3SpaceRolesByGUIDAndTypeReturnsOnCall map[int]struct {
		result1 []cfclient.V3User
		result2 error
	}
	SupportsSpaceSupporterRoleStub        func() (bool, error)
	supportsSpaceSupporterRoleMutex       sync.RWMutex
	supportsSpaceSupporterRoleArgsForCall []struct {
	}
	supportsSpaceSupporterRoleReturns struct {
		result1 bool
		result2 error
	}
	supportsSpaceSupporterRoleReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFClient) CreateV3OrganizationRole(arg1 string, arg2 string, arg3 string) (*cfclient.V3Role, error) {
	fake.createV3OrganizationRoleMutex.Lock()
	ret, specificReturn := fake.createV3OrganizationRoleReturnsOnCall[len(fake.createV3OrganizationRoleArgsForCall)]
	fake.createV3OrganizationRoleArgsForCall = append(fake.createV3OrganizationRoleArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("CreateV3OrganizationRole", []interface{}{arg1, arg2, arg3})
	fake.createV3OrganizationRoleMutex.Unlock()
	if fake.CreateV3OrganizationRoleStub != nil {
		return fake.CreateV3OrganizationRoleStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createV3OrganizationRoleReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) CreateV3OrganizationRoleCallCount() int {
	fake.createV3OrganizationRoleMutex.RLock()
	defer fake.createV3OrganizationRoleMutex.RUnlock()
	return len(fake.createV3OrganizationRoleArgsForCall)
}

func (fake *FakeCFClient) CreateV3OrganizationRoleCalls(stub func(string, string, string) (*cfclient.V3Role, error)) {
	fake.createV3OrganizationRoleMutex.Lock()
	defer fake.createV3OrganizationRoleMutex.Unlock()
	fake.CreateV3OrganizationRoleStub = stub
}

func (fake *FakeCFClient) CreateV3OrganizationRoleArgsForCall(i int) (string, string, string) {
	fake.createV3OrganizationRoleMutex.RLock()
	defer fake.createV3OrganizationRoleMutex.RUnlock()
	argsForCall := fake.createV3OrganizationRoleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCFClient) CreateV3OrganizationRoleReturns(result1 *cfclient.V3Role, result2 error) {
	fake.createV3OrganizationRoleMutex.Lock()
	defer fake.createV3OrganizationRoleMutex.Unlock()
	fake.CreateV3OrganizationRoleStub = nil
	fake.createV3OrganizationRoleReturns = struct {
		result1 *cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) CreateV3OrganizationRoleReturnsOnCall(i int, result1 *cfclient.V3Role, result2 error) {
	fake.createV3OrganizationRoleMutex.Lock()
	defer fake.createV3OrganizationRoleMutex.Unlock()
	fake.CreateV3OrganizationRoleStub = nil
	if fake.createV3OrganizationRoleReturnsOnCall == nil {
		fake.createV3OrganizationRoleReturnsOnCall = make(map[int]struct {
			result1 *cfclient.V3Role
			result2 error
		})
	}
	fake.createV3OrganizationRoleReturnsOnCall[i] = struct {
		result1 *cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) CreateV3SpaceRole(arg1 string, arg2 string, arg3 string) (*cfclient.V3Role, error) {
	fake.createV3SpaceRoleMutex.Lock()
	ret, specificReturn := fake.createV3SpaceRoleReturnsOnCall[len(fake.createV3SpaceRoleArgsForCall)]
	fake.createV3SpaceRoleArgsForCall = append(fake.createV3SpaceRoleArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("CreateV3SpaceRole", []interface{}{arg1, arg2, arg3})
	fake.createV3SpaceRoleMutex.Unlock()
	if fake.CreateV3SpaceRoleStub != nil {
		return fake.CreateV3SpaceRoleStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createV3SpaceRoleReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) CreateV3SpaceRoleCallCount() int {
	fake.createV3SpaceRoleMutex.RLock()
	defer fake.createV3SpaceRoleMutex.RUnlock()
	return len(fake.createV3SpaceRoleArgsForCall)
}

func (fake *FakeCFClient) CreateV3SpaceRoleCalls(stub func(string, string, string) (*cfclient.V3Role, error)) {
	fake.createV3SpaceRoleMutex.Lock()
	defer fake.createV3SpaceRoleMutex.Unlock()
	fake.CreateV3SpaceRoleStub = stub
}

func (fake *FakeCFClient) CreateV3SpaceRoleArgsForCall(i int) (string, string, string) {
	fake.createV3SpaceRoleMutex.RLock()
	defer fake.createV3SpaceRoleMutex.RUnlock()
	argsForCall := fake.createV3SpaceRoleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCFClient) CreateV3SpaceRoleReturns(result1 *cfclient.V3Role, result2 error) {
	fake.createV3SpaceRoleMutex.Lock()
	defer fake.createV3SpaceRoleMutex.Unlock()
	fake.CreateV3SpaceRoleStub = nil
	fake.createV3SpaceRoleReturns = struct {
		result1 *cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) CreateV3SpaceRoleReturnsOnCall(i int, result1 *cfclient.V3Role, result2 error) {
	fake.createV3SpaceRoleMutex.Lock()
	defer fake.createV3SpaceRoleMutex.Unlock()
	fake.CreateV3SpaceRoleStub = nil
	if fake.createV3SpaceRoleReturnsOnCall == nil {
		fake.createV3SpaceRoleReturnsOnCall = make(map[int]struct {
			result1 *cfclient.V3Role
			result2 error
		})
	}
	fake.createV3SpaceRoleReturnsOnCall[i] = struct {
		result1 *cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) DeleteUser(arg1 string) error {
	fake.deleteUserMutex.Lock()
	ret, specificReturn := fake.deleteUserReturnsOnCall[len(fake.deleteUserArgsForCall)]
	fake.deleteUserArgsForCall = append(fake.deleteUserArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteUser", []interface{}{arg1})
	fake.deleteUserMutex.Unlock()
	if fake.DeleteUserStub != nil {
		return fake.DeleteUserStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteUserReturns
	return fakeReturns.result1
}

func (fake *FakeCFClient) DeleteUserCallCount() int {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	return len(fake.deleteUserArgsForCall)
}

func (fake *FakeCFClient) DeleteUserCalls(stub func(string) error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = stub
}

func (fake *FakeCFClient) DeleteUserArgsForCall(i int) string {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	argsForCall := fake.deleteUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) DeleteUserReturns(result1 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	fake.deleteUserReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) DeleteUserReturnsOnCall(i int, result1 error) {
	fake.deleteUserMutex.Lock()
	defer fake.deleteUserMutex.Unlock()
	fake.DeleteUserStub = nil
	if fake.deleteUserReturnsOnCall == nil {
		fake.deleteUserReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteUserReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) DeleteV3Role(arg1 string) error {
	fake.deleteV3RoleMutex.Lock()
	ret, specificReturn := fake.deleteV3RoleReturnsOnCall[len(fake.deleteV3RoleArgsForCall)]
	fake.deleteV3RoleArgsForCall = append(fake.deleteV3RoleArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteV3Role", []interface{}{arg1})
	fake.deleteV3RoleMutex.Unlock()
	if fake.DeleteV3RoleStub != nil {
		return fake.DeleteV3RoleStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteV3RoleReturns
	return fakeReturns.result1
}

func (fake *FakeCFClient) DeleteV3RoleCallCount() int {
	fake.deleteV3RoleMutex.RLock()
	defer fake.deleteV3RoleMutex.RUnlock()
	return len(fake.deleteV3RoleArgsForCall)
}

func (fake *FakeCFClient) DeleteV3RoleCalls(stub func(string) error) {
	fake.deleteV3RoleMutex.Lock()
	defer fake.deleteV3RoleMutex.Unlock()
	fake.DeleteV3RoleStub = stub
}

func (fake *FakeCFClient) DeleteV3RoleArgsForCall(i int) string {
	fake.deleteV3RoleMutex.RLock()
	defer fake.deleteV3RoleMutex.RUnlock()
	argsForCall := fake.deleteV3RoleArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) DeleteV3RoleReturns(result1 error) {
	fake.deleteV3RoleMutex.Lock()
	defer fake.deleteV3RoleMutex.Unlock()
	fake.DeleteV3RoleStub = nil
	fake.deleteV3RoleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) DeleteV3RoleReturnsOnCall(i int, result1 error) {
	fake.deleteV3RoleMutex.Lock()
	defer fake.deleteV3RoleMutex.Unlock()
	fake.DeleteV3RoleStub = nil
	if fake.deleteV3RoleReturnsOnCall == nil {
		fake.deleteV3RoleReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteV3RoleReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) ListSpacesByQuery(arg1 url.Values) ([]cfclient.Space, error) {
	fake.listSpacesByQueryMutex.Lock()
	ret, specificReturn := fake.listSpacesByQueryReturnsOnCall[len(fake.listSpacesByQueryArgsForCall)]
	fake.listSpacesByQueryArgsForCall = append(fake.listSpacesByQueryArgsForCall, struct {
		arg1 url.Values
	}{arg1})
	fake.recordInvocation("ListSpacesByQuery", []interface{}{arg1})
	fake.listSpacesByQueryMutex.Unlock()
	if fake.ListSpacesByQueryStub != nil {
		return fake.ListSpacesByQueryStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listSpacesByQueryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) ListSpacesByQueryCallCount() int {
	fake.listSpacesByQueryMutex.RLock()
	defer fake.listSpacesByQueryMutex.RUnlock()
	return len(fake.listSpacesByQueryArgsForCall)
}

func (fake *FakeCFClient) ListSpacesByQueryCalls(stub func(url.Values) ([]cfclient.Space, error)) {
	fake.listSpacesByQueryMutex.Lock()
	defer fake.listSpacesByQueryMutex.Unlock()
	fake.ListSpacesByQueryStub = stub
}

func (fake *FakeCFClient) ListSpacesByQueryArgsForCall(i int) url.Values {
	fake.listSpacesByQueryMutex.RLock()
	defer fake.listSpacesByQueryMutex.RUnlock()
	argsForCall := fake.listSpacesByQueryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) ListSpacesByQueryReturns(result1 []cfclient.Space, result2 error) {
	fake.listSpacesByQueryMutex.Lock()
	defer fake.listSpacesByQueryMutex.Unlock()
	fake.ListSpacesByQueryStub = nil
	fake.listSpacesByQueryReturns = struct {
		result1 []cfclient.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListSpacesByQueryReturnsOnCall(i int, result1 []cfclient.Space, result2 error) {
	fake.listSpacesByQueryMutex.Lock()
	defer fake.listSpacesByQueryMutex.Unlock()
	fake.ListSpacesByQueryStub = nil
	if fake.listSpacesByQueryReturnsOnCall == nil {
		fake.listSpacesByQueryReturnsOnCall = make(map[int]struct {
			result1 []cfclient.Space
			result2 error
		})
	}
	fake.listSpacesByQueryReturnsOnCall[i] = struct {
		result1 []cfclient.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndType(arg1 string, arg2 string) ([]cfclient.V3User, error) {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.Lock()
	ret, specificReturn := fake.listV3OrganizationRolesByGUIDAndTypeReturnsOnCall[len(fake.listV3OrganizationRolesByGUIDAndTypeArgsForCall)]
	fake.listV3OrganizationRolesByGUIDAndTypeArgsForCall = append(fake.listV3OrganizationRolesByGUIDAndTypeArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ListV3OrganizationRolesByGUIDAndType", []interface{}{arg1, arg2})
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.Unlock()
	if fake.ListV3OrganizationRolesByGUIDAndTypeStub != nil {
		return fake.ListV3OrganizationRolesByGUIDAndTypeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listV3OrganizationRolesByGUIDAndTypeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndTypeCallCount() int {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.RUnlock()
	return len(fake.listV3OrganizationRolesByGUIDAndTypeArgsForCall)
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndTypeCalls(stub func(string, string) ([]cfclient.V3User, error)) {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3OrganizationRolesByGUIDAndTypeStub = stub
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndTypeArgsForCall(i int) (string, string) {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.RUnlock()
	argsForCall := fake.listV3OrganizationRolesByGUIDAndTypeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndTypeReturns(result1 []cfclient.V3User, result2 error) {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3OrganizationRolesByGUIDAndTypeStub = nil
	fake.listV3OrganizationRolesByGUIDAndTypeReturns = struct {
		result1 []cfclient.V3User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3OrganizationRolesByGUIDAndTypeReturnsOnCall(i int, result1 []cfclient.V3User, result2 error) {
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3OrganizationRolesByGUIDAndTypeStub = nil
	if fake.listV3OrganizationRolesByGUIDAndTypeReturnsOnCall == nil {
		fake.listV3OrganizationRolesByGUIDAndTypeReturnsOnCall = make(map[int]struct {
			result1 []cfclient.V3User
			result2 error
		})
	}
	fake.listV3OrganizationRolesByGUIDAndTypeReturnsOnCall[i] = struct {
		result1 []cfclient.V3User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3RolesByQuery(arg1 url.Values) ([]cfclient.V3Role, error) {
	fake.listV3RolesByQueryMutex.Lock()
	ret, specificReturn := fake.listV3RolesByQueryReturnsOnCall[len(fake.listV3RolesByQueryArgsForCall)]
	fake.listV3RolesByQueryArgsForCall = append(fake.listV3RolesByQueryArgsForCall, struct {
		arg1 url.Values
	}{arg1})
	fake.recordInvocation("ListV3RolesByQuery", []interface{}{arg1})
	fake.listV3RolesByQueryMutex.Unlock()
	if fake.ListV3RolesByQueryStub != nil {
		return fake.ListV3RolesByQueryStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listV3RolesByQueryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) ListV3RolesByQueryCallCount() int {
	fake.listV3RolesByQueryMutex.RLock()
	defer fake.listV3RolesByQueryMutex.RUnlock()
	return len(fake.listV3RolesByQueryArgsForCall)
}

func (fake *FakeCFClient) ListV3RolesByQueryCalls(stub func(url.Values) ([]cfclient.V3Role, error)) {
	fake.listV3RolesByQueryMutex.Lock()
	defer fake.listV3RolesByQueryMutex.Unlock()
	fake.ListV3RolesByQueryStub = stub
}

func (fake *FakeCFClient) ListV3RolesByQueryArgsForCall(i int) url.Values {
	fake.listV3RolesByQueryMutex.RLock()
	defer fake.listV3RolesByQueryMutex.RUnlock()
	argsForCall := fake.listV3RolesByQueryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) ListV3RolesByQueryReturns(result1 []cfclient.V3Role, result2 error) {
	fake.listV3RolesByQueryMutex.Lock()
	defer fake.listV3RolesByQueryMutex.Unlock()
	fake.ListV3RolesByQueryStub = nil
	fake.listV3RolesByQueryReturns = struct {
		result1 []cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3RolesByQueryReturnsOnCall(i int, result1 []cfclient.V3Role, result2 error) {
	fake.listV3RolesByQueryMutex.Lock()
	defer fake.listV3RolesByQueryMutex.Unlock()
	fake.ListV3RolesByQueryStub = nil
	if fake.listV3RolesByQueryReturnsOnCall == nil {
		fake.listV3RolesByQueryReturnsOnCall = make(map[int]struct {
			result1 []cfclient.V3Role
			result2 error
		})
	}
	fake.listV3RolesByQueryReturnsOnCall[i] = struct {
		result1 []cfclient.V3Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndType(arg1 string, arg2 string) ([]cfclient.V3User, error) {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.Lock()
	ret, specificReturn := fake.listV3SpaceRolesByGUIDAndTypeReturnsOnCall[len(fake.listV3SpaceRolesByGUIDAndTypeArgsForCall)]
	fake.listV3SpaceRolesByGUIDAndTypeArgsForCall = append(fake.listV3SpaceRolesByGUIDAndTypeArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ListV3SpaceRolesByGUIDAndType", []interface{}{arg1, arg2})
	fake.listV3SpaceRolesByGUIDAndTypeMutex.Unlock()
	if fake.ListV3SpaceRolesByGUIDAndTypeStub != nil {
		return fake.ListV3SpaceRolesByGUIDAndTypeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listV3SpaceRolesByGUIDAndTypeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndTypeCallCount() int {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.RUnlock()
	return len(fake.listV3SpaceRolesByGUIDAndTypeArgsForCall)
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndTypeCalls(stub func(string, string) ([]cfclient.V3User, error)) {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3SpaceRolesByGUIDAndTypeStub = stub
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndTypeArgsForCall(i int) (string, string) {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.RUnlock()
	argsForCall := fake.listV3SpaceRolesByGUIDAndTypeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndTypeReturns(result1 []cfclient.V3User, result2 error) {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3SpaceRolesByGUIDAndTypeStub = nil
	fake.listV3SpaceRolesByGUIDAndTypeReturns = struct {
		result1 []cfclient.V3User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListV3SpaceRolesByGUIDAndTypeReturnsOnCall(i int, result1 []cfclient.V3User, result2 error) {
	fake.listV3SpaceRolesByGUIDAndTypeMutex.Lock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.Unlock()
	fake.ListV3SpaceRolesByGUIDAndTypeStub = nil
	if fake.listV3SpaceRolesByGUIDAndTypeReturnsOnCall == nil {
		fake.listV3SpaceRolesByGUIDAndTypeReturnsOnCall = make(map[int]struct {
			result1 []cfclient.V3User
			result2 error
		})
	}
	fake.listV3SpaceRolesByGUIDAndTypeReturnsOnCall[i] = struct {
		result1 []cfclient.V3User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) SupportsSpaceSupporterRole() (bool, error) {
	fake.supportsSpaceSupporterRoleMutex.Lock()
	ret, specificReturn := fake.supportsSpaceSupporterRoleReturnsOnCall[len(fake.supportsSpaceSupporterRoleArgsForCall)]
	fake.supportsSpaceSupporterRoleArgsForCall = append(fake.supportsSpaceSupporterRoleArgsForCall, struct {
	}{})
	fake.recordInvocation("SupportsSpaceSupporterRole", []interface{}{})
	fake.supportsSpaceSupporterRoleMutex.Unlock()
	if fake.SupportsSpaceSupporterRoleStub != nil {
		return fake.SupportsSpaceSupporterRoleStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.supportsSpaceSupporterRoleReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) SupportsSpaceSupporterRoleCallCount() int {
	fake.supportsSpaceSupporterRoleMutex.RLock()
	defer fake.supportsSpaceSupporterRoleMutex.RUnlock()
	return len(fake.supportsSpaceSupporterRoleArgsForCall)
}

func (fake *FakeCFClient) SupportsSpaceSupporterRoleCalls(stub func() (bool, error)) {
	fake.supportsSpaceSupporterRoleMutex.Lock()
	defer fake.supportsSpaceSupporterRoleMutex.Unlock()
	fake.SupportsSpaceSupporterRoleStub = stub
}

func (fake *FakeCFClient) SupportsSpaceSupporterRoleReturns(result1 bool, result2 error) {
	fake.supportsSpaceSupporterRoleMutex.Lock()
	defer fake.supportsSpaceSupporterRoleMutex.Unlock()
	fake.SupportsSpaceSupporterRoleStub = nil
	fake.supportsSpaceSupporterRoleReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) SupportsSpaceSupporterRoleReturnsOnCall(i int, result1 bool, result2 error) {
	fake.supportsSpaceSupporterRoleMutex.Lock()
	defer fake.supportsSpaceSupporterRoleMutex.Unlock()
	fake.SupportsSpaceSupporterRoleStub = nil
	if fake.supportsSpaceSupporterRoleReturnsOnCall == nil {
		fake.supportsSpaceSupporterRoleReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.supportsSpaceSupporterRoleReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createV3OrganizationRoleMutex.RLock()
	defer fake.createV3OrganizationRoleMutex.RUnlock()
	fake.createV3SpaceRoleMutex.RLock()
	defer fake.createV3SpaceRoleMutex.RUnlock()
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	fake.deleteV3RoleMutex.RLock()
	defer fake.deleteV3RoleMutex.RUnlock()
	fake.listSpacesByQueryMutex.RLock()
	defer fake.listSpacesByQueryMutex.RUnlock()
	fake.listV3OrganizationRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3OrganizationRolesByGUIDAndTypeMutex.RUnlock()
	fake.listV3RolesByQueryMutex.RLock()
	defer fake.listV3RolesByQueryMutex.RUnlock()
	fake.listV3SpaceRolesByGUIDAndTypeMutex.RLock()
	defer fake.listV3SpaceRolesByGUIDAndTypeMutex.RUnlock()
	fake.supportsSpaceSupporterRoleMutex.RLock()
	defer fake.supportsSpaceSupporterRoleMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCFClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ user.CFClient = new(FakeCFClient)
