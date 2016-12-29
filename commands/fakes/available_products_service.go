// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/rsamban/om/api"
)

type AvailableProductsService struct {
	ListStub        func() (api.AvailableProductsOutput, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 api.AvailableProductsOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *AvailableProductsService) List() (api.AvailableProductsOutput, error) {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	fake.recordInvocation("List", []interface{}{})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub()
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *AvailableProductsService) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *AvailableProductsService) ListReturns(result1 api.AvailableProductsOutput, result2 error) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 api.AvailableProductsOutput
		result2 error
	}{result1, result2}
}

func (fake *AvailableProductsService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return fake.invocations
}

func (fake *AvailableProductsService) recordInvocation(key string, args []interface{}) {
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
