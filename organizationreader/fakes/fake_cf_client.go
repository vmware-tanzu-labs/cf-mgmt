// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
)

type FakeCFClient struct {
	DeleteOrgStub        func(string, bool, bool) error
	deleteOrgMutex       sync.RWMutex
	deleteOrgArgsForCall []struct {
		arg1 string
		arg2 bool
		arg3 bool
	}
	deleteOrgReturns struct {
		result1 error
	}
	deleteOrgReturnsOnCall map[int]struct {
		result1 error
	}
	GetOrgByGuidStub        func(string) (cfclient.Org, error)
	getOrgByGuidMutex       sync.RWMutex
	getOrgByGuidArgsForCall []struct {
		arg1 string
	}
	getOrgByGuidReturns struct {
		result1 cfclient.Org
		result2 error
	}
	getOrgByGuidReturnsOnCall map[int]struct {
		result1 cfclient.Org
		result2 error
	}
	ListOrgsStub        func() ([]cfclient.Org, error)
	listOrgsMutex       sync.RWMutex
	listOrgsArgsForCall []struct {
	}
	listOrgsReturns struct {
		result1 []cfclient.Org
		result2 error
	}
	listOrgsReturnsOnCall map[int]struct {
		result1 []cfclient.Org
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFClient) DeleteOrg(arg1 string, arg2 bool, arg3 bool) error {
	fake.deleteOrgMutex.Lock()
	ret, specificReturn := fake.deleteOrgReturnsOnCall[len(fake.deleteOrgArgsForCall)]
	fake.deleteOrgArgsForCall = append(fake.deleteOrgArgsForCall, struct {
		arg1 string
		arg2 bool
		arg3 bool
	}{arg1, arg2, arg3})
	stub := fake.DeleteOrgStub
	fakeReturns := fake.deleteOrgReturns
	fake.recordInvocation("DeleteOrg", []interface{}{arg1, arg2, arg3})
	fake.deleteOrgMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCFClient) DeleteOrgCallCount() int {
	fake.deleteOrgMutex.RLock()
	defer fake.deleteOrgMutex.RUnlock()
	return len(fake.deleteOrgArgsForCall)
}

func (fake *FakeCFClient) DeleteOrgCalls(stub func(string, bool, bool) error) {
	fake.deleteOrgMutex.Lock()
	defer fake.deleteOrgMutex.Unlock()
	fake.DeleteOrgStub = stub
}

func (fake *FakeCFClient) DeleteOrgArgsForCall(i int) (string, bool, bool) {
	fake.deleteOrgMutex.RLock()
	defer fake.deleteOrgMutex.RUnlock()
	argsForCall := fake.deleteOrgArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCFClient) DeleteOrgReturns(result1 error) {
	fake.deleteOrgMutex.Lock()
	defer fake.deleteOrgMutex.Unlock()
	fake.DeleteOrgStub = nil
	fake.deleteOrgReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) DeleteOrgReturnsOnCall(i int, result1 error) {
	fake.deleteOrgMutex.Lock()
	defer fake.deleteOrgMutex.Unlock()
	fake.DeleteOrgStub = nil
	if fake.deleteOrgReturnsOnCall == nil {
		fake.deleteOrgReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteOrgReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) GetOrgByGuid(arg1 string) (cfclient.Org, error) {
	fake.getOrgByGuidMutex.Lock()
	ret, specificReturn := fake.getOrgByGuidReturnsOnCall[len(fake.getOrgByGuidArgsForCall)]
	fake.getOrgByGuidArgsForCall = append(fake.getOrgByGuidArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetOrgByGuidStub
	fakeReturns := fake.getOrgByGuidReturns
	fake.recordInvocation("GetOrgByGuid", []interface{}{arg1})
	fake.getOrgByGuidMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) GetOrgByGuidCallCount() int {
	fake.getOrgByGuidMutex.RLock()
	defer fake.getOrgByGuidMutex.RUnlock()
	return len(fake.getOrgByGuidArgsForCall)
}

func (fake *FakeCFClient) GetOrgByGuidCalls(stub func(string) (cfclient.Org, error)) {
	fake.getOrgByGuidMutex.Lock()
	defer fake.getOrgByGuidMutex.Unlock()
	fake.GetOrgByGuidStub = stub
}

func (fake *FakeCFClient) GetOrgByGuidArgsForCall(i int) string {
	fake.getOrgByGuidMutex.RLock()
	defer fake.getOrgByGuidMutex.RUnlock()
	argsForCall := fake.getOrgByGuidArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) GetOrgByGuidReturns(result1 cfclient.Org, result2 error) {
	fake.getOrgByGuidMutex.Lock()
	defer fake.getOrgByGuidMutex.Unlock()
	fake.GetOrgByGuidStub = nil
	fake.getOrgByGuidReturns = struct {
		result1 cfclient.Org
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) GetOrgByGuidReturnsOnCall(i int, result1 cfclient.Org, result2 error) {
	fake.getOrgByGuidMutex.Lock()
	defer fake.getOrgByGuidMutex.Unlock()
	fake.GetOrgByGuidStub = nil
	if fake.getOrgByGuidReturnsOnCall == nil {
		fake.getOrgByGuidReturnsOnCall = make(map[int]struct {
			result1 cfclient.Org
			result2 error
		})
	}
	fake.getOrgByGuidReturnsOnCall[i] = struct {
		result1 cfclient.Org
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListOrgs() ([]cfclient.Org, error) {
	fake.listOrgsMutex.Lock()
	ret, specificReturn := fake.listOrgsReturnsOnCall[len(fake.listOrgsArgsForCall)]
	fake.listOrgsArgsForCall = append(fake.listOrgsArgsForCall, struct {
	}{})
	stub := fake.ListOrgsStub
	fakeReturns := fake.listOrgsReturns
	fake.recordInvocation("ListOrgs", []interface{}{})
	fake.listOrgsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) ListOrgsCallCount() int {
	fake.listOrgsMutex.RLock()
	defer fake.listOrgsMutex.RUnlock()
	return len(fake.listOrgsArgsForCall)
}

func (fake *FakeCFClient) ListOrgsCalls(stub func() ([]cfclient.Org, error)) {
	fake.listOrgsMutex.Lock()
	defer fake.listOrgsMutex.Unlock()
	fake.ListOrgsStub = stub
}

func (fake *FakeCFClient) ListOrgsReturns(result1 []cfclient.Org, result2 error) {
	fake.listOrgsMutex.Lock()
	defer fake.listOrgsMutex.Unlock()
	fake.ListOrgsStub = nil
	fake.listOrgsReturns = struct {
		result1 []cfclient.Org
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) ListOrgsReturnsOnCall(i int, result1 []cfclient.Org, result2 error) {
	fake.listOrgsMutex.Lock()
	defer fake.listOrgsMutex.Unlock()
	fake.ListOrgsStub = nil
	if fake.listOrgsReturnsOnCall == nil {
		fake.listOrgsReturnsOnCall = make(map[int]struct {
			result1 []cfclient.Org
			result2 error
		})
	}
	fake.listOrgsReturnsOnCall[i] = struct {
		result1 []cfclient.Org
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteOrgMutex.RLock()
	defer fake.deleteOrgMutex.RUnlock()
	fake.getOrgByGuidMutex.RLock()
	defer fake.getOrgByGuidMutex.RUnlock()
	fake.listOrgsMutex.RLock()
	defer fake.listOrgsMutex.RUnlock()
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

var _ organizationreader.CFClient = new(FakeCFClient)
