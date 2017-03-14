package resource

import (
	"github.com/zpatrick/go-bytesize"
)

type ResourceProviderManager interface {
	GetResourceProviders(environmentID string) ([]*ResourceProvider, error)
	MemoryPerProvider() bytesize.Bytesize
	ScaleTo(environmentID string, scale int, unusedProviders ...*ResourceProvider) (int, error)
}
