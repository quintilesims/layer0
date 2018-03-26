package client

import (
	"net/url"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/rclient"
)

func (c *APIClient) ReadConfig() (*models.APIConfig, error) {
	var config *models.APIConfig
	if err := c.client.Get("/admin/config", &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *APIClient) ReadLayer0InstanceLogs(query url.Values) ([]models.LogFile, error) {
	var logs []models.LogFile
	if err := c.client.Get("/admin/instancelogs", &logs, rclient.Query(query)); err != nil {
		return nil, err
	}

	return logs, nil
}
