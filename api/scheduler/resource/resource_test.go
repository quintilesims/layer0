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
	getPendingResources := func() ([]ResourceConsumer, error) {
		return r.PendingResources, nil
	}

	providerManager := &TestResourceProviderManager{
		GetResourceProvidersf: func() ([]*ResourceProvider, error) {
			return r.ResourceProviders, nil
		},
		MemoryPerProviderf: func() bytesize.Bytesize {
			return r.MemoryPerProvider
		},
		ScaleTof: func(scale int) error {
			assert.Equal(t, r.ExpectedScale, scale)
			return nil
		},
	}

	return NewResourceManager(providerManager, getPendingResources)
}

type TestResourceProviderManager struct {
	GetResourceProvidersf func() ([]*ResourceProvider, error)
	MemoryPerProviderf    func() bytesize.Bytesize
	ScaleTof              func(int) error
}

func (t *TestResourceProviderManager) GetResourceProviders() ([]*ResourceProvider, error) {
	return t.GetResourceProvidersf()
}

func (t *TestResourceProviderManager) MemoryPerProvider() bytesize.Bytesize {
	return t.MemoryPerProviderf()
}

func (t *TestResourceProviderManager) ScaleTo(scale int) error {
	return t.ScaleTof(scale)
}
