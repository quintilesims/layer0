package models

import (
	"fmt"

	bytesize "github.com/zpatrick/go-bytesize"
)

type ResourceProvider struct {
	AgentIsConnected bool              `json:"agent_connected"`
	AvailableCPU     bytesize.Bytesize `json:"available_cpu"`
	AvailableMemory  bytesize.Bytesize `json:"available_memory"`
	ID               string            `json:"id"`
	InUse            bool              `json:"in_use"`
	Status           string            `json:"status"`
	UsedPorts        []int             `json:"used_ports"`
}

func (r *ResourceProvider) HasResourcesFor(consumer ResourceConsumer) bool {
	if !r.AgentIsConnected {
		return false
	}

	for _, wanted := range consumer.Ports {
		for _, used := range r.UsedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return consumer.CPU <= r.AvailableCPU && consumer.Memory <= r.AvailableMemory
}

func (r *ResourceProvider) SubtractResourcesFor(consumer ResourceConsumer) error {
	if !r.HasResourcesFor(consumer) {
		return fmt.Errorf("Cannot subtract resources for consumer '%s' from provider '%s'.", consumer.ID, r.ID)
	}

	r.AvailableCPU -= consumer.CPU
	r.AvailableMemory -= consumer.Memory
	r.InUse = true
	r.UsedPorts = append(r.UsedPorts, consumer.Ports...)

	return nil
}
