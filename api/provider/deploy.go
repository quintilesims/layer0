package provider

import "github.com/quintilesims/layer0/common/models"

type DeployProvider interface {
	Create(req models.CreateDeployRequest) (*models.Deploy, error)
	Read(deployID string) (*models.Deploy, error)
	List() ([]models.DeploySummary, error)
	Delete(deployID string) error
	Update(deployID string) error
}
