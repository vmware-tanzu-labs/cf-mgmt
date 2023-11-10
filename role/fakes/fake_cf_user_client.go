// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
)

type FakeCFUserClient struct {
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
	ListAllStub        func(context.Context, *client.UserListOptions) ([]*resource.User, error)
	listAllMutex       sync.RWMutex
	listAllArgsForCall []struct {
		arg1 context.Context
		arg2 *client.UserListOptions
	}
	listAllReturns struct {
		result1 []*resource.User
		result2 error
	}
	listAllReturnsOnCall map[int]struct {
		result1 []*resource.User
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFUserClient) Delete(arg1 context.Context, arg2 string) (string, error) {
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

func (fake *FakeCFUserClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeCFUserClient) DeleteCalls(stub func(context.Context, string) (string, error)) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeCFUserClient) DeleteArgsForCall(i int) (context.Context, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFUserClient) DeleteReturns(result1 string, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCFUserClient) DeleteReturnsOnCall(i int, result1 string, result2 error) {
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

func (fake *FakeCFUserClient) ListAll(arg1 context.Context, arg2 *client.UserListOptions) ([]*resource.User, error) {
	fake.listAllMutex.Lock()
	ret, specificReturn := fake.listAllReturnsOnCall[len(fake.listAllArgsForCall)]
	fake.listAllArgsForCall = append(fake.listAllArgsForCall, struct {
		arg1 context.Context
		arg2 *client.UserListOptions
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

func (fake *FakeCFUserClient) ListAllCallCount() int {
	fake.listAllMutex.RLock()
	defer fake.listAllMutex.RUnlock()
	return len(fake.listAllArgsForCall)
}

func (fake *FakeCFUserClient) ListAllCalls(stub func(context.Context, *client.UserListOptions) ([]*resource.User, error)) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = stub
}

func (fake *FakeCFUserClient) ListAllArgsForCall(i int) (context.Context, *client.UserListOptions) {
	fake.listAllMutex.RLock()
	defer fake.listAllMutex.RUnlock()
	argsForCall := fake.listAllArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFUserClient) ListAllReturns(result1 []*resource.User, result2 error) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = nil
	fake.listAllReturns = struct {
		result1 []*resource.User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFUserClient) ListAllReturnsOnCall(i int, result1 []*resource.User, result2 error) {
	fake.listAllMutex.Lock()
	defer fake.listAllMutex.Unlock()
	fake.ListAllStub = nil
	if fake.listAllReturnsOnCall == nil {
		fake.listAllReturnsOnCall = make(map[int]struct {
			result1 []*resource.User
			result2 error
		})
	}
	fake.listAllReturnsOnCall[i] = struct {
		result1 []*resource.User
		result2 error
	}{result1, result2}
}

func (fake *FakeCFUserClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
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

func (fake *FakeCFUserClient) recordInvocation(key string, args []interface{}) {
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

var _ role.CFUserClient = new(FakeCFUserClient)
