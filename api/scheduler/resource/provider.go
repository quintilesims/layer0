package resource

import (
	"errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/go-bytesize"
)

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

type ByMemory []*ResourceProvider

func (m ByMemory) Len() int {
	return len(m)
}

func (m ByMemory) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m ByMemory) Less(i, j int) bool {
	return m[i].availableMemory < m[j].availableMemory
}

type ByUsage []*ResourceProvider

func (m ByUsage) Len() int {
	return len(m)
}

func (m ByUsage) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m ByUsage) Less(i, j int) bool {
	return m[i].inUse && !m[j].inUse
}
