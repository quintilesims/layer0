package resource

import (
	"errors"
	"github.com/zpatrick/go-bytesize"
)

type ResourceProvider struct {
	usedPorts       []int
	availableMemory bytesize.Bytesize
}

func NewResourceProvider(availableMemory bytesize.Bytesize, usedPorts []int) *ResourceProvider {
	if usedPorts == nil {
		usedPorts = []int{}
	}

	return &ResourceProvider{
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

	return nil
}
