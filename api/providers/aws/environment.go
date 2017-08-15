package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

type EnvironmentProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
}

func NewEnvironmentProvider(a *awsc.Client, t tag_store.TagStore) *EnvironmentProvider {
	return &EnvironmentProvider{
		AWS:      a,
		TagStore: t,
	}
}
