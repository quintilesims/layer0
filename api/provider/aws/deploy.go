package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/models"
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

func (d *DeployProvider) Update(deployID string) error {
	return nil
}
