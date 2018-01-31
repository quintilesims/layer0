package client

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateDeploy(req models.CreateDeployRequest) (string, error) {
	var resp models.CreateEntityResponse
	if err := c.client.Post("/deploy", req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) DeleteDeploy(deployID string) error {
	path := fmt.Sprintf("/deploy/%s", deployID)
	if err := c.client.Delete(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ListDeploys() ([]models.DeploySummary, error) {
	var deploys []models.DeploySummary
	if err := c.client.Get("/deploy", &deploys); err != nil {
		return nil, err
	}

	return deploys, nil
}

func (c *APIClient) ReadDeploy(deployID string) (*models.Deploy, error) {
	var deploy *models.Deploy
	path := fmt.Sprintf("/deploy/%s", deployID)
	if err := c.client.Get(path, &deploy); err != nil {
		return nil, err
	}

	return deploy, nil
}
