package client

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/types"
	"gitlab.imshealth.com/xfra/layer0/common/waitutils"
	"time"
)

func (c *APIClient) DeleteJob(id string) error {
	var response *string
	if err := c.Execute(c.Sling("job/").Delete(id), &response); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) GetJob(id string) (*models.Job, error) {
	var job *models.Job
	if err := c.Execute(c.Sling("job/").Get(id), &job); err != nil {
		return nil, err
	}

	return job, nil
}

func (c *APIClient) ListJobs() ([]*models.Job, error) {
	var jobs []*models.Job
	if err := c.Execute(c.Sling("job/").Get(""), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *APIClient) WaitForJob(jobID string, timeout time.Duration) error {
	waiter := waitutils.Waiter{
		Name:    "WaitForJob",
		Timeout: timeout,
		Delay:   time.Second * 5,
		Clock:   c.Clock,
		Check: func() (bool, error) {
			job, err := c.GetJob(jobID)
			if err != nil {
				return false, err
			}

			if types.JobStatus(job.JobStatus) == types.Error {
				text := "An error occured during the job's execution. \n"
				text += fmt.Sprintf("Use 'l0 job logs %s' for more information", job.JobID)
				return false, fmt.Errorf(text)
			}

			if types.JobStatus(job.JobStatus) == types.Completed {
				return true, nil
			}

			return false, nil
		},
	}

	return waiter.Wait()
}
