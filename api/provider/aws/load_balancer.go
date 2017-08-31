package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
)

type LoadBalancerProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Config   config.APIConfig
}

func NewLoadBalancerProvider(a *awsc.Client, t tag.Store, c config.APIConfig) *LoadBalancerProvider {
	return &LoadBalancerProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
