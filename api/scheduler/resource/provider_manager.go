package resource

import (
	"github.com/zpatrick/go-bytesize"
)

type ResourceProviderManager interface {
	GetResourceProviders() ([]*ResourceProvider, error)
	MemoryPerProvider() bytesize.Bytesize
	ScaleTo(int) error
}
