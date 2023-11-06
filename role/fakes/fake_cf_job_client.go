// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
)

type FakeCFJobClient struct {
	PollCompleteStub        func(context.Context, string, *client.PollingOptions) error
	pollCompleteMutex       sync.RWMutex
	pollCompleteArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 *client.PollingOptions
	}
	pollCompleteReturns struct {
		result1 error
	}
	pollCompleteReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFJobClient) PollComplete(arg1 context.Context, arg2 string, arg3 *client.PollingOptions) error {
	fake.pollCompleteMutex.Lock()
	ret, specificReturn := fake.pollCompleteReturnsOnCall[len(fake.pollCompleteArgsForCall)]
	fake.pollCompleteArgsForCall = append(fake.pollCompleteArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 *client.PollingOptions
	}{arg1, arg2, arg3})
	stub := fake.PollCompleteStub
	fakeReturns := fake.pollCompleteReturns
	fake.recordInvocation("PollComplete", []interface{}{arg1, arg2, arg3})
	fake.pollCompleteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCFJobClient) PollCompleteCallCount() int {
	fake.pollCompleteMutex.RLock()
	defer fake.pollCompleteMutex.RUnlock()
	return len(fake.pollCompleteArgsForCall)
}

func (fake *FakeCFJobClient) PollCompleteCalls(stub func(context.Context, string, *client.PollingOptions) error) {
	fake.pollCompleteMutex.Lock()
	defer fake.pollCompleteMutex.Unlock()
	fake.PollCompleteStub = stub
}

func (fake *FakeCFJobClient) PollCompleteArgsForCall(i int) (context.Context, string, *client.PollingOptions) {
	fake.pollCompleteMutex.RLock()
	defer fake.pollCompleteMutex.RUnlock()
	argsForCall := fake.pollCompleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCFJobClient) PollCompleteReturns(result1 error) {
	fake.pollCompleteMutex.Lock()
	defer fake.pollCompleteMutex.Unlock()
	fake.PollCompleteStub = nil
	fake.pollCompleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFJobClient) PollCompleteReturnsOnCall(i int, result1 error) {
	fake.pollCompleteMutex.Lock()
	defer fake.pollCompleteMutex.Unlock()
	fake.PollCompleteStub = nil
	if fake.pollCompleteReturnsOnCall == nil {
		fake.pollCompleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pollCompleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFJobClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pollCompleteMutex.RLock()
	defer fake.pollCompleteMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCFJobClient) recordInvocation(key string, args []interface{}) {
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

var _ role.CFJobClient = new(FakeCFJobClient)
