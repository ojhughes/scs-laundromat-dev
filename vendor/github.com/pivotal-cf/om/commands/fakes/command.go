// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	jhandacommands "github.com/pivotal-cf/jhanda/commands"
)

type Command struct {
	ExecuteStub        func(args []string) error
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		args []string
	}
	executeReturns struct {
		result1 error
	}
	executeReturnsOnCall map[int]struct {
		result1 error
	}
	UsageStub        func() jhandacommands.Usage
	usageMutex       sync.RWMutex
	usageArgsForCall []struct{}
	usageReturns     struct {
		result1 jhandacommands.Usage
	}
	usageReturnsOnCall map[int]struct {
		result1 jhandacommands.Usage
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Command) Execute(args []string) error {
	var argsCopy []string
	if args != nil {
		argsCopy = make([]string, len(args))
		copy(argsCopy, args)
	}
	fake.executeMutex.Lock()
	ret, specificReturn := fake.executeReturnsOnCall[len(fake.executeArgsForCall)]
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		args []string
	}{argsCopy})
	fake.recordInvocation("Execute", []interface{}{argsCopy})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub(args)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.executeReturns.result1
}

func (fake *Command) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *Command) ExecuteArgsForCall(i int) []string {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].args
}

func (fake *Command) ExecuteReturns(result1 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 error
	}{result1}
}

func (fake *Command) ExecuteReturnsOnCall(i int, result1 error) {
	fake.ExecuteStub = nil
	if fake.executeReturnsOnCall == nil {
		fake.executeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.executeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Command) Usage() jhandacommands.Usage {
	fake.usageMutex.Lock()
	ret, specificReturn := fake.usageReturnsOnCall[len(fake.usageArgsForCall)]
	fake.usageArgsForCall = append(fake.usageArgsForCall, struct{}{})
	fake.recordInvocation("Usage", []interface{}{})
	fake.usageMutex.Unlock()
	if fake.UsageStub != nil {
		return fake.UsageStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.usageReturns.result1
}

func (fake *Command) UsageCallCount() int {
	fake.usageMutex.RLock()
	defer fake.usageMutex.RUnlock()
	return len(fake.usageArgsForCall)
}

func (fake *Command) UsageReturns(result1 jhandacommands.Usage) {
	fake.UsageStub = nil
	fake.usageReturns = struct {
		result1 jhandacommands.Usage
	}{result1}
}

func (fake *Command) UsageReturnsOnCall(i int, result1 jhandacommands.Usage) {
	fake.UsageStub = nil
	if fake.usageReturnsOnCall == nil {
		fake.usageReturnsOnCall = make(map[int]struct {
			result1 jhandacommands.Usage
		})
	}
	fake.usageReturnsOnCall[i] = struct {
		result1 jhandacommands.Usage
	}{result1}
}

func (fake *Command) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	fake.usageMutex.RLock()
	defer fake.usageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Command) recordInvocation(key string, args []interface{}) {
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

var _ jhandacommands.Command = new(Command)
