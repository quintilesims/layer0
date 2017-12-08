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

func (l *Layer0TestClient) CreateTask(req models.CreateTaskRequest) string {
	jobID, err := l.Client.CreateTask(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateEnvironment(req models.CreateEnvironmentRequest) string {
	jobID, err := l.Client.CreateEnvironment(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateDeploy(req models.CreateDeployRequest) string {
	jobID, err := l.Client.CreateDeploy(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateLoadBalancer(req models.CreateLoadBalancerRequest) string {
	jobID, err := l.Client.CreateLoadBalancer(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateService(req models.CreateServiceRequest) string {
	jobID, err := l.Client.CreateService(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) ReadDeploy(deployID string) *models.Deploy {
	deploy, err := l.Client.ReadDeploy(deployID)
	if err != nil {
		l.T.Fatal(err)
	}

	return deploy
}

func (l *Layer0TestClient) ReadLoadBalancer(loadBalancerID string) *models.LoadBalancer {
	loadBalancer, err := l.Client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		l.T.Fatal(err)
	}

	return loadBalancer
}

func (l *Layer0TestClient) ReadService(serviceID string) *models.Service {
	service, err := l.Client.ReadService(serviceID)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) ReadTask(taskID string) *models.Task {
	task, err := l.Client.ReadTask(taskID)
	if err != nil {
		l.T.Fatal(err)
	}

	return task
}

func (l *Layer0TestClient) ReadEnvironment(environmentID string) *models.Environment {
	environment, err := l.Client.ReadEnvironment(environmentID)
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

func (l *Layer0TestClient) UpdateEnvironmentLink(environmentID string, req models.UpdateEnvironmentRequest) string {
	jobID, err := l.Client.UpdateEnvironment(environmentID, req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) UpdateService(serviceID string, req models.UpdateServiceRequest) string {
	jobID, err := l.Client.UpdateService(serviceID, req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}
