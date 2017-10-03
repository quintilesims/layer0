package provider

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type ServiceProvider interface {
	Create(req models.CreateServiceRequest) (*models.Service, error)
	Delete(serviceID string) error
	List() ([]models.ServiceSummary, error)
	Logs(taskID string, tail int, start, end time.Time) ([]models.LogFile, error)
	Read(serviceID string) (*models.Service, error)
	Update(req models.UpdateServiceRequest) error
}
