package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) GetVersion() (string, error) {
	var version string
	if err := c.Execute(c.Sling("admin/").Get("version"), &version); err != nil {
		return "", err
	}

	return version, nil
}

func (c *APIClient) GetConfig() (*models.APIConfig, error) {
	var config *models.APIConfig
	if err := c.Execute(c.Sling("admin/").Get("config"), &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *APIClient) UpdateSQL() error {
	req := models.SQLVersion{
		Version: "latest",
	}

	if err := c.Execute(c.Sling("admin/").Post("sql").BodyJSON(req), nil); err != nil {
		return err
	}

	return nil
}
