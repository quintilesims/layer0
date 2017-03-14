package resource

import (
	"errors"
	"github.com/zpatrick/go-bytesize"
)

type ResourceProvider struct {
	ID              string
	InUse           bool
	UsedPorts       []int
	AvailableMemory bytesize.Bytesize
}

func NewResourceProvider(id string, inUse bool, availableMemory bytesize.Bytesize, usedPorts []int) *ResourceProvider {
	if usedPorts == nil {
		usedPorts = []int{}
	}

	return &ResourceProvider{
		ID:              id,
		InUse:           inUse,
		UsedPorts:       usedPorts,
		AvailableMemory: availableMemory,
	}
}

func (r *ResourceProvider) HasResourcesFor(consumer ResourceConsumer) bool {
	for _, wanted := range consumer.Ports {
		for _, used := range r.UsedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.Memory <= r.AvailableMemory
}

func (r *ResourceProvider) SubtractResourcesFor(consumer ResourceConsumer) error {
	if !r.HasResourcesFor(consumer) {
		return errors.New("Provider does not have adequate resources to subtract")
	}

	r.UsedPorts = append(r.UsedPorts, consumer.Ports...)
	r.AvailableMemory -= consumer.Memory
	r.InUse = true

	return nil
}

func (r *ResourceProvider) IsInUse() bool {
	return r.InUse
}

type ByMemory []*ResourceProvider

func (m ByMemory) Len() int {
	return len(m)
}

func (m ByMemory) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m ByMemory) Less(i, j int) bool {
	return m[i].AvailableMemory < m[j].AvailableMemory
}

type ByUsage []*ResourceProvider

func (m ByUsage) Len() int {
	return len(m)
}

func (m ByUsage) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m ByUsage) Less(i, j int) bool {
	return m[i].IsInUse() && !m[j].IsInUse()
}
