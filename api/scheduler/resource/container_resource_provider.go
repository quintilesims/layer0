package resource

import (
	"errors"
	"github.com/zpatrick/go-bytesize"
)

type ContainerResourceProvider struct {
	usedPorts       []int
	availableMemory bytesize.Bytesize
}

func NewContainerResourceProvider(availableMemory bytesize.Bytesize, usedPorts []int) *ContainerResourceProvider {
	return &ContainerResourceProvider{
		usedPorts:       usedPorts,
		availableMemory: availableMemory,
	}
}

func (c *ContainerResourceProvider) HasResourcesFor(resource ContainerResource) bool {
	for _, wanted := range resource.Ports {
		for _, used := range c.usedPorts {
			if wanted == used {
				return false
			}
		}
	}

	return resource.Memory <= c.availableMemory
}

func (c *ContainerResourceProvider) SubtractResourcesFor(resource ContainerResource) error {
	if !c.HasResourcesFor(resource) {
		return errors.New("Provider does not have adequate resources to subtract")
	}

	c.usedPorts = append(c.usedPorts, resource.Ports...)
	c.availableMemory -= resource.Memory

	return nil
}

func (c *ContainerResourceProvider) UsedPorts() []int {
	return c.usedPorts
}

func (c *ContainerResourceProvider) AvailableMemory() bytesize.Bytesize {
	return c.availableMemory
}
