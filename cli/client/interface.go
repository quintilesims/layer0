package client

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type Client interface {
	CreateDeploy(name string, content []byte) (*models.Deploy, error)
	DeleteDeploy(id string) error
	GetDeploy(id string) (*models.Deploy, error)
	ListDeploys() ([]*models.DeploySummary, error)

	CreateEnvironment(name, instanceSize string, minCount int, userData []byte, os, amiID string) (*models.Environment, error)
	DeleteEnvironment(id string) (string, error)
	GetEnvironment(id string) (*models.Environment, error)
	ListEnvironments() ([]*models.EnvironmentSummary, error)
	UpdateEnvironment(id string, minCount int) (*models.Environment, error)
	CreateLink(sourceID string, destinationID string) error
	DeleteLink(sourceID string, destinationID string) error

	Delete(id string) error
	GetJob(id string) (*models.Job, error)
	ListJobs() ([]*models.Job, error)
	WaitForJob(jobID string, timeout time.Duration) error

	CreateLoadBalancer(name, environmentID string, healthCheck models.HealthCheck, ports []models.Port, isPublic bool, idleTimeout int) (*models.LoadBalancer, error)
	DeleteLoadBalancer(id string) (string, error)
	GetLoadBalancer(id string) (*models.LoadBalancer, error)
	ListLoadBalancers() ([]*models.LoadBalancerSummary, error)
	UpdateLoadBalancerHealthCheck(id string, healthCheck models.HealthCheck) (*models.LoadBalancer, error)
	UpdateLoadBalancerPorts(id string, ports []models.Port) (*models.LoadBalancer, error)

	CreateService(name, environmentID, deployID, loadBalancerID string) (*models.Service, error)
	DeleteService(id string) (string, error)
	UpdateService(serviceID, deployID string) (*models.Service, error)
	GetService(id string) (*models.Service, error)
	GetServiceLogs(id, start, end string, tail int) ([]*models.LogFile, error)
	ListServices() ([]*models.ServiceSummary, error)
	ScaleService(id string, scale int) (*models.Service, error)
	WaitForDeployment(serviceID string, timeout time.Duration) (*models.Service, error)

	CreateTask(name, environmentID, deployID string, overrides []models.ContainerOverride) (string, error)
	DeleteTask(id string) error
	GetTask(id string) (*models.Task, error)
	GetTaskLogs(id, start, end string, tail int) ([]*models.LogFile, error)
	ListTasks() ([]*models.TaskSummary, error)

	SelectByQuery(params map[string]string) ([]*models.EntityWithTags, error)
	GetVersion() (string, error)
	GetConfig() (*models.APIConfig, error)
	UpdateSQL() error
	RunScaler(environmentID string) (*models.ScalerRunInfo, error)
}
