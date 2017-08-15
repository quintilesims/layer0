package provider

import "github.com/quintilesims/layer0/common/models"

type JobProvider interface {
	Create(req models.CreateJobRequest) error
	Read(jobID string) (*models.Job, error)
	List() ([]*models.Job, error)
	Delete(jobID string) error
}
