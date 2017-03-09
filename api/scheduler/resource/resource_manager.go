package resource

type PendingResourceGetter interface {
	GetPendingResources() ([]ContainerResource, error)
}

type ResourceProviderGetter interface {
	GetResourceProviders() ([]ContainerResource, error)
}

type ResourceManager struct {
	pendingResourceGetter  PendingResourceGetter
	resourceProviderGetter ResourceProviderGetter
}

func NewResourceManager(prg PendingResourceGetter, rpg ResourceProviderGetter) *ResourceManager {
	return &ResourceManager{
		pendingResourceGetter:  prg,
		resourceProviderGetter: rpg,
	}
}

func (r *ResourceManager) Run() error {
	pendingResources, err := r.pendingResourceGetter.GetPendingResources()
	if err != nil {
		return err
	}

	resourceProviders, err := r.resourceProviderGetter.GetResourceProviders()
	if err != nil {
		return err
	}

	print(pendingResources, resourceProviders)
	return nil
}
