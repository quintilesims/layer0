package client

import (
	"github.com/quintilesims/layer0/common/models"
)

type Client interface {
	CreateEnvironment(req models.CreateEnvironmentRequest) (string, error)
	DeleteEnvironment(environmentID string) (string, error)
	ListEnvironments() ([]*models.EnvironmentSummary, error)
	ReadEnvironment(environmentID string) (*models.Environment, error)
	UpdateEnvironment(req models.UpdateEnvironmentRequest) (string, error)

	DeleteJob(jobID string) error
	ReadJob(jobID string) (*models.Job, error)
	ListJobs() ([]*models.Job, error)

	CreateService(req models.CreateServiceRequest) (string, error)
	DeleteService(serviceID string) (string, error)
	ListServices() ([]*models.ServiceSummary, error)
	ReadService(serviceID string) (*models.Service, error)
	UpdateService(req models.UpdateServiceRequest) (string, error)
}
