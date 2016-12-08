package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateEnvironment(name, instanceSize string, minCount int, userData []byte) (*models.Environment, error) {
	req := models.CreateEnvironmentRequest{
		EnvironmentName:  name,
		InstanceSize:     instanceSize,
		MinClusterCount:  minCount,
		UserDataTemplate: userData,
	}

	var environment *models.Environment
	if err := c.Execute(c.Sling("environment/").Post("").BodyJSON(req), &environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (c *APIClient) DeleteEnvironment(id string) (string, error) {
	jobID, err := c.ExecuteWithJob(c.Sling("environment/").Delete(id))
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (c *APIClient) GetEnvironment(id string) (*models.Environment, error) {
	var environment *models.Environment
	if err := c.Execute(c.Sling("environment/").Get(id), &environment); err != nil {
		return nil, err
	}

	return environment, nil
}

func (c *APIClient) ListEnvironments() ([]*models.Environment, error) {
	var environments []*models.Environment
	if err := c.Execute(c.Sling("environment/").Get(""), &environments); err != nil {
		return nil, err
	}

	return environments, nil
}

func (c *APIClient) UpdateEnvironment(id string, minCount int) (*models.Environment, error) {
	req := models.UpdateEnvironmentRequest{
		MinClusterCount: minCount,
	}

	var environment *models.Environment
	if err := c.Execute(c.Sling("environment/").Put(id).BodyJSON(req), &environment); err != nil {
		return nil, err
	}

	return environment, nil
}
