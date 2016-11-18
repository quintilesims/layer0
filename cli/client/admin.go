package client

import (
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

func (c *APIClient) GetVersion() (string, error) {
	var version string
	if err := c.Execute(c.Sling("admin/").Get("version"), &version); err != nil {
		return "", err
	}

	return version, nil
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
