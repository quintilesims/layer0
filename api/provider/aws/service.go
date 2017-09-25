package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
)

type ServiceProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Config   config.APIConfig
}

func NewServiceProvider(a *awsc.Client, t tag.Store, c config.APIConfig) *ServiceProvider {
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
