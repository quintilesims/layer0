package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/urfave/cli"
)

type DeployProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewDeployProvider(a *awsc.Client, t tag.Store, c *cli.Context) *DeployProvider {
	return &DeployProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}
