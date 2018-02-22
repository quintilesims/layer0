package backend

import (
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/models"
)

type Backend interface {
	CreateEnvironment(environmentName, instanceSize, operatingSystem, amiID string, minClusterCount int, userData []byte) (*models.Environment, error)
	UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error)
	DeleteEnvironment(environmentID string) error
	GetEnvironment(environmentID string) (*models.Environment, error)
	ListEnvironments() ([]id.ECSEnvironmentID, error)
	CreateEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error
	DeleteEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error

	ListDeploys() ([]*models.Deploy, error)
	GetDeploy(deployID string) (*models.Deploy, error)
	CreateDeploy(name string, body []byte) (*models.Deploy, error)
	DeleteDeploy(deployID string) error

	ListServices() ([]id.ECSServiceID, error)
	GetService(environmentID, serviceID string) (*models.Service, error)
	GetEnvironmentServices(environmentID string) ([]*models.Service, error)
	CreateService(serviceName, environmentID, deployID, loadBalancerID string) (*models.Service, error)
	DeleteService(environmentID, serviceID string) error
	ScaleService(environmentID, serviceID string, count int) (*models.Service, error)
	UpdateService(environmentID, serviceID, deployID string) (*models.Service, error)
	GetServiceLogs(environmentID, serviceID, start, end string, tail int) ([]*models.LogFile, error)

	CreateTask(environmentID, deployID string, overrides []models.ContainerOverride) (string, error)
	ListTasks() ([]string, error)
	GetTask(environmentID, taskARN string) (*models.Task, error)
	GetEnvironmentTasks(environmentID string) (map[string]*models.Task, error)
	DeleteTask(environmentID, taskARN string) error
	GetTaskLogs(environmentID, taskARN, start, end string, tail int) ([]*models.LogFile, error)

	ListLoadBalancers() ([]*models.LoadBalancer, error)
	GetLoadBalancer(id string) (*models.LoadBalancer, error)
	DeleteLoadBalancer(id string) error
	CreateLoadBalancer(loadBalancerName, environmentID string, isPublic bool, ports []models.Port, healthCheck models.HealthCheck, idleTimeout int) (*models.LoadBalancer, error)
	UpdateLoadBalancerPorts(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error)
	UpdateLoadBalancerHealthCheck(loadBalancerID string, healthCheck models.HealthCheck) (*models.LoadBalancer, error)
}
