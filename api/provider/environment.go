package provider

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentProvider interface {
	Create(req models.CreateEnvironmentRequest) (string, error)
	Delete(environmentID string) error
	List() ([]models.EnvironmentSummary, error)
	Logs(environmentID string, tail int, start, end time.Time) ([]models.LogFile, error)
	Read(environmentID string) (*models.Environment, error)
	Update(environmentID string, req models.UpdateEnvironmentRequest) error
}
