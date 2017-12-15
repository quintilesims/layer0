package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/urfave/cli"
)

type LoadBalancerProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewLoadBalancerProvider(a *awsc.Client, t tag.Store, c *cli.Context) *LoadBalancerProvider {
	return &LoadBalancerProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}
