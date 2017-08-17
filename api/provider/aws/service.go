package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
)

type ServiceProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
}

func NewServiceProvider(a *awsc.Client, t tag_store.TagStore) *ServiceProvider {
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
