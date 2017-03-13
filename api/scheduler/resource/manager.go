package resource

import (
	"fmt"
	"github.com/quintilesims/layer0/common/errors"
	"sort"
)

type ResourceConsumerGetter func(environmentID string) ([]ResourceConsumer, error)

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

// todo: this doesn't check reserved CPU units
func (r *ResourceManager) Run(environmentID string) error {
	pendingResources, err := r.getPendingResources(environmentID)
	if err != nil {
		return err
	}

	resourceProviders, err := r.providerManager.GetResourceProviders(environmentID)
	if err != nil {
		return err
	}

	var errs []error

	// check if we need to scale up
	for _, consumer := range pendingResources {
		hasRoom := false

		// first, sort by memory so we pack tasks by memory as tightly as possible
		sort.Sort(ByMemory(resourceProviders))

		// next, place any unused providers in the back of the list
		// that way, we can can delete them if we avoid placing any tasks in them
		sort.Sort(ByUsage(resourceProviders))

		for _, provider := range resourceProviders {
			if provider.HasResourcesFor(consumer) {
				hasRoom = true
				provider.SubtractResourcesFor(consumer)
				break
			}
		}

		if !hasRoom {
			memory := r.providerManager.MemoryPerProvider()
			newProvider := NewResourceProvider("", false, memory, nil)

			if !newProvider.HasResourcesFor(consumer) {
				err := fmt.Errorf("Resource '%s' is too large for current provider size %v!", consumer.ID, memory)
				errs = append(errs, err)
				continue
			}

			newProvider.SubtractResourcesFor(consumer)
			resourceProviders = append(resourceProviders, newProvider)
		}
	}

	// check if we need to scale down
	unusedProviders := []*ResourceProvider{}
	for i := 0; i < len(resourceProviders); i++ {
		if !resourceProviders[i].IsInUse() {
			unusedProviders = append(unusedProviders, resourceProviders[i])
		}
	}

	newScale := len(resourceProviders) - len(unusedProviders)
	if err := r.providerManager.ScaleTo(environmentID, newScale, unusedProviders...); err != nil {
		errs = append(errs, err)
	}

	return errors.MultiError(errs)
}
