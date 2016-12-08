package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateLoadBalancer(name, environmentID string, ports []models.Port, isPublic bool) (*models.LoadBalancer, error) {
	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: name,
		EnvironmentID:    environmentID,
		Ports:            ports,
		IsPublic:         isPublic,
	}

	var loadBalancer *models.LoadBalancer
	if err := c.Execute(c.Sling("loadbalancer/").Post("").BodyJSON(req), &loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (c *APIClient) DeleteLoadBalancer(id string) (string, error) {
	jobID, err := c.ExecuteWithJob(c.Sling("loadbalancer/").Delete(id))
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (c *APIClient) GetLoadBalancer(id string) (*models.LoadBalancer, error) {
	var loadBalancer *models.LoadBalancer
	if err := c.Execute(c.Sling("loadbalancer/").Get(id), &loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (c *APIClient) ListLoadBalancers() ([]*models.LoadBalancer, error) {
	var loadBalancers []*models.LoadBalancer
	if err := c.Execute(c.Sling("loadbalancer/").Get(""), &loadBalancers); err != nil {
		return nil, err
	}

	return loadBalancers, nil
}

func (c *APIClient) UpdateLoadBalancer(id string, ports []models.Port) (*models.LoadBalancer, error) {
	req := models.UpdateLoadBalancerRequest{
		Ports: ports,
	}

	var loadBalancer *models.LoadBalancer
	if err := c.Execute(c.Sling("loadbalancer/").Put(id).BodyJSON(req), &loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}
