package provider

import "github.com/quintilesims/layer0/common/models"

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) (string, error)
	Delete(environmentID string) error
	List() ([]models.EnvironmentSummary, error)
	Read(environmentID string) (*models.Environment, error)
	Update(req models.UpdateEnvironmentRequest) error
	Link(req models.CreateEnvironmentLinkRequest) error
	Unlink(req models.DeleteEnvironmentLinkRequest) error
}
