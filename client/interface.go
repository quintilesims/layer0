package client

import (
	"net/url"

	"github.com/quintilesims/layer0/common/models"
)

type Client interface {
	CreateDeploy(req models.CreateDeployRequest) (string, error)
	DeleteDeploy(deployID string) error
	ListDeploys() ([]models.DeploySummary, error)
	ReadDeploy(deployID string) (*models.Deploy, error)

	CreateEnvironment(req models.CreateEnvironmentRequest) (string, error)
	DeleteEnvironment(environmentID string) error
	ListEnvironments() ([]models.EnvironmentSummary, error)
	ReadEnvironment(environmentID string) (*models.Environment, error)
	ReadEnvironmentLogs(environmentID string, query url.Values) ([]models.LogFile, error)
	UpdateEnvironment(environmentID string, req models.UpdateEnvironmentRequest) error

	CreateLoadBalancer(req models.CreateLoadBalancerRequest) (string, error)
	DeleteLoadBalancer(loadBalancerID string) error
	ListLoadBalancers() ([]models.LoadBalancerSummary, error)
	ReadLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error)
	UpdateLoadBalancer(loadBalancerID string, req models.UpdateLoadBalancerRequest) error

	CreateService(req models.CreateServiceRequest) (string, error)
	DeleteService(serviceID string) error
	ListServices() ([]models.ServiceSummary, error)
	ReadService(serviceID string) (*models.Service, error)
	ReadServiceLogs(serviceID string, query url.Values) ([]models.LogFile, error)
	UpdateService(serviceID string, req models.UpdateServiceRequest) error

	CreateTask(req models.CreateTaskRequest) (string, error)
	DeleteTask(taskID string) error
	ListTasks() ([]models.TaskSummary, error)
	ReadTask(taskID string) (*models.Task, error)
	ReadTaskLogs(taskID string, query url.Values) ([]models.LogFile, error)

	ListTags(query url.Values) (models.Tags, error)

	ReadConfig() (*models.APIConfig, error)
	ReadAdminLogs(query url.Values) ([]models.LogFile, error)
}
