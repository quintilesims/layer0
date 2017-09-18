package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/models"
)

type ServiceProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
}

func NewServiceProvider(a *awsc.Client, t tag.Store) *ServiceProvider {
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
	}
}

func (s *ServiceProvider) Create(req models.CreateServiceRequest) (*models.Service, error) {
	return nil, nil
}

func (s *ServiceProvider) Read(ServiceID string) (*models.Service, error) {
	return nil, nil
}

func (s *ServiceProvider) List() ([]models.ServiceSummary, error) {
	return nil, nil
}

func (s *ServiceProvider) Delete(ServiceID string) error {
	return nil
}

func (s *ServiceProvider) Update(req models.UpdateServiceRequest) error {
	return nil
}
