// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
)

type FakeCFRoleClient struct {
	CreateOrganizationRoleStub        func(context.Context, string, string, resource.OrganizationRoleType) (*resource.Role, error)
	createOrganizationRoleMutex       sync.RWMutex
	createOrganizationRoleArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 resource.OrganizationRoleType
	}
	createOrganizationRoleReturns struct {
		result1 *resource.Role
		result2 error
	}
	createOrganizationRoleReturnsOnCall map[int]struct {
		result1 *resource.Role
		result2 error
	}
	CreateSpaceRoleStub        func(context.Context, string, string, resource.SpaceRoleType) (*resource.Role, error)
	createSpaceRoleMutex       sync.RWMutex
	createSpaceRoleArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 resource.SpaceRoleType
	}
	createSpaceRoleReturns struct {
		result1 *resource.Role
		result2 error
	}
	createSpaceRoleReturnsOnCall map[int]struct {
		result1 *resource.Role
		result2 error
	}
	DeleteStub        func(context.Context, string) (string, error)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	deleteReturns struct {
		result1 string
		result2 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	ListAllStub        func(context.Context, *client.RoleListOptions) ([]*resource.Role, error)
	listAllMutex       sync.RWMutex
	listAllArgsForCall []struct {
		arg1 context.Context
		arg2 *client.RoleListOptions
	}
	listAllReturns struct {
		result1 []*resource.Role
		result2 error
	}
	listAllReturnsOnCall map[int]struct {
		result1 []*resource.Role
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFRoleClient) CreateOrganizationRole(arg1 context.Context, arg2 string, arg3 string, arg4 resource.OrganizationRoleType) (*resource.Role, error) {
	fake.createOrganizationRoleMutex.Lock()
	ret, specificReturn := fake.createOrganizationRoleReturnsOnCall[len(fake.createOrganizationRoleArgsForCall)]
	fake.createOrganizationRoleArgsForCall = append(fake.createOrganizationRoleArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 resource.OrganizationRoleType
	}{arg1, arg2, arg3, arg4})
	stub := fake.CreateOrganizationRoleStub
	fakeReturns := fake.createOrganizationRoleReturns
	fake.recordInvocation("CreateOrganizationRole", []interface{}{arg1, arg2, arg3, arg4})
	fake.createOrganizationRoleMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFRoleClient) CreateOrganizationRoleCallCount() int {
	fake.createOrganizationRoleMutex.RLock()
	defer fake.createOrganizationRoleMutex.RUnlock()
	return len(fake.createOrganizationRoleArgsForCall)
}

func (fake *FakeCFRoleClient) CreateOrganizationRoleCalls(stub func(context.Context, string, string, resource.OrganizationRoleType) (*resource.Role, error)) {
	fake.createOrganizationRoleMutex.Lock()
	defer fake.createOrganizationRoleMutex.Unlock()
	fake.CreateOrganizationRoleStub = stub
}

func (fake *FakeCFRoleClient) CreateOrganizationRoleArgsForCall(i int) (context.Context, string, string, resource.OrganizationRoleType) {
	fake.createOrganizationRoleMutex.RLock()
	defer fake.createOrganizationRoleMutex.RUnlock()
	argsForCall := fake.createOrganizationRoleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeCFRoleClient) CreateOrganizationRoleReturns(result1 *resource.Role, result2 error) {
	fake.createOrganizationRoleMutex.Lock()
	defer fake.createOrganizationRoleMutex.Unlock()
	fake.CreateOrganizationRoleStub = nil
	fake.createOrganizationRoleReturns = struct {
		result1 *resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) CreateOrganizationRoleReturnsOnCall(i int, result1 *resource.Role, result2 error) {
	fake.createOrganizationRoleMutex.Lock()
	defer fake.createOrganizationRoleMutex.Unlock()
	fake.CreateOrganizationRoleStub = nil
	if fake.createOrganizationRoleReturnsOnCall == nil {
		fake.createOrganizationRoleReturnsOnCall = make(map[int]struct {
			result1 *resource.Role
			result2 error
		})
	}
	fake.createOrganizationRoleReturnsOnCall[i] = struct {
		result1 *resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) CreateSpaceRole(arg1 context.Context, arg2 string, arg3 string, arg4 resource.SpaceRoleType) (*resource.Role, error) {
	fake.createSpaceRoleMutex.Lock()
	ret, specificReturn := fake.createSpaceRoleReturnsOnCall[len(fake.createSpaceRoleArgsForCall)]
	fake.createSpaceRoleArgsForCall = append(fake.createSpaceRoleArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 resource.SpaceRoleType
	}{arg1, arg2, arg3, arg4})
	stub := fake.CreateSpaceRoleStub
	fakeReturns := fake.createSpaceRoleReturns
	fake.recordInvocation("CreateSpaceRole", []interface{}{arg1, arg2, arg3, arg4})
	fake.createSpaceRoleMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFRoleClient) CreateSpaceRoleCallCount() int {
	fake.createSpaceRoleMutex.RLock()
	defer fake.createSpaceRoleMutex.RUnlock()
	return len(fake.createSpaceRoleArgsForCall)
}

func (fake *FakeCFRoleClient) CreateSpaceRoleCalls(stub func(context.Context, string, string, resource.SpaceRoleType) (*resource.Role, error)) {
	fake.createSpaceRoleMutex.Lock()
	defer fake.createSpaceRoleMutex.Unlock()
	fake.CreateSpaceRoleStub = stub
}

func (fake *FakeCFRoleClient) CreateSpaceRoleArgsForCall(i int) (context.Context, string, string, resource.SpaceRoleType) {
	fake.createSpaceRoleMutex.RLock()
	defer fake.createSpaceRoleMutex.RUnlock()
	argsForCall := fake.createSpaceRoleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeCFRoleClient) CreateSpaceRoleReturns(result1 *resource.Role, result2 error) {
	fake.createSpaceRoleMutex.Lock()
	defer fake.createSpaceRoleMutex.Unlock()
	fake.CreateSpaceRoleStub = nil
	fake.createSpaceRoleReturns = struct {
		result1 *resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) CreateSpaceRoleReturnsOnCall(i int, result1 *resource.Role, result2 error) {
	fake.createSpaceRoleMutex.Lock()
	defer fake.createSpaceRoleMutex.Unlock()
	fake.CreateSpaceRoleStub = nil
	if fake.createSpaceRoleReturnsOnCall == nil {
		fake.createSpaceRoleReturnsOnCall = make(map[int]struct {
			result1 *resource.Role
			result2 error
		})
	}
	fake.createSpaceRoleReturnsOnCall[i] = struct {
		result1 *resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) Delete(arg1 context.Context, arg2 string) (string, error) {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.DeleteStub
	fakeReturns := fake.deleteReturns
	fake.recordInvocation("Delete", []interface{}{arg1, arg2})
	fake.deleteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFRoleClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeCFRoleClient) DeleteCalls(stub func(context.Context, string) (string, error)) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeCFRoleClient) DeleteArgsForCall(i int) (context.Context, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFRoleClient) DeleteReturns(result1 string, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) DeleteReturnsOnCall(i int, result1 string, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) ListAll(arg1 context.Context, arg2 *client.RoleListOptions) ([]*resource.Role, error) {
	fake.listAllMutex.Lock()
	ret, specificReturn := fake.listAllReturnsOnCall[len(fake.listAllArgsForCall)]
	fake.listAllArgsForCall = append(fake.listAllArgsForCall, struct {
		arg1 context.Context
		arg2 *client.RoleListOptions
	}{arg1, arg2})
	stub := fake.ListAllStub
	fakeReturns := fake.listAllReturns
	fake.recordInvocation("ListAll", []interface{}{arg1, arg2})
	fake.listAllMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFRoleClient) ListAllCallCount() int {
	fake.listAllMutex.RLock()
	defer fake.listAllMutex.RUnlock()
	return len(fake.listAllArgsForCall)
}

func (fake *FakeCFRoleClient) ListAllCalls(stub func(context.Context, *client.RoleListOptions) ([]*resource.Role, error)) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = stub
}

func (fake *FakeCFRoleClient) ListAllArgsForCall(i int) (context.Context, *client.RoleListOptions) {
	fake.listAllMutex.RLock()
	defer fake.listAllMutex.RUnlock()
	argsForCall := fake.listAllArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFRoleClient) ListAllReturns(result1 []*resource.Role, result2 error) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = nil
	fake.listAllReturns = struct {
		result1 []*resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) ListAllReturnsOnCall(i int, result1 []*resource.Role, result2 error) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = nil
	if fake.listAllReturnsOnCall == nil {
		fake.listAllReturnsOnCall = make(map[int]struct {
			result1 []*resource.Role
			result2 error
		})
	}
	fake.listAllReturnsOnCall[i] = struct {
		result1 []*resource.Role
		result2 error
	}{result1, result2}
}

func (fake *FakeCFRoleClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createOrganizationRoleMutex.RLock()
	defer fake.createOrganizationRoleMutex.RUnlock()
	fake.createSpaceRoleMutex.RLock()
	defer fake.createSpaceRoleMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.listAllMutex.RLock()
	defer fake.listAllMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCFRoleClient) recordInvocation(key string, args []interface{}) {
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

var _ role.CFRoleClient = new(FakeCFRoleClient)
