package resource

import (
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/go-bytesize"
	"testing"
)

type TestResourceManager struct {
	PendingResources  []ResourceConsumer
	ResourceProviders []*ResourceProvider
	ExpectedScale     int
	MemoryPerProvider bytesize.Bytesize
}

func (r *TestResourceManager) Manager(t *testing.T) *ResourceManager {
	getPendingResources := func(string) ([]ResourceConsumer, error) {
		return r.PendingResources, nil
	}

	providerManager := &TestResourceProviderManager{
		GetResourceProvidersf: func(string) ([]*ResourceProvider, error) {
			return r.ResourceProviders, nil
		},
		MemoryPerProviderf: func() bytesize.Bytesize {
			return r.MemoryPerProvider
		},
		ScaleTof: func(environmentID string, scale int, unusedProviders ...*ResourceProvider) error {
			assert.Equal(t, r.ExpectedScale, scale)
			return nil
		},
	}

	return NewResourceManager(providerManager, getPendingResources)
}

type TestResourceProviderManager struct {
	GetResourceProvidersf func(string) ([]*ResourceProvider, error)
	MemoryPerProviderf    func() bytesize.Bytesize
	ScaleTof              func(string, int, ...*ResourceProvider) error
}

func (t *TestResourceProviderManager) GetResourceProviders(environmentID string) ([]*ResourceProvider, error) {
	return t.GetResourceProvidersf(environmentID)
}

func (t *TestResourceProviderManager) MemoryPerProvider() bytesize.Bytesize {
	return t.MemoryPerProviderf()
}

func (t *TestResourceProviderManager) ScaleTo(environmentID string, scale int, unusedProviders ...*ResourceProvider) error {
	return t.ScaleTof(environmentID, scale, unusedProviders...)
}
