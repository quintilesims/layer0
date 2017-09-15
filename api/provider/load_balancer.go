package provider

import "github.com/quintilesims/layer0/common/models"

type LoadBalancerProvider interface {
	Create(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error)
	Delete(loadBalancerID string) error
	List() ([]models.LoadBalancerSummary, error)
	Read(loadBalancerID string) (*models.LoadBalancer, error)
	Update(req models.UpdateLoadBalancerRequest) error
}
