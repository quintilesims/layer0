package provider

import "github.com/quintilesims/layer0/common/models"

type LoadBalancerProvider interface {
	Create(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error)
	Read(loadBalancerID string) (*models.LoadBalancer, error)
	List() ([]models.LoadBalancerSummary, error)
	Delete(loadBalancerID string) error
	Update(req models.UpdateLoadBalancerRequest) (*models.LoadBalancer, error)
}
