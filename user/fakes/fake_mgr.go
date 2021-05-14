// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/user"
)

type FakeManager struct {
	CleanupOrgUsersStub        func() error
	cleanupOrgUsersMutex       sync.RWMutex
	cleanupOrgUsersArgsForCall []struct {
	}
	cleanupOrgUsersReturns struct {
		result1 error
	}
	cleanupOrgUsersReturnsOnCall map[int]struct {
		result1 error
	}
	DeinitializeLdapStub        func() error
	deinitializeLdapMutex       sync.RWMutex
	deinitializeLdapArgsForCall []struct {
	}
	deinitializeLdapReturns struct {
		result1 error
	}
	deinitializeLdapReturnsOnCall map[int]struct {
		result1 error
	}
	InitializeLdapStub        func(string, string, string) error
	initializeLdapMutex       sync.RWMutex
	initializeLdapArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	initializeLdapReturns struct {
		result1 error
	}
	initializeLdapReturnsOnCall map[int]struct {
		result1 error
	}
	ListOrgAuditorsStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listOrgAuditorsMutex       sync.RWMutex
	listOrgAuditorsArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listOrgAuditorsReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listOrgAuditorsReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	ListOrgBillingManagersStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listOrgBillingManagersMutex       sync.RWMutex
	listOrgBillingManagersArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listOrgBillingManagersReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listOrgBillingManagersReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	ListOrgManagersStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listOrgManagersMutex       sync.RWMutex
	listOrgManagersArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listOrgManagersReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listOrgManagersReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	ListSpaceAuditorsStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listSpaceAuditorsMutex       sync.RWMutex
	listSpaceAuditorsArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listSpaceAuditorsReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listSpaceAuditorsReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	ListSpaceDevelopersStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listSpaceDevelopersMutex       sync.RWMutex
	listSpaceDevelopersArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listSpaceDevelopersReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listSpaceDevelopersReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	ListSpaceManagersStub        func(string, *uaa.Users) (*user.RoleUsers, error)
	listSpaceManagersMutex       sync.RWMutex
	listSpaceManagersArgsForCall []struct {
		arg1 string
		arg2 *uaa.Users
	}
	listSpaceManagersReturns struct {
		result1 *user.RoleUsers
		result2 error
	}
	listSpaceManagersReturnsOnCall map[int]struct {
		result1 *user.RoleUsers
		result2 error
	}
	UpdateOrgUsersStub        func() error
	updateOrgUsersMutex       sync.RWMutex
	updateOrgUsersArgsForCall []struct {
	}
	updateOrgUsersReturns struct {
		result1 error
	}
	updateOrgUsersReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateSpaceUsersStub        func() error
	updateSpaceUsersMutex       sync.RWMutex
	updateSpaceUsersArgsForCall []struct {
	}
	updateSpaceUsersReturns struct {
		result1 error
	}
	updateSpaceUsersReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManager) CleanupOrgUsers() error {
	fake.cleanupOrgUsersMutex.Lock()
	ret, specificReturn := fake.cleanupOrgUsersReturnsOnCall[len(fake.cleanupOrgUsersArgsForCall)]
	fake.cleanupOrgUsersArgsForCall = append(fake.cleanupOrgUsersArgsForCall, struct {
	}{})
	stub := fake.CleanupOrgUsersStub
	fakeReturns := fake.cleanupOrgUsersReturns
	fake.recordInvocation("CleanupOrgUsers", []interface{}{})
	fake.cleanupOrgUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) CleanupOrgUsersCallCount() int {
	fake.cleanupOrgUsersMutex.RLock()
	defer fake.cleanupOrgUsersMutex.RUnlock()
	return len(fake.cleanupOrgUsersArgsForCall)
}

func (fake *FakeManager) CleanupOrgUsersCalls(stub func() error) {
	fake.cleanupOrgUsersMutex.Lock()
	defer fake.cleanupOrgUsersMutex.Unlock()
	fake.CleanupOrgUsersStub = stub
}

func (fake *FakeManager) CleanupOrgUsersReturns(result1 error) {
	fake.cleanupOrgUsersMutex.Lock()
	defer fake.cleanupOrgUsersMutex.Unlock()
	fake.CleanupOrgUsersStub = nil
	fake.cleanupOrgUsersReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) CleanupOrgUsersReturnsOnCall(i int, result1 error) {
	fake.cleanupOrgUsersMutex.Lock()
	defer fake.cleanupOrgUsersMutex.Unlock()
	fake.CleanupOrgUsersStub = nil
	if fake.cleanupOrgUsersReturnsOnCall == nil {
		fake.cleanupOrgUsersReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.cleanupOrgUsersReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeinitializeLdap() error {
	fake.deinitializeLdapMutex.Lock()
	ret, specificReturn := fake.deinitializeLdapReturnsOnCall[len(fake.deinitializeLdapArgsForCall)]
	fake.deinitializeLdapArgsForCall = append(fake.deinitializeLdapArgsForCall, struct {
	}{})
	stub := fake.DeinitializeLdapStub
	fakeReturns := fake.deinitializeLdapReturns
	fake.recordInvocation("DeinitializeLdap", []interface{}{})
	fake.deinitializeLdapMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) DeinitializeLdapCallCount() int {
	fake.deinitializeLdapMutex.RLock()
	defer fake.deinitializeLdapMutex.RUnlock()
	return len(fake.deinitializeLdapArgsForCall)
}

func (fake *FakeManager) DeinitializeLdapCalls(stub func() error) {
	fake.deinitializeLdapMutex.Lock()
	defer fake.deinitializeLdapMutex.Unlock()
	fake.DeinitializeLdapStub = stub
}

func (fake *FakeManager) DeinitializeLdapReturns(result1 error) {
	fake.deinitializeLdapMutex.Lock()
	defer fake.deinitializeLdapMutex.Unlock()
	fake.DeinitializeLdapStub = nil
	fake.deinitializeLdapReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeinitializeLdapReturnsOnCall(i int, result1 error) {
	fake.deinitializeLdapMutex.Lock()
	defer fake.deinitializeLdapMutex.Unlock()
	fake.DeinitializeLdapStub = nil
	if fake.deinitializeLdapReturnsOnCall == nil {
		fake.deinitializeLdapReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deinitializeLdapReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) InitializeLdap(arg1 string, arg2 string, arg3 string) error {
	fake.initializeLdapMutex.Lock()
	ret, specificReturn := fake.initializeLdapReturnsOnCall[len(fake.initializeLdapArgsForCall)]
	fake.initializeLdapArgsForCall = append(fake.initializeLdapArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.InitializeLdapStub
	fakeReturns := fake.initializeLdapReturns
	fake.recordInvocation("InitializeLdap", []interface{}{arg1, arg2, arg3})
	fake.initializeLdapMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) InitializeLdapCallCount() int {
	fake.initializeLdapMutex.RLock()
	defer fake.initializeLdapMutex.RUnlock()
	return len(fake.initializeLdapArgsForCall)
}

func (fake *FakeManager) InitializeLdapCalls(stub func(string, string, string) error) {
	fake.initializeLdapMutex.Lock()
	defer fake.initializeLdapMutex.Unlock()
	fake.InitializeLdapStub = stub
}

func (fake *FakeManager) InitializeLdapArgsForCall(i int) (string, string, string) {
	fake.initializeLdapMutex.RLock()
	defer fake.initializeLdapMutex.RUnlock()
	argsForCall := fake.initializeLdapArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeManager) InitializeLdapReturns(result1 error) {
	fake.initializeLdapMutex.Lock()
	defer fake.initializeLdapMutex.Unlock()
	fake.InitializeLdapStub = nil
	fake.initializeLdapReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) InitializeLdapReturnsOnCall(i int, result1 error) {
	fake.initializeLdapMutex.Lock()
	defer fake.initializeLdapMutex.Unlock()
	fake.InitializeLdapStub = nil
	if fake.initializeLdapReturnsOnCall == nil {
		fake.initializeLdapReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initializeLdapReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) ListOrgAuditors(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listOrgAuditorsMutex.Lock()
	ret, specificReturn := fake.listOrgAuditorsReturnsOnCall[len(fake.listOrgAuditorsArgsForCall)]
	fake.listOrgAuditorsArgsForCall = append(fake.listOrgAuditorsArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListOrgAuditorsStub
	fakeReturns := fake.listOrgAuditorsReturns
	fake.recordInvocation("ListOrgAuditors", []interface{}{arg1, arg2})
	fake.listOrgAuditorsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListOrgAuditorsCallCount() int {
	fake.listOrgAuditorsMutex.RLock()
	defer fake.listOrgAuditorsMutex.RUnlock()
	return len(fake.listOrgAuditorsArgsForCall)
}

func (fake *FakeManager) ListOrgAuditorsCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listOrgAuditorsMutex.Lock()
	defer fake.listOrgAuditorsMutex.Unlock()
	fake.ListOrgAuditorsStub = stub
}

func (fake *FakeManager) ListOrgAuditorsArgsForCall(i int) (string, *uaa.Users) {
	fake.listOrgAuditorsMutex.RLock()
	defer fake.listOrgAuditorsMutex.RUnlock()
	argsForCall := fake.listOrgAuditorsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListOrgAuditorsReturns(result1 *user.RoleUsers, result2 error) {
	fake.listOrgAuditorsMutex.Lock()
	defer fake.listOrgAuditorsMutex.Unlock()
	fake.ListOrgAuditorsStub = nil
	fake.listOrgAuditorsReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListOrgAuditorsReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listOrgAuditorsMutex.Lock()
	defer fake.listOrgAuditorsMutex.Unlock()
	fake.ListOrgAuditorsStub = nil
	if fake.listOrgAuditorsReturnsOnCall == nil {
		fake.listOrgAuditorsReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listOrgAuditorsReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListOrgBillingManagers(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listOrgBillingManagersMutex.Lock()
	ret, specificReturn := fake.listOrgBillingManagersReturnsOnCall[len(fake.listOrgBillingManagersArgsForCall)]
	fake.listOrgBillingManagersArgsForCall = append(fake.listOrgBillingManagersArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListOrgBillingManagersStub
	fakeReturns := fake.listOrgBillingManagersReturns
	fake.recordInvocation("ListOrgBillingManagers", []interface{}{arg1, arg2})
	fake.listOrgBillingManagersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListOrgBillingManagersCallCount() int {
	fake.listOrgBillingManagersMutex.RLock()
	defer fake.listOrgBillingManagersMutex.RUnlock()
	return len(fake.listOrgBillingManagersArgsForCall)
}

func (fake *FakeManager) ListOrgBillingManagersCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listOrgBillingManagersMutex.Lock()
	defer fake.listOrgBillingManagersMutex.Unlock()
	fake.ListOrgBillingManagersStub = stub
}

func (fake *FakeManager) ListOrgBillingManagersArgsForCall(i int) (string, *uaa.Users) {
	fake.listOrgBillingManagersMutex.RLock()
	defer fake.listOrgBillingManagersMutex.RUnlock()
	argsForCall := fake.listOrgBillingManagersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListOrgBillingManagersReturns(result1 *user.RoleUsers, result2 error) {
	fake.listOrgBillingManagersMutex.Lock()
	defer fake.listOrgBillingManagersMutex.Unlock()
	fake.ListOrgBillingManagersStub = nil
	fake.listOrgBillingManagersReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListOrgBillingManagersReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listOrgBillingManagersMutex.Lock()
	defer fake.listOrgBillingManagersMutex.Unlock()
	fake.ListOrgBillingManagersStub = nil
	if fake.listOrgBillingManagersReturnsOnCall == nil {
		fake.listOrgBillingManagersReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listOrgBillingManagersReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListOrgManagers(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listOrgManagersMutex.Lock()
	ret, specificReturn := fake.listOrgManagersReturnsOnCall[len(fake.listOrgManagersArgsForCall)]
	fake.listOrgManagersArgsForCall = append(fake.listOrgManagersArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListOrgManagersStub
	fakeReturns := fake.listOrgManagersReturns
	fake.recordInvocation("ListOrgManagers", []interface{}{arg1, arg2})
	fake.listOrgManagersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListOrgManagersCallCount() int {
	fake.listOrgManagersMutex.RLock()
	defer fake.listOrgManagersMutex.RUnlock()
	return len(fake.listOrgManagersArgsForCall)
}

func (fake *FakeManager) ListOrgManagersCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listOrgManagersMutex.Lock()
	defer fake.listOrgManagersMutex.Unlock()
	fake.ListOrgManagersStub = stub
}

func (fake *FakeManager) ListOrgManagersArgsForCall(i int) (string, *uaa.Users) {
	fake.listOrgManagersMutex.RLock()
	defer fake.listOrgManagersMutex.RUnlock()
	argsForCall := fake.listOrgManagersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListOrgManagersReturns(result1 *user.RoleUsers, result2 error) {
	fake.listOrgManagersMutex.Lock()
	defer fake.listOrgManagersMutex.Unlock()
	fake.ListOrgManagersStub = nil
	fake.listOrgManagersReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListOrgManagersReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listOrgManagersMutex.Lock()
	defer fake.listOrgManagersMutex.Unlock()
	fake.ListOrgManagersStub = nil
	if fake.listOrgManagersReturnsOnCall == nil {
		fake.listOrgManagersReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listOrgManagersReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceAuditors(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listSpaceAuditorsMutex.Lock()
	ret, specificReturn := fake.listSpaceAuditorsReturnsOnCall[len(fake.listSpaceAuditorsArgsForCall)]
	fake.listSpaceAuditorsArgsForCall = append(fake.listSpaceAuditorsArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListSpaceAuditorsStub
	fakeReturns := fake.listSpaceAuditorsReturns
	fake.recordInvocation("ListSpaceAuditors", []interface{}{arg1, arg2})
	fake.listSpaceAuditorsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListSpaceAuditorsCallCount() int {
	fake.listSpaceAuditorsMutex.RLock()
	defer fake.listSpaceAuditorsMutex.RUnlock()
	return len(fake.listSpaceAuditorsArgsForCall)
}

func (fake *FakeManager) ListSpaceAuditorsCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listSpaceAuditorsMutex.Lock()
	defer fake.listSpaceAuditorsMutex.Unlock()
	fake.ListSpaceAuditorsStub = stub
}

func (fake *FakeManager) ListSpaceAuditorsArgsForCall(i int) (string, *uaa.Users) {
	fake.listSpaceAuditorsMutex.RLock()
	defer fake.listSpaceAuditorsMutex.RUnlock()
	argsForCall := fake.listSpaceAuditorsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListSpaceAuditorsReturns(result1 *user.RoleUsers, result2 error) {
	fake.listSpaceAuditorsMutex.Lock()
	defer fake.listSpaceAuditorsMutex.Unlock()
	fake.ListSpaceAuditorsStub = nil
	fake.listSpaceAuditorsReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceAuditorsReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listSpaceAuditorsMutex.Lock()
	defer fake.listSpaceAuditorsMutex.Unlock()
	fake.ListSpaceAuditorsStub = nil
	if fake.listSpaceAuditorsReturnsOnCall == nil {
		fake.listSpaceAuditorsReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listSpaceAuditorsReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceDevelopers(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listSpaceDevelopersMutex.Lock()
	ret, specificReturn := fake.listSpaceDevelopersReturnsOnCall[len(fake.listSpaceDevelopersArgsForCall)]
	fake.listSpaceDevelopersArgsForCall = append(fake.listSpaceDevelopersArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListSpaceDevelopersStub
	fakeReturns := fake.listSpaceDevelopersReturns
	fake.recordInvocation("ListSpaceDevelopers", []interface{}{arg1, arg2})
	fake.listSpaceDevelopersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListSpaceDevelopersCallCount() int {
	fake.listSpaceDevelopersMutex.RLock()
	defer fake.listSpaceDevelopersMutex.RUnlock()
	return len(fake.listSpaceDevelopersArgsForCall)
}

func (fake *FakeManager) ListSpaceDevelopersCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listSpaceDevelopersMutex.Lock()
	defer fake.listSpaceDevelopersMutex.Unlock()
	fake.ListSpaceDevelopersStub = stub
}

func (fake *FakeManager) ListSpaceDevelopersArgsForCall(i int) (string, *uaa.Users) {
	fake.listSpaceDevelopersMutex.RLock()
	defer fake.listSpaceDevelopersMutex.RUnlock()
	argsForCall := fake.listSpaceDevelopersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListSpaceDevelopersReturns(result1 *user.RoleUsers, result2 error) {
	fake.listSpaceDevelopersMutex.Lock()
	defer fake.listSpaceDevelopersMutex.Unlock()
	fake.ListSpaceDevelopersStub = nil
	fake.listSpaceDevelopersReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceDevelopersReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listSpaceDevelopersMutex.Lock()
	defer fake.listSpaceDevelopersMutex.Unlock()
	fake.ListSpaceDevelopersStub = nil
	if fake.listSpaceDevelopersReturnsOnCall == nil {
		fake.listSpaceDevelopersReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listSpaceDevelopersReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceManagers(arg1 string, arg2 *uaa.Users) (*user.RoleUsers, error) {
	fake.listSpaceManagersMutex.Lock()
	ret, specificReturn := fake.listSpaceManagersReturnsOnCall[len(fake.listSpaceManagersArgsForCall)]
	fake.listSpaceManagersArgsForCall = append(fake.listSpaceManagersArgsForCall, struct {
		arg1 string
		arg2 *uaa.Users
	}{arg1, arg2})
	stub := fake.ListSpaceManagersStub
	fakeReturns := fake.listSpaceManagersReturns
	fake.recordInvocation("ListSpaceManagers", []interface{}{arg1, arg2})
	fake.listSpaceManagersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListSpaceManagersCallCount() int {
	fake.listSpaceManagersMutex.RLock()
	defer fake.listSpaceManagersMutex.RUnlock()
	return len(fake.listSpaceManagersArgsForCall)
}

func (fake *FakeManager) ListSpaceManagersCalls(stub func(string, *uaa.Users) (*user.RoleUsers, error)) {
	fake.listSpaceManagersMutex.Lock()
	defer fake.listSpaceManagersMutex.Unlock()
	fake.ListSpaceManagersStub = stub
}

func (fake *FakeManager) ListSpaceManagersArgsForCall(i int) (string, *uaa.Users) {
	fake.listSpaceManagersMutex.RLock()
	defer fake.listSpaceManagersMutex.RUnlock()
	argsForCall := fake.listSpaceManagersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) ListSpaceManagersReturns(result1 *user.RoleUsers, result2 error) {
	fake.listSpaceManagersMutex.Lock()
	defer fake.listSpaceManagersMutex.Unlock()
	fake.ListSpaceManagersStub = nil
	fake.listSpaceManagersReturns = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaceManagersReturnsOnCall(i int, result1 *user.RoleUsers, result2 error) {
	fake.listSpaceManagersMutex.Lock()
	defer fake.listSpaceManagersMutex.Unlock()
	fake.ListSpaceManagersStub = nil
	if fake.listSpaceManagersReturnsOnCall == nil {
		fake.listSpaceManagersReturnsOnCall = make(map[int]struct {
			result1 *user.RoleUsers
			result2 error
		})
	}
	fake.listSpaceManagersReturnsOnCall[i] = struct {
		result1 *user.RoleUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) UpdateOrgUsers() error {
	fake.updateOrgUsersMutex.Lock()
	ret, specificReturn := fake.updateOrgUsersReturnsOnCall[len(fake.updateOrgUsersArgsForCall)]
	fake.updateOrgUsersArgsForCall = append(fake.updateOrgUsersArgsForCall, struct {
	}{})
	stub := fake.UpdateOrgUsersStub
	fakeReturns := fake.updateOrgUsersReturns
	fake.recordInvocation("UpdateOrgUsers", []interface{}{})
	fake.updateOrgUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) UpdateOrgUsersCallCount() int {
	fake.updateOrgUsersMutex.RLock()
	defer fake.updateOrgUsersMutex.RUnlock()
	return len(fake.updateOrgUsersArgsForCall)
}

func (fake *FakeManager) UpdateOrgUsersCalls(stub func() error) {
	fake.updateOrgUsersMutex.Lock()
	defer fake.updateOrgUsersMutex.Unlock()
	fake.UpdateOrgUsersStub = stub
}

func (fake *FakeManager) UpdateOrgUsersReturns(result1 error) {
	fake.updateOrgUsersMutex.Lock()
	defer fake.updateOrgUsersMutex.Unlock()
	fake.UpdateOrgUsersStub = nil
	fake.updateOrgUsersReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateOrgUsersReturnsOnCall(i int, result1 error) {
	fake.updateOrgUsersMutex.Lock()
	defer fake.updateOrgUsersMutex.Unlock()
	fake.UpdateOrgUsersStub = nil
	if fake.updateOrgUsersReturnsOnCall == nil {
		fake.updateOrgUsersReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateOrgUsersReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateSpaceUsers() error {
	fake.updateSpaceUsersMutex.Lock()
	ret, specificReturn := fake.updateSpaceUsersReturnsOnCall[len(fake.updateSpaceUsersArgsForCall)]
	fake.updateSpaceUsersArgsForCall = append(fake.updateSpaceUsersArgsForCall, struct {
	}{})
	stub := fake.UpdateSpaceUsersStub
	fakeReturns := fake.updateSpaceUsersReturns
	fake.recordInvocation("UpdateSpaceUsers", []interface{}{})
	fake.updateSpaceUsersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) UpdateSpaceUsersCallCount() int {
	fake.updateSpaceUsersMutex.RLock()
	defer fake.updateSpaceUsersMutex.RUnlock()
	return len(fake.updateSpaceUsersArgsForCall)
}

func (fake *FakeManager) UpdateSpaceUsersCalls(stub func() error) {
	fake.updateSpaceUsersMutex.Lock()
	defer fake.updateSpaceUsersMutex.Unlock()
	fake.UpdateSpaceUsersStub = stub
}

func (fake *FakeManager) UpdateSpaceUsersReturns(result1 error) {
	fake.updateSpaceUsersMutex.Lock()
	defer fake.updateSpaceUsersMutex.Unlock()
	fake.UpdateSpaceUsersStub = nil
	fake.updateSpaceUsersReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateSpaceUsersReturnsOnCall(i int, result1 error) {
	fake.updateSpaceUsersMutex.Lock()
	defer fake.updateSpaceUsersMutex.Unlock()
	fake.UpdateSpaceUsersStub = nil
	if fake.updateSpaceUsersReturnsOnCall == nil {
		fake.updateSpaceUsersReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateSpaceUsersReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cleanupOrgUsersMutex.RLock()
	defer fake.cleanupOrgUsersMutex.RUnlock()
	fake.deinitializeLdapMutex.RLock()
	defer fake.deinitializeLdapMutex.RUnlock()
	fake.initializeLdapMutex.RLock()
	defer fake.initializeLdapMutex.RUnlock()
	fake.listOrgAuditorsMutex.RLock()
	defer fake.listOrgAuditorsMutex.RUnlock()
	fake.listOrgBillingManagersMutex.RLock()
	defer fake.listOrgBillingManagersMutex.RUnlock()
	fake.listOrgManagersMutex.RLock()
	defer fake.listOrgManagersMutex.RUnlock()
	fake.listSpaceAuditorsMutex.RLock()
	defer fake.listSpaceAuditorsMutex.RUnlock()
	fake.listSpaceDevelopersMutex.RLock()
	defer fake.listSpaceDevelopersMutex.RUnlock()
	fake.listSpaceManagersMutex.RLock()
	defer fake.listSpaceManagersMutex.RUnlock()
	fake.updateOrgUsersMutex.RLock()
	defer fake.updateOrgUsersMutex.RUnlock()
	fake.updateSpaceUsersMutex.RLock()
	defer fake.updateSpaceUsersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManager) recordInvocation(key string, args []interface{}) {
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

var _ user.Manager = new(FakeManager)
