package aws

import (
<<<<<<< bace4b6d6d143ceb9cf4c0f877df724826e78d1a
<<<<<<< 0f6b259f78b9a344646eca80c79aa47cb44bc13e
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
=======
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
>>>>>>> initial change
	"github.com/quintilesims/layer0/common/models"
=======
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
>>>>>>> Break CRUD actions up into separate files
)

type ServiceProvider struct {
	AWS      *awsc.Client
<<<<<<< 0f6b259f78b9a344646eca80c79aa47cb44bc13e
	TagStore tag.Store
}

func NewServiceProvider(a *awsc.Client, t tag.Store) *ServiceProvider {
=======
	TagStore tag_store.TagStore
	Config   config.APIConfig
}

func NewServiceProvider(a *awsc.Client, t tag_store.TagStore, c config.APIConfig) *ServiceProvider {
>>>>>>> initial change
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}
