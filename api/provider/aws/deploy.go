package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

type DeployProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
}

func NewDeployProvider(a *awsc.Client, t tag_store.TagStore) *DeployProvider {
	return &DeployProvider{
		AWS:      a,
		TagStore: t,
	}
}
