package provider

import "github.com/quintilesims/layer0/common/models"

type ServiceProvider interface {
	Create(req models.CreateServiceRequest) (*models.Service, error)
	Read(serviceID string) (*models.Service, error)
	List() ([]models.ServiceSummary, error)
	Delete(serviceID string) error
	Update(req models.UpdateServiceRequest) error
}
