package client

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateEnvironment(req models.CreateEnvironmentRequest) (string, error) {
	var resp models.CreateJobResponse
	if err := c.client.Post("/environment", req, &resp); err != nil {
		return "", err
	}

	return resp.JobID, nil
}

func (c *APIClient) DeleteEnvironment(environmentID string) (string, error) {
	var resp models.CreateJobResponse
	path := fmt.Sprintf("/environment/%s", environmentID)
	if err := c.client.Delete(path, nil, &resp); err != nil {
		return "", err
	}

	return resp.JobID, nil
}

func (c *APIClient) ListEnvironments() ([]*models.EnvironmentSummary, error) {
	var environments []*models.EnvironmentSummary
	if err := c.client.Get("/environment", &environments); err != nil {
		return nil, err
	}

	return environments, nil
}

func (c *APIClient) ReadEnvironment(environmentID string) (*models.Environment, error) {
	var environment *models.Environment
	path := fmt.Sprintf("/environment/%s", environmentID)
	if err := c.client.Get(path, &environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (c *APIClient) UpdateEnvironment(req models.UpdateEnvironmentRequest) (string, error) {
	var resp models.CreateJobResponse
	if err := c.client.Put("/environment", req, &resp); err != nil {
		return "", err
	}

	return resp.JobID, nil
}

func (c *APIClient) CreateLink(req models.CreateEnvironmentLinkRequest) (string, error) {
	var resp models.CreateJobResponse
	path := fmt.Sprintf("/environment/%s/link", req.SourceEnvironmentID)
	if err := c.client.Post(path, req, &resp); err != nil {
		return "", err
	}

	return resp.JobID, nil
}

func (c *APIClient) DeleteLink(req models.DeleteEnvironmentLinkRequest) (string, error) {
	var resp models.CreateJobResponse
	path := fmt.Sprintf("/environment/%s/link/%s", req.SourceEnvironmentID, req.DestEnvironmentID)
	if err := c.client.Delete(path, nil, &resp); err != nil {
		return "", err
	}

	return resp.JobID, nil
}
