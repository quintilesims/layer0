package client

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateLoadBalancer(req models.CreateLoadBalancerRequest) (string, error) {
	var resp models.CreateEntityResponse
	if err := c.client.Post("/loadbalancer", req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) DeleteLoadBalancer(loadBalancerID string) (error) {
	path := fmt.Sprintf("/loadbalancer/%s", loadBalancerID)
	if err := c.client.Delete(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ListLoadBalancers() ([]*models.LoadBalancerSummary, error) {
	var loadbalancers []*models.LoadBalancerSummary
	if err := c.client.Get("/loadbalancer", &loadbalancers); err != nil {
		return nil, err
	}

	return loadbalancers, nil
}

func (c *APIClient) ReadLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error) {
	var loadbalancer *models.LoadBalancer
	path := fmt.Sprintf("/loadbalancer/%s", loadBalancerID)
	if err := c.client.Get(path, &loadbalancer); err != nil {
		return nil, err
	}

	return loadbalancer, nil
}

func (c *APIClient) UpdateLoadBalancer(loadBalancerID string, req models.UpdateLoadBalancerRequest) ( error) {
	path := fmt.Sprintf("/loadbalancer/%s", loadBalancerID)
	if err := c.client.Patch(path, req, nil); err != nil {
		return err
	}

	return nil
}
