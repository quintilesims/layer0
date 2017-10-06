package client

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) DeleteJob(jobID string) error {
	path := fmt.Sprintf("/job/%s", jobID)
	if err := c.client.Delete(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ListJobs() ([]*models.Job, error) {
	var jobs []*models.Job
	if err := c.client.Get("/job", &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *APIClient) ReadJob(jobID string) (*models.Job, error) {
	var job *models.Job
	path := fmt.Sprintf("/job/%s", jobID)
	if err := c.client.Get(path, &job); err != nil {
		return nil, err
	}

	return job, nil
}
