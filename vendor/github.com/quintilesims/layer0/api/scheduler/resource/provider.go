package resource

import (
	"errors"
	"sort"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/go-bytesize"
)

type ProviderManager interface {
	CalculateNewProvider(environmentID string) (*ResourceProvider, error)
	GetProviders(environmentID string) ([]*ResourceProvider, error)
	ScaleTo(environmentID string, size int, unusedProviders []*ResourceProvider) (int, error)
}

type ResourceProvider struct {
	ID              string
	inUse           bool
	usedPorts       []int
	availableMemory bytesize.Bytesize
}

func NewResourceProvider(id string, inUse bool, availableMemory bytesize.Bytesize, usedPorts []int) *ResourceProvider {
	if usedPorts == nil {
		usedPorts = []int{}
	}

	return &ResourceProvider{
		ID:              id,
		inUse:           inUse,
		usedPorts:       usedPorts,
		availableMemory: availableMemory,
	}
}

func (r *ResourceProvider) HasResourcesFor(consumer ResourceConsumer) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range r.usedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.Memory <= r.availableMemory
}

func (r *ResourceProvider) SubtractResourcesFor(consumer ResourceConsumer) error {
	if !r.HasResourcesFor(consumer) {
		return errors.New("Provider does not have adequate resources to subtract")
	}

	r.usedPorts = append(r.usedPorts, consumer.Ports...)
	r.availableMemory -= consumer.Memory
	r.inUse = true

	return nil
}

func (r *ResourceProvider) IsInUse() bool {
	return r.inUse
}

func (r ResourceProvider) ToModel() models.ResourceProvider {
	return models.ResourceProvider{
		ID:              r.ID,
		InUse:           r.inUse,
		UsedPorts:       r.usedPorts,
		AvailableMemory: r.availableMemory.Format("mib"),
	}
}

func SortProvidersByMemory(p []*ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *ResourceProvider, j *ResourceProvider) bool {
			return i.availableMemory < j.availableMemory
		},
	}

	sort.Sort(sorter)
}

func SortProvidersByUsage(p []*ResourceProvider) {
	sorter := &ResourceProviderSorter{
		Providers: p,
		lessThan: func(i *ResourceProvider, j *ResourceProvider) bool {
			return i.inUse && !j.inUse
		},
	}

	sort.Sort(sorter)
}

type ResourceProviderSorter struct {
	Providers []*ResourceProvider
	lessThan  func(*ResourceProvider, *ResourceProvider) bool
}

func (r *ResourceProviderSorter) Len() int {
	return len(r.Providers)
}

func (r *ResourceProviderSorter) Swap(i, j int) {
	r.Providers[i], r.Providers[j] = r.Providers[j], r.Providers[i]
}

func (r *ResourceProviderSorter) Less(i, j int) bool {
	return r.lessThan(r.Providers[i], r.Providers[j])
}
