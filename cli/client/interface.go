package client

import (
	"github.com/quintilesims/layer0/common/models"
	"time"
)

type Client interface {
	CreateDeploy(name string, content []byte) (*models.Deploy, error)
	DeleteDeploy(id string) error
	GetDeploy(id string) (*models.Deploy, error)
	ListDeploys() ([]*models.Deploy, error)

	CreateCertificate(name string, public, private, chain []byte) (*models.Certificate, error)
	DeleteCertificate(id string) error
	GetCertificate(id string) (*models.Certificate, error)
	ListCertificates() ([]*models.Certificate, error)

	CreateEnvironment(name, instanceSize string, minCount int, userData []byte) (*models.Environment, error)
	DeleteEnvironment(id string) (string, error)
	GetEnvironment(id string) (*models.Environment, error)
	ListEnvironments() ([]*models.Environment, error)
	UpdateEnvironment(id string, minCount int) (*models.Environment, error)

	Delete(id string) error
	SelectByID(id string) (*models.Job, error)
	SelectAll() ([]*models.Job, error)
	WaitForJob(jobID string, timeout time.Duration) error

	CreateLoadBalancer(name, environmentID string, ports []models.Port, isPublic bool) (*models.LoadBalancer, error)
	DeleteLoadBalancer(id string) (string, error)
	GetLoadBalancer(id string) (*models.LoadBalancer, error)
	ListLoadBalancers() ([]*models.LoadBalancer, error)
	UpdateLoadBalancer(id string, ports []models.Port) (*models.LoadBalancer, error)

	CreateService(name, environmentID, deployID, loadBalancerID string) (*models.Service, error)
	DeleteService(id string) (string, error)
	UpdateService(serviceID, deployID string) (*models.Service, error)
	GetService(id string) (*models.Service, error)
	GetServiceLogs(id string, tail int) ([]*models.LogFile, error)
	ListServices() ([]*models.Service, error)
	ScaleService(id string, scale int) (*models.Service, error)
	WaitForDeployment(serviceID string, timeout time.Duration) (*models.Service, error)

	CreateTask(name, environmentID, deployID string, copies int, overrides []models.ContainerOverride) (*models.Task, error)
	DeleteTask(id string) error
	GetTask(id string) (*models.Task, error)
	GetTaskLogs(id string, tail int) ([]*models.LogFile, error)
	ListTasks() ([]*models.Task, error)

	SelectByQuery(params map[string]string) ([]*models.EntityWithTags, error)
	GetVersion() (string, error)
	GetConfig() (*models.APIConfig, error)
	UpdateSQL() error
}
