package client

import (
	"fmt"
	"net/url"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/rclient"
)

func (c *APIClient) CreateService(req models.CreateServiceRequest) (string, error) {
	var resp models.CreateEntityResponse
	if err := c.client.Post("/service", req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) DeleteService(serviceID string) (string, error) {
	var resp models.CreateEntityResponse
	path := fmt.Sprintf("/service/%s", serviceID)
	if err := c.client.Delete(path, nil, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) ListServices() ([]*models.ServiceSummary, error) {
	var services []*models.ServiceSummary
	if err := c.client.Get("/service", &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (c *APIClient) ReadService(serviceID string) (*models.Service, error) {
	var service *models.Service
	path := fmt.Sprintf("/service/%s", serviceID)
	if err := c.client.Get(path, &service); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *APIClient) ReadServiceLogs(serviceID string, query url.Values) ([]*models.LogFile, error) {
	var logs []*models.LogFile
	path := fmt.Sprintf("/service/%s/logs", serviceID)
	if err := c.client.Get(path, &logs, rclient.Query(query)); err != nil {
		return nil, err
	}

	return logs, nil
}

func (c *APIClient) UpdateService(serviceID string, req models.UpdateServiceRequest) (string, error) {
	var resp models.CreateEntityResponse
	path := fmt.Sprintf("/service/%s", serviceID)
	if err := c.client.Patch(path, req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}
