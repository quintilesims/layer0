package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
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

func (d *DeployProvider) Create(req models.CreateDeployRequest) (*models.Deploy, error) {
	return nil, nil
}

func (d *DeployProvider) Read(deployID string) (*models.Deploy, error) {
	return nil, nil
}

func (d *DeployProvider) List() ([]models.DeploySummary, error) {
	return nil, nil
}

func (d *DeployProvider) Delete(deployID string) error {
	return nil
}
