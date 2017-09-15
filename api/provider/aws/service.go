package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

type ServiceProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
	Config   config.APIConfig
}

func NewServiceProvider(a *awsc.Client, t tag_store.TagStore, c config.APIConfig) *ServiceProvider {
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
