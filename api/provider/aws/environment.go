package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
)

type EnvironmentProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Config   config.APIConfig
}

func NewEnvironmentProvider(a *awsc.Client, t tag.Store, c config.APIConfig) *EnvironmentProvider {
	return &EnvironmentProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
