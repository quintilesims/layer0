package provider

import "github.com/quintilesims/layer0/common/models"

type LoadBalancerProvider interface {
	Create(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error)
	Read(LoadBalancerID string) (*models.LoadBalancer, error)
	List() ([]models.LoadBalancerSummary, error)
	Delete(LoadBalancerID string) error
}
