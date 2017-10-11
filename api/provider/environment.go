package provider

import "github.com/quintilesims/layer0/common/models"

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) (string, error)
	Delete(environmentID string) error
	List() ([]models.EnvironmentSummary, error)
	Read(environmentID string) (*models.Environment, error)
	Update(req models.UpdateEnvironmentRequest) error
	Link(sourceEnvironmentID, destEnvironmentID string) error
	Unlink(sourceEnvironmentID, destEnvironmentID string) error
}
