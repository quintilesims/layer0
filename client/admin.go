package client

import (
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) ReadConfig() (*models.APIConfig, error) {
	var config *models.APIConfig
	if err := c.client.Get("/admin/config", &config); err != nil {
		return nil, err
	}

	return config, nil
}
