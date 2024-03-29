// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
)

type FakeManager struct {
	CreateSpacesStub        func() error
	createSpacesMutex       sync.RWMutex
	createSpacesArgsForCall []struct {
	}
	createSpacesReturns struct {
		result1 error
	}
	createSpacesReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteSpacesStub        func() error
	deleteSpacesMutex       sync.RWMutex
	deleteSpacesArgsForCall []struct {
	}
	deleteSpacesReturns struct {
		result1 error
	}
	deleteSpacesReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteSpacesForOrgStub        func(string, string) error
	deleteSpacesForOrgMutex       sync.RWMutex
	deleteSpacesForOrgArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteSpacesForOrgReturns struct {
		result1 error
	}
	deleteSpacesForOrgReturnsOnCall map[int]struct {
		result1 error
	}
	FindSpaceStub        func(string, string) (*resource.Space, error)
	findSpaceMutex       sync.RWMutex
	findSpaceArgsForCall []struct {
		arg1 string
		arg2 string
	}
	findSpaceReturns struct {
		result1 *resource.Space
		result2 error
	}
	findSpaceReturnsOnCall map[int]struct {
		result1 *resource.Space
		result2 error
	}
	GetSpaceIsolationSegmentGUIDStub        func(*resource.Space) (string, error)
	getSpaceIsolationSegmentGUIDMutex       sync.RWMutex
	getSpaceIsolationSegmentGUIDArgsForCall []struct {
		arg1 *resource.Space
	}
	getSpaceIsolationSegmentGUIDReturns struct {
		result1 string
		result2 error
	}
	getSpaceIsolationSegmentGUIDReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	IsSSHEnabledStub        func(*resource.Space) (bool, error)
	isSSHEnabledMutex       sync.RWMutex
	isSSHEnabledArgsForCall []struct {
		arg1 *resource.Space
	}
	isSSHEnabledReturns struct {
		result1 bool
		result2 error
	}
	isSSHEnabledReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	ListSpacesStub        func(string) ([]*resource.Space, error)
	listSpacesMutex       sync.RWMutex
	listSpacesArgsForCall []struct {
		arg1 string
	}
	listSpacesReturns struct {
		result1 []*resource.Space
		result2 error
	}
	listSpacesReturnsOnCall map[int]struct {
		result1 []*resource.Space
		result2 error
	}
	UpdateSpacesStub        func() error
	updateSpacesMutex       sync.RWMutex
	updateSpacesArgsForCall []struct {
	}
	updateSpacesReturns struct {
		result1 error
	}
	updateSpacesReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateSpacesMetadataStub        func() error
	updateSpacesMetadataMutex       sync.RWMutex
	updateSpacesMetadataArgsForCall []struct {
	}
	updateSpacesMetadataReturns struct {
		result1 error
	}
	updateSpacesMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManager) CreateSpaces() error {
	fake.createSpacesMutex.Lock()
	ret, specificReturn := fake.createSpacesReturnsOnCall[len(fake.createSpacesArgsForCall)]
	fake.createSpacesArgsForCall = append(fake.createSpacesArgsForCall, struct {
	}{})
	stub := fake.CreateSpacesStub
	fakeReturns := fake.createSpacesReturns
	fake.recordInvocation("CreateSpaces", []interface{}{})
	fake.createSpacesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) CreateSpacesCallCount() int {
	fake.createSpacesMutex.RLock()
	defer fake.createSpacesMutex.RUnlock()
	return len(fake.createSpacesArgsForCall)
}

func (fake *FakeManager) CreateSpacesCalls(stub func() error) {
	fake.createSpacesMutex.Lock()
	defer fake.createSpacesMutex.Unlock()
	fake.CreateSpacesStub = stub
}

func (fake *FakeManager) CreateSpacesReturns(result1 error) {
	fake.createSpacesMutex.Lock()
	defer fake.createSpacesMutex.Unlock()
	fake.CreateSpacesStub = nil
	fake.createSpacesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) CreateSpacesReturnsOnCall(i int, result1 error) {
	fake.createSpacesMutex.Lock()
	defer fake.createSpacesMutex.Unlock()
	fake.CreateSpacesStub = nil
	if fake.createSpacesReturnsOnCall == nil {
		fake.createSpacesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createSpacesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeleteSpaces() error {
	fake.deleteSpacesMutex.Lock()
	ret, specificReturn := fake.deleteSpacesReturnsOnCall[len(fake.deleteSpacesArgsForCall)]
	fake.deleteSpacesArgsForCall = append(fake.deleteSpacesArgsForCall, struct {
	}{})
	stub := fake.DeleteSpacesStub
	fakeReturns := fake.deleteSpacesReturns
	fake.recordInvocation("DeleteSpaces", []interface{}{})
	fake.deleteSpacesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) DeleteSpacesCallCount() int {
	fake.deleteSpacesMutex.RLock()
	defer fake.deleteSpacesMutex.RUnlock()
	return len(fake.deleteSpacesArgsForCall)
}

func (fake *FakeManager) DeleteSpacesCalls(stub func() error) {
	fake.deleteSpacesMutex.Lock()
	defer fake.deleteSpacesMutex.Unlock()
	fake.DeleteSpacesStub = stub
}

func (fake *FakeManager) DeleteSpacesReturns(result1 error) {
	fake.deleteSpacesMutex.Lock()
	defer fake.deleteSpacesMutex.Unlock()
	fake.DeleteSpacesStub = nil
	fake.deleteSpacesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeleteSpacesReturnsOnCall(i int, result1 error) {
	fake.deleteSpacesMutex.Lock()
	defer fake.deleteSpacesMutex.Unlock()
	fake.DeleteSpacesStub = nil
	if fake.deleteSpacesReturnsOnCall == nil {
		fake.deleteSpacesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteSpacesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeleteSpacesForOrg(arg1 string, arg2 string) error {
	fake.deleteSpacesForOrgMutex.Lock()
	ret, specificReturn := fake.deleteSpacesForOrgReturnsOnCall[len(fake.deleteSpacesForOrgArgsForCall)]
	fake.deleteSpacesForOrgArgsForCall = append(fake.deleteSpacesForOrgArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.DeleteSpacesForOrgStub
	fakeReturns := fake.deleteSpacesForOrgReturns
	fake.recordInvocation("DeleteSpacesForOrg", []interface{}{arg1, arg2})
	fake.deleteSpacesForOrgMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) DeleteSpacesForOrgCallCount() int {
	fake.deleteSpacesForOrgMutex.RLock()
	defer fake.deleteSpacesForOrgMutex.RUnlock()
	return len(fake.deleteSpacesForOrgArgsForCall)
}

func (fake *FakeManager) DeleteSpacesForOrgCalls(stub func(string, string) error) {
	fake.deleteSpacesForOrgMutex.Lock()
	defer fake.deleteSpacesForOrgMutex.Unlock()
	fake.DeleteSpacesForOrgStub = stub
}

func (fake *FakeManager) DeleteSpacesForOrgArgsForCall(i int) (string, string) {
	fake.deleteSpacesForOrgMutex.RLock()
	defer fake.deleteSpacesForOrgMutex.RUnlock()
	argsForCall := fake.deleteSpacesForOrgArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) DeleteSpacesForOrgReturns(result1 error) {
	fake.deleteSpacesForOrgMutex.Lock()
	defer fake.deleteSpacesForOrgMutex.Unlock()
	fake.DeleteSpacesForOrgStub = nil
	fake.deleteSpacesForOrgReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeleteSpacesForOrgReturnsOnCall(i int, result1 error) {
	fake.deleteSpacesForOrgMutex.Lock()
	defer fake.deleteSpacesForOrgMutex.Unlock()
	fake.DeleteSpacesForOrgStub = nil
	if fake.deleteSpacesForOrgReturnsOnCall == nil {
		fake.deleteSpacesForOrgReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteSpacesForOrgReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) FindSpace(arg1 string, arg2 string) (*resource.Space, error) {
	fake.findSpaceMutex.Lock()
	ret, specificReturn := fake.findSpaceReturnsOnCall[len(fake.findSpaceArgsForCall)]
	fake.findSpaceArgsForCall = append(fake.findSpaceArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.FindSpaceStub
	fakeReturns := fake.findSpaceReturns
	fake.recordInvocation("FindSpace", []interface{}{arg1, arg2})
	fake.findSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) FindSpaceCallCount() int {
	fake.findSpaceMutex.RLock()
	defer fake.findSpaceMutex.RUnlock()
	return len(fake.findSpaceArgsForCall)
}

func (fake *FakeManager) FindSpaceCalls(stub func(string, string) (*resource.Space, error)) {
	fake.findSpaceMutex.Lock()
	defer fake.findSpaceMutex.Unlock()
	fake.FindSpaceStub = stub
}

func (fake *FakeManager) FindSpaceArgsForCall(i int) (string, string) {
	fake.findSpaceMutex.RLock()
	defer fake.findSpaceMutex.RUnlock()
	argsForCall := fake.findSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) FindSpaceReturns(result1 *resource.Space, result2 error) {
	fake.findSpaceMutex.Lock()
	defer fake.findSpaceMutex.Unlock()
	fake.FindSpaceStub = nil
	fake.findSpaceReturns = struct {
		result1 *resource.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) FindSpaceReturnsOnCall(i int, result1 *resource.Space, result2 error) {
	fake.findSpaceMutex.Lock()
	defer fake.findSpaceMutex.Unlock()
	fake.FindSpaceStub = nil
	if fake.findSpaceReturnsOnCall == nil {
		fake.findSpaceReturnsOnCall = make(map[int]struct {
			result1 *resource.Space
			result2 error
		})
	}
	fake.findSpaceReturnsOnCall[i] = struct {
		result1 *resource.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUID(arg1 *resource.Space) (string, error) {
	fake.getSpaceIsolationSegmentGUIDMutex.Lock()
	ret, specificReturn := fake.getSpaceIsolationSegmentGUIDReturnsOnCall[len(fake.getSpaceIsolationSegmentGUIDArgsForCall)]
	fake.getSpaceIsolationSegmentGUIDArgsForCall = append(fake.getSpaceIsolationSegmentGUIDArgsForCall, struct {
		arg1 *resource.Space
	}{arg1})
	stub := fake.GetSpaceIsolationSegmentGUIDStub
	fakeReturns := fake.getSpaceIsolationSegmentGUIDReturns
	fake.recordInvocation("GetSpaceIsolationSegmentGUID", []interface{}{arg1})
	fake.getSpaceIsolationSegmentGUIDMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUIDCallCount() int {
	fake.getSpaceIsolationSegmentGUIDMutex.RLock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.RUnlock()
	return len(fake.getSpaceIsolationSegmentGUIDArgsForCall)
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUIDCalls(stub func(*resource.Space) (string, error)) {
	fake.getSpaceIsolationSegmentGUIDMutex.Lock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.Unlock()
	fake.GetSpaceIsolationSegmentGUIDStub = stub
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUIDArgsForCall(i int) *resource.Space {
	fake.getSpaceIsolationSegmentGUIDMutex.RLock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.RUnlock()
	argsForCall := fake.getSpaceIsolationSegmentGUIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUIDReturns(result1 string, result2 error) {
	fake.getSpaceIsolationSegmentGUIDMutex.Lock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.Unlock()
	fake.GetSpaceIsolationSegmentGUIDStub = nil
	fake.getSpaceIsolationSegmentGUIDReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetSpaceIsolationSegmentGUIDReturnsOnCall(i int, result1 string, result2 error) {
	fake.getSpaceIsolationSegmentGUIDMutex.Lock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.Unlock()
	fake.GetSpaceIsolationSegmentGUIDStub = nil
	if fake.getSpaceIsolationSegmentGUIDReturnsOnCall == nil {
		fake.getSpaceIsolationSegmentGUIDReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getSpaceIsolationSegmentGUIDReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) IsSSHEnabled(arg1 *resource.Space) (bool, error) {
	fake.isSSHEnabledMutex.Lock()
	ret, specificReturn := fake.isSSHEnabledReturnsOnCall[len(fake.isSSHEnabledArgsForCall)]
	fake.isSSHEnabledArgsForCall = append(fake.isSSHEnabledArgsForCall, struct {
		arg1 *resource.Space
	}{arg1})
	stub := fake.IsSSHEnabledStub
	fakeReturns := fake.isSSHEnabledReturns
	fake.recordInvocation("IsSSHEnabled", []interface{}{arg1})
	fake.isSSHEnabledMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) IsSSHEnabledCallCount() int {
	fake.isSSHEnabledMutex.RLock()
	defer fake.isSSHEnabledMutex.RUnlock()
	return len(fake.isSSHEnabledArgsForCall)
}

func (fake *FakeManager) IsSSHEnabledCalls(stub func(*resource.Space) (bool, error)) {
	fake.isSSHEnabledMutex.Lock()
	defer fake.isSSHEnabledMutex.Unlock()
	fake.IsSSHEnabledStub = stub
}

func (fake *FakeManager) IsSSHEnabledArgsForCall(i int) *resource.Space {
	fake.isSSHEnabledMutex.RLock()
	defer fake.isSSHEnabledMutex.RUnlock()
	argsForCall := fake.isSSHEnabledArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) IsSSHEnabledReturns(result1 bool, result2 error) {
	fake.isSSHEnabledMutex.Lock()
	defer fake.isSSHEnabledMutex.Unlock()
	fake.IsSSHEnabledStub = nil
	fake.isSSHEnabledReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) IsSSHEnabledReturnsOnCall(i int, result1 bool, result2 error) {
	fake.isSSHEnabledMutex.Lock()
	defer fake.isSSHEnabledMutex.Unlock()
	fake.IsSSHEnabledStub = nil
	if fake.isSSHEnabledReturnsOnCall == nil {
		fake.isSSHEnabledReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.isSSHEnabledReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpaces(arg1 string) ([]*resource.Space, error) {
	fake.listSpacesMutex.Lock()
	ret, specificReturn := fake.listSpacesReturnsOnCall[len(fake.listSpacesArgsForCall)]
	fake.listSpacesArgsForCall = append(fake.listSpacesArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ListSpacesStub
	fakeReturns := fake.listSpacesReturns
	fake.recordInvocation("ListSpaces", []interface{}{arg1})
	fake.listSpacesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) ListSpacesCallCount() int {
	fake.listSpacesMutex.RLock()
	defer fake.listSpacesMutex.RUnlock()
	return len(fake.listSpacesArgsForCall)
}

func (fake *FakeManager) ListSpacesCalls(stub func(string) ([]*resource.Space, error)) {
	fake.listSpacesMutex.Lock()
	defer fake.listSpacesMutex.Unlock()
	fake.ListSpacesStub = stub
}

func (fake *FakeManager) ListSpacesArgsForCall(i int) string {
	fake.listSpacesMutex.RLock()
	defer fake.listSpacesMutex.RUnlock()
	argsForCall := fake.listSpacesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) ListSpacesReturns(result1 []*resource.Space, result2 error) {
	fake.listSpacesMutex.Lock()
	defer fake.listSpacesMutex.Unlock()
	fake.ListSpacesStub = nil
	fake.listSpacesReturns = struct {
		result1 []*resource.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) ListSpacesReturnsOnCall(i int, result1 []*resource.Space, result2 error) {
	fake.listSpacesMutex.Lock()
	defer fake.listSpacesMutex.Unlock()
	fake.ListSpacesStub = nil
	if fake.listSpacesReturnsOnCall == nil {
		fake.listSpacesReturnsOnCall = make(map[int]struct {
			result1 []*resource.Space
			result2 error
		})
	}
	fake.listSpacesReturnsOnCall[i] = struct {
		result1 []*resource.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) UpdateSpaces() error {
	fake.updateSpacesMutex.Lock()
	ret, specificReturn := fake.updateSpacesReturnsOnCall[len(fake.updateSpacesArgsForCall)]
	fake.updateSpacesArgsForCall = append(fake.updateSpacesArgsForCall, struct {
	}{})
	stub := fake.UpdateSpacesStub
	fakeReturns := fake.updateSpacesReturns
	fake.recordInvocation("UpdateSpaces", []interface{}{})
	fake.updateSpacesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) UpdateSpacesCallCount() int {
	fake.updateSpacesMutex.RLock()
	defer fake.updateSpacesMutex.RUnlock()
	return len(fake.updateSpacesArgsForCall)
}

func (fake *FakeManager) UpdateSpacesCalls(stub func() error) {
	fake.updateSpacesMutex.Lock()
	defer fake.updateSpacesMutex.Unlock()
	fake.UpdateSpacesStub = stub
}

func (fake *FakeManager) UpdateSpacesReturns(result1 error) {
	fake.updateSpacesMutex.Lock()
	defer fake.updateSpacesMutex.Unlock()
	fake.UpdateSpacesStub = nil
	fake.updateSpacesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateSpacesReturnsOnCall(i int, result1 error) {
	fake.updateSpacesMutex.Lock()
	defer fake.updateSpacesMutex.Unlock()
	fake.UpdateSpacesStub = nil
	if fake.updateSpacesReturnsOnCall == nil {
		fake.updateSpacesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateSpacesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateSpacesMetadata() error {
	fake.updateSpacesMetadataMutex.Lock()
	ret, specificReturn := fake.updateSpacesMetadataReturnsOnCall[len(fake.updateSpacesMetadataArgsForCall)]
	fake.updateSpacesMetadataArgsForCall = append(fake.updateSpacesMetadataArgsForCall, struct {
	}{})
	stub := fake.UpdateSpacesMetadataStub
	fakeReturns := fake.updateSpacesMetadataReturns
	fake.recordInvocation("UpdateSpacesMetadata", []interface{}{})
	fake.updateSpacesMetadataMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManager) UpdateSpacesMetadataCallCount() int {
	fake.updateSpacesMetadataMutex.RLock()
	defer fake.updateSpacesMetadataMutex.RUnlock()
	return len(fake.updateSpacesMetadataArgsForCall)
}

func (fake *FakeManager) UpdateSpacesMetadataCalls(stub func() error) {
	fake.updateSpacesMetadataMutex.Lock()
	defer fake.updateSpacesMetadataMutex.Unlock()
	fake.UpdateSpacesMetadataStub = stub
}

func (fake *FakeManager) UpdateSpacesMetadataReturns(result1 error) {
	fake.updateSpacesMetadataMutex.Lock()
	defer fake.updateSpacesMetadataMutex.Unlock()
	fake.UpdateSpacesMetadataStub = nil
	fake.updateSpacesMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UpdateSpacesMetadataReturnsOnCall(i int, result1 error) {
	fake.updateSpacesMetadataMutex.Lock()
	defer fake.updateSpacesMetadataMutex.Unlock()
	fake.UpdateSpacesMetadataStub = nil
	if fake.updateSpacesMetadataReturnsOnCall == nil {
		fake.updateSpacesMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateSpacesMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createSpacesMutex.RLock()
	defer fake.createSpacesMutex.RUnlock()
	fake.deleteSpacesMutex.RLock()
	defer fake.deleteSpacesMutex.RUnlock()
	fake.deleteSpacesForOrgMutex.RLock()
	defer fake.deleteSpacesForOrgMutex.RUnlock()
	fake.findSpaceMutex.RLock()
	defer fake.findSpaceMutex.RUnlock()
	fake.getSpaceIsolationSegmentGUIDMutex.RLock()
	defer fake.getSpaceIsolationSegmentGUIDMutex.RUnlock()
	fake.isSSHEnabledMutex.RLock()
	defer fake.isSSHEnabledMutex.RUnlock()
	fake.listSpacesMutex.RLock()
	defer fake.listSpacesMutex.RUnlock()
	fake.updateSpacesMutex.RLock()
	defer fake.updateSpacesMutex.RUnlock()
	fake.updateSpacesMetadataMutex.RLock()
	defer fake.updateSpacesMetadataMutex.RUnlock()
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

var _ space.Manager = new(FakeManager)
