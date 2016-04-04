// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/properties"
)

type FakeMapPersister struct {
	LoadMapStub        func(string) (map[string]string, error)
	loadMapMutex       sync.RWMutex
	loadMapArgsForCall []struct {
		arg1 string
	}
	loadMapReturns struct {
		result1 map[string]string
		result2 error
	}
	SaveMapStub        func(string, map[string]string) error
	saveMapMutex       sync.RWMutex
	saveMapArgsForCall []struct {
		arg1 string
		arg2 map[string]string
	}
	saveMapReturns struct {
		result1 error
	}
	DeleteMapStub        func(string) error
	deleteMapMutex       sync.RWMutex
	deleteMapArgsForCall []struct {
		arg1 string
	}
	deleteMapReturns struct {
		result1 error
	}
	IsMapPersistedStub        func(string) bool
	isMapPersistedMutex       sync.RWMutex
	isMapPersistedArgsForCall []struct {
		arg1 string
	}
	isMapPersistedReturns struct {
		result1 bool
	}
}

func (fake *FakeMapPersister) LoadMap(arg1 string) (map[string]string, error) {
	fake.loadMapMutex.Lock()
	fake.loadMapArgsForCall = append(fake.loadMapArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.loadMapMutex.Unlock()
	if fake.LoadMapStub != nil {
		return fake.LoadMapStub(arg1)
	} else {
		return fake.loadMapReturns.result1, fake.loadMapReturns.result2
	}
}

func (fake *FakeMapPersister) LoadMapCallCount() int {
	fake.loadMapMutex.RLock()
	defer fake.loadMapMutex.RUnlock()
	return len(fake.loadMapArgsForCall)
}

func (fake *FakeMapPersister) LoadMapArgsForCall(i int) string {
	fake.loadMapMutex.RLock()
	defer fake.loadMapMutex.RUnlock()
	return fake.loadMapArgsForCall[i].arg1
}

func (fake *FakeMapPersister) LoadMapReturns(result1 map[string]string, result2 error) {
	fake.LoadMapStub = nil
	fake.loadMapReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeMapPersister) SaveMap(arg1 string, arg2 map[string]string) error {
	fake.saveMapMutex.Lock()
	fake.saveMapArgsForCall = append(fake.saveMapArgsForCall, struct {
		arg1 string
		arg2 map[string]string
	}{arg1, arg2})
	fake.saveMapMutex.Unlock()
	if fake.SaveMapStub != nil {
		return fake.SaveMapStub(arg1, arg2)
	} else {
		return fake.saveMapReturns.result1
	}
}

func (fake *FakeMapPersister) SaveMapCallCount() int {
	fake.saveMapMutex.RLock()
	defer fake.saveMapMutex.RUnlock()
	return len(fake.saveMapArgsForCall)
}

func (fake *FakeMapPersister) SaveMapArgsForCall(i int) (string, map[string]string) {
	fake.saveMapMutex.RLock()
	defer fake.saveMapMutex.RUnlock()
	return fake.saveMapArgsForCall[i].arg1, fake.saveMapArgsForCall[i].arg2
}

func (fake *FakeMapPersister) SaveMapReturns(result1 error) {
	fake.SaveMapStub = nil
	fake.saveMapReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMapPersister) DeleteMap(arg1 string) error {
	fake.deleteMapMutex.Lock()
	fake.deleteMapArgsForCall = append(fake.deleteMapArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.deleteMapMutex.Unlock()
	if fake.DeleteMapStub != nil {
		return fake.DeleteMapStub(arg1)
	} else {
		return fake.deleteMapReturns.result1
	}
}

func (fake *FakeMapPersister) DeleteMapCallCount() int {
	fake.deleteMapMutex.RLock()
	defer fake.deleteMapMutex.RUnlock()
	return len(fake.deleteMapArgsForCall)
}

func (fake *FakeMapPersister) DeleteMapArgsForCall(i int) string {
	fake.deleteMapMutex.RLock()
	defer fake.deleteMapMutex.RUnlock()
	return fake.deleteMapArgsForCall[i].arg1
}

func (fake *FakeMapPersister) DeleteMapReturns(result1 error) {
	fake.DeleteMapStub = nil
	fake.deleteMapReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMapPersister) IsMapPersisted(arg1 string) bool {
	fake.isMapPersistedMutex.Lock()
	fake.isMapPersistedArgsForCall = append(fake.isMapPersistedArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.isMapPersistedMutex.Unlock()
	if fake.IsMapPersistedStub != nil {
		return fake.IsMapPersistedStub(arg1)
	} else {
		return fake.isMapPersistedReturns.result1
	}
}

func (fake *FakeMapPersister) IsMapPersistedCallCount() int {
	fake.isMapPersistedMutex.RLock()
	defer fake.isMapPersistedMutex.RUnlock()
	return len(fake.isMapPersistedArgsForCall)
}

func (fake *FakeMapPersister) IsMapPersistedArgsForCall(i int) string {
	fake.isMapPersistedMutex.RLock()
	defer fake.isMapPersistedMutex.RUnlock()
	return fake.isMapPersistedArgsForCall[i].arg1
}

func (fake *FakeMapPersister) IsMapPersistedReturns(result1 bool) {
	fake.IsMapPersistedStub = nil
	fake.isMapPersistedReturns = struct {
		result1 bool
	}{result1}
}

var _ properties.MapPersister = new(FakeMapPersister)
