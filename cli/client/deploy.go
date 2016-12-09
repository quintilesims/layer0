package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateDeploy(name string, content []byte) (*models.Deploy, error) {
	req := models.CreateDeployRequest{
		DeployName: name,
		Dockerrun:  content,
	}

	var deploy *models.Deploy
	if err := c.Execute(c.Sling("deploy").Post("").BodyJSON(req), &deploy); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (c *APIClient) DeleteDeploy(id string) error {
	var response *string
	if err := c.Execute(c.Sling("deploy/").Delete(id), &response); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) GetDeploy(id string) (*models.Deploy, error) {
	var deploy *models.Deploy
	if err := c.Execute(c.Sling("deploy/").Get(id), &deploy); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (c *APIClient) ListDeploys() ([]*models.Deploy, error) {
	var deploys []*models.Deploy
	if err := c.Execute(c.Sling("deploy/").Get(""), &deploys); err != nil {
		return nil, err
	}

	return deploys, nil
}
