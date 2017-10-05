package provider

import "github.com/quintilesims/layer0/common/models"

type DeployProvider interface {
	Create(req models.CreateDeployRequest) (string, error)
	Delete(deployID string) error
	List() ([]models.DeploySummary, error)
	Read(deployID string) (*models.Deploy, error)
}
