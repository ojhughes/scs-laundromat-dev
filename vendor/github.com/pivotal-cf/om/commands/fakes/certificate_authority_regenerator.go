// This file was generated by counterfeiter
package fakes

import (
	"sync"
)

type CertificateAuthorityRegenerator struct {
	RegenerateStub        func() error
	regenerateMutex       sync.RWMutex
	regenerateArgsForCall []struct{}
	regenerateReturns     struct {
		result1 error
	}
	regenerateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CertificateAuthorityRegenerator) Regenerate() error {
	fake.regenerateMutex.Lock()
	ret, specificReturn := fake.regenerateReturnsOnCall[len(fake.regenerateArgsForCall)]
	fake.regenerateArgsForCall = append(fake.regenerateArgsForCall, struct{}{})
	fake.recordInvocation("Regenerate", []interface{}{})
	fake.regenerateMutex.Unlock()
	if fake.RegenerateStub != nil {
		return fake.RegenerateStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.regenerateReturns.result1
}

func (fake *CertificateAuthorityRegenerator) RegenerateCallCount() int {
	fake.regenerateMutex.RLock()
	defer fake.regenerateMutex.RUnlock()
	return len(fake.regenerateArgsForCall)
}

func (fake *CertificateAuthorityRegenerator) RegenerateReturns(result1 error) {
	fake.RegenerateStub = nil
	fake.regenerateReturns = struct {
		result1 error
	}{result1}
}

func (fake *CertificateAuthorityRegenerator) RegenerateReturnsOnCall(i int, result1 error) {
	fake.RegenerateStub = nil
	if fake.regenerateReturnsOnCall == nil {
		fake.regenerateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.regenerateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CertificateAuthorityRegenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.regenerateMutex.RLock()
	defer fake.regenerateMutex.RUnlock()
	return fake.invocations
}

func (fake *CertificateAuthorityRegenerator) recordInvocation(key string, args []interface{}) {
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
