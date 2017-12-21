package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/urfave/cli"
)

type EnvironmentProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewEnvironmentProvider(a *awsc.Client, t tag.Store, c *cli.Context) *EnvironmentProvider {
	return &EnvironmentProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}
