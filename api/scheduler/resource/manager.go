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

type RunInfo struct {
	EnvironmentID           string              `json:"environment_id"`
	ScaleBeforeRun          int                 `json:"scale_before_run"`
	DesiredScaleAfterRun    int                 `json:"desired_scale_after_run"`
	ActualScaleAfterRun     int                 `json:"actual_scale_after_run"`
	UnusedResourceProviders int                 `json:"unused_resource_providers"`
	PendingResources        []ResourceConsumer  `json:"pending_resources"`
	ResourceProviders       []*ResourceProvider `json:"resource_providers"`
}

// todo: this doesn't check reserved CPU units
func (r *ResourceManager) Run(environmentID string) (*RunInfo, error) {
	pendingResources, err := r.getPendingResources(environmentID)
	if err != nil {
		return nil, err
	}

	resourceProviders, err := r.providerManager.GetResourceProviders(environmentID)
	if err != nil {
		return nil, err
	}

	scaleBeforeRun := len(resourceProviders)
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
			newProvider := NewResourceProvider("<new resource provider>", false, memory, nil)

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

	desiredScale := len(resourceProviders) - len(unusedProviders)
	actualScale, err := r.providerManager.ScaleTo(environmentID, desiredScale, unusedProviders...)
	if err != nil {
		errs = append(errs, err)
	}

	info := &RunInfo{
		EnvironmentID:           environmentID,
		PendingResources:        pendingResources,
		ResourceProviders:       resourceProviders,
		ScaleBeforeRun:          scaleBeforeRun,
		DesiredScaleAfterRun:    desiredScale,
		ActualScaleAfterRun:     actualScale,
		UnusedResourceProviders: len(unusedProviders),
	}

	return info, errors.MultiError(errs)
}
