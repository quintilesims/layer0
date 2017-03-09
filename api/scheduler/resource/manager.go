package resource

import (
	"fmt"
	"github.com/quintilesims/layer0/common/errors"
)

type ResourceConsumerGetter func() ([]ResourceConsumer, error)

type ResourceManager struct {
	getPendingResources ResourceConsumerGetter
	providerManager     ResourceProviderManager
}

func NewResourceManager(p ResourceProviderManager, g ResourceConsumerGetter) *ResourceManager {
	return &ResourceManager{
		getPendingResources: g,
		providerManager:     p,
	}
}

func (r *ResourceManager) Run() error {
	pendingResources, err := r.getPendingResources()
	if err != nil {
		return err
	}

	resourceProviders, err := r.providerManager.GetResourceProviders()
	if err != nil {
		return err
	}

	var errs []error
	for _, consumer := range pendingResources {
		hasRoom := false

		for _, provider := range resourceProviders {
			if provider.HasResourcesFor(consumer) {
				hasRoom = true
				provider.SubtractResourcesFor(consumer)
				break
			}
		}

		if !hasRoom {
			memory := r.providerManager.MemoryPerProvider()
			newProvider := NewResourceProvider(memory, nil)

			if !newProvider.HasResourcesFor(consumer) {
				err := fmt.Errorf("Resource '%s' is too large for current provider size %v!", consumer.ID, memory)
				errs = append(errs, err)
				continue
			}

			newProvider.SubtractResourcesFor(consumer)
			resourceProviders = append(resourceProviders, newProvider)
		}
	}

	newScale := len(resourceProviders)
	if err := r.providerManager.ScaleTo(newScale); err != nil {
		errs = append(errs, err)
	}

	return errors.MultiError(errs)
}
