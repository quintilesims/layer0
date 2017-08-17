package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

type LoadBalancerProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
	Config   config.APIConfig
}

func NewLoadBalancerProvider(a *awsc.Client, t tag_store.TagStore, c config.APIConfig) *LoadBalancerProvider {
	return &LoadBalancerProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}

func (l *LoadBalancerProvider) Delete(loadBalancerID string) error {
	return nil
}
