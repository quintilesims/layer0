package client

import (
	"fmt"
	"net/url"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/rclient"
)

func (c *APIClient) CreateEnvironment(req models.CreateEnvironmentRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	var resp models.CreateEntityResponse
	if err := c.client.Post("/environment", req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) DeleteEnvironment(environmentID string) error {
	path := fmt.Sprintf("/environment/%s", environmentID)
	if err := c.client.Delete(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ListEnvironments() ([]models.EnvironmentSummary, error) {
	var environments []models.EnvironmentSummary
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

func (c *APIClient) ReadEnvironmentLogs(environmentID string, query url.Values) ([]models.LogFile, error) {
	var logs []models.LogFile
	path := fmt.Sprintf("/environment/%s/logs", environmentID)
	if err := c.client.Get(path, &logs, rclient.Query(query)); err != nil {
		return nil, err
	}

	return logs, nil
}

func (c *APIClient) UpdateEnvironment(environmentID string, req models.UpdateEnvironmentRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/environment/%s", environmentID)
	if err := c.client.Patch(path, req, nil); err != nil {
		return err
	}

	return nil
}
