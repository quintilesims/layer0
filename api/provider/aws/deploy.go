package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
)

type DeployProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
}

func NewDeployProvider(a *awsc.Client, t tag.Store) *DeployProvider {
	return &DeployProvider{
		AWS:      a,
		TagStore: t,
	}
}
