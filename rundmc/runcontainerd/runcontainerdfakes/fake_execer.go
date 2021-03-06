// Code generated by counterfeiter. DO NOT EDIT.
package runcontainerdfakes

import (
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/rundmc/runcontainerd"
	"code.cloudfoundry.org/lager"
)

type FakeExecer struct {
	ExecStub        func(log lager.Logger, bundlePath string, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error)
	execMutex       sync.RWMutex
	execArgsForCall []struct {
		log        lager.Logger
		bundlePath string
		id         string
		spec       garden.ProcessSpec
		io         garden.ProcessIO
	}
	execReturns struct {
		result1 garden.Process
		result2 error
	}
	execReturnsOnCall map[int]struct {
		result1 garden.Process
		result2 error
	}
	AttachStub        func(log lager.Logger, bundlePath string, id string, processId string, io garden.ProcessIO) (garden.Process, error)
	attachMutex       sync.RWMutex
	attachArgsForCall []struct {
		log        lager.Logger
		bundlePath string
		id         string
		processId  string
		io         garden.ProcessIO
	}
	attachReturns struct {
		result1 garden.Process
		result2 error
	}
	attachReturnsOnCall map[int]struct {
		result1 garden.Process
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExecer) Exec(log lager.Logger, bundlePath string, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
	fake.execMutex.Lock()
	ret, specificReturn := fake.execReturnsOnCall[len(fake.execArgsForCall)]
	fake.execArgsForCall = append(fake.execArgsForCall, struct {
		log        lager.Logger
		bundlePath string
		id         string
		spec       garden.ProcessSpec
		io         garden.ProcessIO
	}{log, bundlePath, id, spec, io})
	fake.recordInvocation("Exec", []interface{}{log, bundlePath, id, spec, io})
	fake.execMutex.Unlock()
	if fake.ExecStub != nil {
		return fake.ExecStub(log, bundlePath, id, spec, io)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.execReturns.result1, fake.execReturns.result2
}

func (fake *FakeExecer) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeExecer) ExecArgsForCall(i int) (lager.Logger, string, string, garden.ProcessSpec, garden.ProcessIO) {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return fake.execArgsForCall[i].log, fake.execArgsForCall[i].bundlePath, fake.execArgsForCall[i].id, fake.execArgsForCall[i].spec, fake.execArgsForCall[i].io
}

func (fake *FakeExecer) ExecReturns(result1 garden.Process, result2 error) {
	fake.ExecStub = nil
	fake.execReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeExecer) ExecReturnsOnCall(i int, result1 garden.Process, result2 error) {
	fake.ExecStub = nil
	if fake.execReturnsOnCall == nil {
		fake.execReturnsOnCall = make(map[int]struct {
			result1 garden.Process
			result2 error
		})
	}
	fake.execReturnsOnCall[i] = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeExecer) Attach(log lager.Logger, bundlePath string, id string, processId string, io garden.ProcessIO) (garden.Process, error) {
	fake.attachMutex.Lock()
	ret, specificReturn := fake.attachReturnsOnCall[len(fake.attachArgsForCall)]
	fake.attachArgsForCall = append(fake.attachArgsForCall, struct {
		log        lager.Logger
		bundlePath string
		id         string
		processId  string
		io         garden.ProcessIO
	}{log, bundlePath, id, processId, io})
	fake.recordInvocation("Attach", []interface{}{log, bundlePath, id, processId, io})
	fake.attachMutex.Unlock()
	if fake.AttachStub != nil {
		return fake.AttachStub(log, bundlePath, id, processId, io)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.attachReturns.result1, fake.attachReturns.result2
}

func (fake *FakeExecer) AttachCallCount() int {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return len(fake.attachArgsForCall)
}

func (fake *FakeExecer) AttachArgsForCall(i int) (lager.Logger, string, string, string, garden.ProcessIO) {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return fake.attachArgsForCall[i].log, fake.attachArgsForCall[i].bundlePath, fake.attachArgsForCall[i].id, fake.attachArgsForCall[i].processId, fake.attachArgsForCall[i].io
}

func (fake *FakeExecer) AttachReturns(result1 garden.Process, result2 error) {
	fake.AttachStub = nil
	fake.attachReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeExecer) AttachReturnsOnCall(i int, result1 garden.Process, result2 error) {
	fake.AttachStub = nil
	if fake.attachReturnsOnCall == nil {
		fake.attachReturnsOnCall = make(map[int]struct {
			result1 garden.Process
			result2 error
		})
	}
	fake.attachReturnsOnCall[i] = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeExecer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeExecer) recordInvocation(key string, args []interface{}) {
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

var _ runcontainerd.Execer = new(FakeExecer)
