package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

type EnvironmentProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
	Config   config.APIConfig
}

func NewEnvironmentProvider(a *awsc.Client, t tag_store.TagStore, c config.APIConfig) *EnvironmentProvider {
	return &EnvironmentProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
