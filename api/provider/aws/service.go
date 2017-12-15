package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/urfave/cli"
)

type ServiceProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewServiceProvider(a *awsc.Client, t tag.Store, c *cli.Context) *ServiceProvider {
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}
