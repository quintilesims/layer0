package aws

import (
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
)

type LoadBalancerProvider struct {
	AWS      *awsc.Client
	TagStore tag_store.TagStore
}

func NewLoadBalancerProvider(a *awsc.Client, t tag_store.TagStore) *LoadBalancerProvider {
	return &LoadBalancerProvider{
		AWS:      a,
		TagStore: t,
	}
}

func (d *LoadBalancerProvider) Create(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error) {
	return nil, nil
}

func (d *LoadBalancerProvider) Read(LoadBalancerID string) (*models.LoadBalancer, error) {
	return nil, nil
}

func (d *LoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	return nil, nil
}

func (d *LoadBalancerProvider) Delete(LoadBalancerID string) error {
	return nil
}
