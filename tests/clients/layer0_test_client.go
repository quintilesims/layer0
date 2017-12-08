package clients

import (
	"testing"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
)

type Layer0TestClient struct {
	T      *testing.T
	Client *client.APIClient
}

func NewLayer0TestClient(t *testing.T, endpoint, token string) *Layer0TestClient {
	return &Layer0TestClient{
		T: t,
		Client: client.NewAPIClient(client.Config{
			Endpoint: endpoint,
			Token:    token,
		}),
	}
}

func (l *Layer0TestClient) CreateTask(taskName, environmentID, deployID string, overrides []models.ContainerOverride) string {
	req := models.CreateTaskRequest{
		ContainerOverrides: overrides,
		TaskName:           taskName,
		EnvironmentID:      environmentID,
		DeployID:           deployID,
	}

	jobID, err := l.Client.CreateTask(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateEnvironment(name string) string {
	req := models.CreateEnvironmentRequest{
		EnvironmentName:  name,
		InstanceSize:     "m3.medium",
		UserDataTemplate: nil,
		MinClusterCount:  0,
		OperatingSystem:  "linux",
		AMIID:            "",
	}

	jobID, err := l.Client.CreateEnvironment(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateDeploy(name string, content []byte) string {
	req := models.CreateDeployRequest{
		DeployName: name,
		DeployFile: content,
	}

	jobID, err := l.Client.CreateDeploy(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateLoadBalancer(name, environmentID string) string {
	hc := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           10,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	ports := []models.Port{{HostPort: 80, ContainerPort: 80, Protocol: "http"}}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: name,
		EnvironmentID:    environmentID,
		IsPublic:         true,
		Ports:            ports,
		HealthCheck:      hc,
	}

	jobID, err := l.Client.CreateLoadBalancer(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateService(name, environmentID, deployID, loadBalancerID string) string {
	req := models.CreateServiceRequest{
		DeployID:       deployID,
		EnvironmentID:  environmentID,
		LoadBalancerID: loadBalancerID,
		ServiceName:    name,
	}

	jobID, err := l.Client.CreateService(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) ReadDeploy(id string) *models.Deploy {
	deploy, err := l.Client.ReadDeploy(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return deploy
}

func (l *Layer0TestClient) ReadLoadBalancer(id string) *models.LoadBalancer {
	loadBalancer, err := l.Client.ReadLoadBalancer(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return loadBalancer
}

func (l *Layer0TestClient) ReadService(id string) *models.Service {
	service, err := l.Client.ReadService(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) ReadTask(id string) *models.Task {
	task, err := l.Client.ReadTask(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return task
}

func (l *Layer0TestClient) ReadEnvironment(id string) *models.Environment {
	environment, err := l.Client.ReadEnvironment(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return environment
}

func (l *Layer0TestClient) ListTasks() []*models.TaskSummary {
	tasks, err := l.Client.ListTasks()
	if err != nil {
		l.T.Fatal(err)
	}

	return tasks
}

func (l *Layer0TestClient) UpdateEnvironmentLink(environmentID string, links []string) string {
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	jobID, err := l.Client.UpdateEnvironment(environmentID, req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) UpdateService(serviceID, deployID string, scale int) string {
	req := models.UpdateServiceRequest{
		DeployID: &deployID,
		Scale:    &scale,
	}

	jobID, err := l.Client.UpdateService(serviceID, req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}
