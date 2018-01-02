package scaler

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	bytesize "github.com/zpatrick/go-bytesize"
)

type ResourceProvider struct {
	AgentIsConnected bool              `json:"agent_connected"`
	AvailableCPU     int               `json:"available_cpu"`
	AvailableMemory  bytesize.Bytesize `json:"available_memory"`
	ID               string            `json:"id"`
	InUse            bool              `json:"in_use"`
	Status           string            `json:"status"`
	UsedPorts        []int             `json:"used_ports"`
}

func NewResourceProvider(agent *bool, cpu int, memory bytesize.Bytesize, id string, inUse *bool, status *string, ports *[]int) *ResourceProvider {
	r := &ResourceProvider{
		AvailableCPU:    cpu,
		AvailableMemory: memory,
		ID:              id,
		UsedPorts:       defaultPorts(),
	}

	if agent != nil {
		r.AgentIsConnected = *agent
	}

	if inUse != nil {
		r.InUse = *inUse
	}

	if status != nil {
		r.Status = *status
	}

	if ports != nil {
		r.UsedPorts = append(r.UsedPorts, *ports...)
	}

	return r
}

func (r *ResourceProvider) String() string {
	s := "&scaler.ResourceProvider{AgentIsConnected:%t, AvailableCPU:%d, AvailableMemory:%v, ID:%s, InUse:%t, Status:%s, UsedPorts:%v}"
	return fmt.Sprintf(s, r.AgentIsConnected, r.AvailableCPU, r.AvailableMemory.Format("mib"), r.ID, r.InUse, r.Status, r.UsedPorts)
}

func (r *ResourceProvider) HasResourcesFor(consumer ResourceConsumer) bool {
	if !r.AgentIsConnected || r.Status != ecs.ContainerInstanceStatusActive {
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
