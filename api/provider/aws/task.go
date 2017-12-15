package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/urfave/cli"
)

type TaskProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewTaskProvider(a *awsc.Client, t tag.Store, c *cli.Context) *TaskProvider {
	return &TaskProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}
