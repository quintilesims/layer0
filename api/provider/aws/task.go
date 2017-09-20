package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
)

type TaskProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Config   config.APIConfig
}

func NewTaskProvider(a *awsc.Client, t tag.Store, c config.APIConfig) *TaskProvider {
	return &TaskProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
