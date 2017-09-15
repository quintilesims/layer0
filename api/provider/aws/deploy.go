package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
)

type DeployProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Config   config.APIConfig
}

func NewDeployProvider(a *awsc.Client, t tag.Store, c config.APIConfig) *DeployProvider {
	return &DeployProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
