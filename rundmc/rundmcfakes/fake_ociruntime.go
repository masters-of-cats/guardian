// Code generated by counterfeiter. DO NOT EDIT.
package rundmcfakes

import (
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gardener"
	"code.cloudfoundry.org/guardian/rundmc"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	"code.cloudfoundry.org/lager"
)

type FakeOCIRuntime struct {
	CreateStub        func(log lager.Logger, bundlePath, id string, io garden.ProcessIO) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		log        lager.Logger
		bundlePath string
		id         string
		io         garden.ProcessIO
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	ExecStub        func(log lager.Logger, bundlePath, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error)
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
	AttachStub        func(log lager.Logger, bundlePath, id, processId string, io garden.ProcessIO) (garden.Process, error)
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
	KillStub        func(log lager.Logger, bundlePath string) error
	killMutex       sync.RWMutex
	killArgsForCall []struct {
		log        lager.Logger
		bundlePath string
	}
	killReturns struct {
		result1 error
	}
	killReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(log lager.Logger, force bool, bundlePath string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		log        lager.Logger
		force      bool
		bundlePath string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	StateStub        func(log lager.Logger, id string) (runrunc.State, error)
	stateMutex       sync.RWMutex
	stateArgsForCall []struct {
		log lager.Logger
		id  string
	}
	stateReturns struct {
		result1 runrunc.State
		result2 error
	}
	stateReturnsOnCall map[int]struct {
		result1 runrunc.State
		result2 error
	}
	StatsStub        func(log lager.Logger, id string) (gardener.ActualContainerMetrics, error)
	statsMutex       sync.RWMutex
	statsArgsForCall []struct {
		log lager.Logger
		id  string
	}
	statsReturns struct {
		result1 gardener.ActualContainerMetrics
		result2 error
	}
	statsReturnsOnCall map[int]struct {
		result1 gardener.ActualContainerMetrics
		result2 error
	}
	WatchEventsStub        func(log lager.Logger, id string, eventsNotifier runrunc.EventsNotifier) error
	watchEventsMutex       sync.RWMutex
	watchEventsArgsForCall []struct {
		log            lager.Logger
		id             string
		eventsNotifier runrunc.EventsNotifier
	}
	watchEventsReturns struct {
		result1 error
	}
	watchEventsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOCIRuntime) Create(log lager.Logger, bundlePath string, id string, io garden.ProcessIO) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		log        lager.Logger
		bundlePath string
		id         string
		io         garden.ProcessIO
	}{log, bundlePath, id, io})
	fake.recordInvocation("Create", []interface{}{log, bundlePath, id, io})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(log, bundlePath, id, io)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createReturns.result1
}

func (fake *FakeOCIRuntime) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeOCIRuntime) CreateArgsForCall(i int) (lager.Logger, string, string, garden.ProcessIO) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].log, fake.createArgsForCall[i].bundlePath, fake.createArgsForCall[i].id, fake.createArgsForCall[i].io
}

func (fake *FakeOCIRuntime) CreateReturns(result1 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) CreateReturnsOnCall(i int, result1 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) Exec(log lager.Logger, bundlePath string, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
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

func (fake *FakeOCIRuntime) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeOCIRuntime) ExecArgsForCall(i int) (lager.Logger, string, string, garden.ProcessSpec, garden.ProcessIO) {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return fake.execArgsForCall[i].log, fake.execArgsForCall[i].bundlePath, fake.execArgsForCall[i].id, fake.execArgsForCall[i].spec, fake.execArgsForCall[i].io
}

func (fake *FakeOCIRuntime) ExecReturns(result1 garden.Process, result2 error) {
	fake.ExecStub = nil
	fake.execReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) ExecReturnsOnCall(i int, result1 garden.Process, result2 error) {
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

func (fake *FakeOCIRuntime) Attach(log lager.Logger, bundlePath string, id string, processId string, io garden.ProcessIO) (garden.Process, error) {
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

func (fake *FakeOCIRuntime) AttachCallCount() int {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return len(fake.attachArgsForCall)
}

func (fake *FakeOCIRuntime) AttachArgsForCall(i int) (lager.Logger, string, string, string, garden.ProcessIO) {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return fake.attachArgsForCall[i].log, fake.attachArgsForCall[i].bundlePath, fake.attachArgsForCall[i].id, fake.attachArgsForCall[i].processId, fake.attachArgsForCall[i].io
}

func (fake *FakeOCIRuntime) AttachReturns(result1 garden.Process, result2 error) {
	fake.AttachStub = nil
	fake.attachReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) AttachReturnsOnCall(i int, result1 garden.Process, result2 error) {
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

func (fake *FakeOCIRuntime) Kill(log lager.Logger, bundlePath string) error {
	fake.killMutex.Lock()
	ret, specificReturn := fake.killReturnsOnCall[len(fake.killArgsForCall)]
	fake.killArgsForCall = append(fake.killArgsForCall, struct {
		log        lager.Logger
		bundlePath string
	}{log, bundlePath})
	fake.recordInvocation("Kill", []interface{}{log, bundlePath})
	fake.killMutex.Unlock()
	if fake.KillStub != nil {
		return fake.KillStub(log, bundlePath)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.killReturns.result1
}

func (fake *FakeOCIRuntime) KillCallCount() int {
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	return len(fake.killArgsForCall)
}

func (fake *FakeOCIRuntime) KillArgsForCall(i int) (lager.Logger, string) {
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	return fake.killArgsForCall[i].log, fake.killArgsForCall[i].bundlePath
}

func (fake *FakeOCIRuntime) KillReturns(result1 error) {
	fake.KillStub = nil
	fake.killReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) KillReturnsOnCall(i int, result1 error) {
	fake.KillStub = nil
	if fake.killReturnsOnCall == nil {
		fake.killReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.killReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) Delete(log lager.Logger, force bool, bundlePath string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		log        lager.Logger
		force      bool
		bundlePath string
	}{log, force, bundlePath})
	fake.recordInvocation("Delete", []interface{}{log, force, bundlePath})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(log, force, bundlePath)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeOCIRuntime) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeOCIRuntime) DeleteArgsForCall(i int) (lager.Logger, bool, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].log, fake.deleteArgsForCall[i].force, fake.deleteArgsForCall[i].bundlePath
}

func (fake *FakeOCIRuntime) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) State(log lager.Logger, id string) (runrunc.State, error) {
	fake.stateMutex.Lock()
	ret, specificReturn := fake.stateReturnsOnCall[len(fake.stateArgsForCall)]
	fake.stateArgsForCall = append(fake.stateArgsForCall, struct {
		log lager.Logger
		id  string
	}{log, id})
	fake.recordInvocation("State", []interface{}{log, id})
	fake.stateMutex.Unlock()
	if fake.StateStub != nil {
		return fake.StateStub(log, id)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.stateReturns.result1, fake.stateReturns.result2
}

func (fake *FakeOCIRuntime) StateCallCount() int {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return len(fake.stateArgsForCall)
}

func (fake *FakeOCIRuntime) StateArgsForCall(i int) (lager.Logger, string) {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return fake.stateArgsForCall[i].log, fake.stateArgsForCall[i].id
}

func (fake *FakeOCIRuntime) StateReturns(result1 runrunc.State, result2 error) {
	fake.StateStub = nil
	fake.stateReturns = struct {
		result1 runrunc.State
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) StateReturnsOnCall(i int, result1 runrunc.State, result2 error) {
	fake.StateStub = nil
	if fake.stateReturnsOnCall == nil {
		fake.stateReturnsOnCall = make(map[int]struct {
			result1 runrunc.State
			result2 error
		})
	}
	fake.stateReturnsOnCall[i] = struct {
		result1 runrunc.State
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) Stats(log lager.Logger, id string) (gardener.ActualContainerMetrics, error) {
	fake.statsMutex.Lock()
	ret, specificReturn := fake.statsReturnsOnCall[len(fake.statsArgsForCall)]
	fake.statsArgsForCall = append(fake.statsArgsForCall, struct {
		log lager.Logger
		id  string
	}{log, id})
	fake.recordInvocation("Stats", []interface{}{log, id})
	fake.statsMutex.Unlock()
	if fake.StatsStub != nil {
		return fake.StatsStub(log, id)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.statsReturns.result1, fake.statsReturns.result2
}

func (fake *FakeOCIRuntime) StatsCallCount() int {
	fake.statsMutex.RLock()
	defer fake.statsMutex.RUnlock()
	return len(fake.statsArgsForCall)
}

func (fake *FakeOCIRuntime) StatsArgsForCall(i int) (lager.Logger, string) {
	fake.statsMutex.RLock()
	defer fake.statsMutex.RUnlock()
	return fake.statsArgsForCall[i].log, fake.statsArgsForCall[i].id
}

func (fake *FakeOCIRuntime) StatsReturns(result1 gardener.ActualContainerMetrics, result2 error) {
	fake.StatsStub = nil
	fake.statsReturns = struct {
		result1 gardener.ActualContainerMetrics
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) StatsReturnsOnCall(i int, result1 gardener.ActualContainerMetrics, result2 error) {
	fake.StatsStub = nil
	if fake.statsReturnsOnCall == nil {
		fake.statsReturnsOnCall = make(map[int]struct {
			result1 gardener.ActualContainerMetrics
			result2 error
		})
	}
	fake.statsReturnsOnCall[i] = struct {
		result1 gardener.ActualContainerMetrics
		result2 error
	}{result1, result2}
}

func (fake *FakeOCIRuntime) WatchEvents(log lager.Logger, id string, eventsNotifier runrunc.EventsNotifier) error {
	fake.watchEventsMutex.Lock()
	ret, specificReturn := fake.watchEventsReturnsOnCall[len(fake.watchEventsArgsForCall)]
	fake.watchEventsArgsForCall = append(fake.watchEventsArgsForCall, struct {
		log            lager.Logger
		id             string
		eventsNotifier runrunc.EventsNotifier
	}{log, id, eventsNotifier})
	fake.recordInvocation("WatchEvents", []interface{}{log, id, eventsNotifier})
	fake.watchEventsMutex.Unlock()
	if fake.WatchEventsStub != nil {
		return fake.WatchEventsStub(log, id, eventsNotifier)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.watchEventsReturns.result1
}

func (fake *FakeOCIRuntime) WatchEventsCallCount() int {
	fake.watchEventsMutex.RLock()
	defer fake.watchEventsMutex.RUnlock()
	return len(fake.watchEventsArgsForCall)
}

func (fake *FakeOCIRuntime) WatchEventsArgsForCall(i int) (lager.Logger, string, runrunc.EventsNotifier) {
	fake.watchEventsMutex.RLock()
	defer fake.watchEventsMutex.RUnlock()
	return fake.watchEventsArgsForCall[i].log, fake.watchEventsArgsForCall[i].id, fake.watchEventsArgsForCall[i].eventsNotifier
}

func (fake *FakeOCIRuntime) WatchEventsReturns(result1 error) {
	fake.WatchEventsStub = nil
	fake.watchEventsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) WatchEventsReturnsOnCall(i int, result1 error) {
	fake.WatchEventsStub = nil
	if fake.watchEventsReturnsOnCall == nil {
		fake.watchEventsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.watchEventsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOCIRuntime) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	fake.statsMutex.RLock()
	defer fake.statsMutex.RUnlock()
	fake.watchEventsMutex.RLock()
	defer fake.watchEventsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOCIRuntime) recordInvocation(key string, args []interface{}) {
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

var _ rundmc.OCIRuntime = new(FakeOCIRuntime)
