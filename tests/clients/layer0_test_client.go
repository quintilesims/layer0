package clients

import (
	"testing"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
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

func (l *Layer0TestClient) jobHelper(fn func() (string, error)) string {
	jobID, err := fn()
	if err != nil {
		l.T.Fatal(err)
	}

	job, err := client.WaitForJob(l.Client, jobID, config.DEFAULT_JOB_EXPIRY)
	if err != nil {
		l.T.Fatal(err)
	}

	return job.Result
}

func (l *Layer0TestClient) CreateTask(req models.CreateTaskRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.CreateTask(req)
	})
}

func (l *Layer0TestClient) CreateEnvironment(req models.CreateEnvironmentRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.CreateEnvironment(req)
	})
}

func (l *Layer0TestClient) CreateDeploy(req models.CreateDeployRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.CreateDeploy(req)
	})
}

func (l *Layer0TestClient) CreateLoadBalancer(req models.CreateLoadBalancerRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.CreateLoadBalancer(req)
	})
}

func (l *Layer0TestClient) CreateService(req models.CreateServiceRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.CreateService(req)
	})
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

func (l *Layer0TestClient) UpdateEnvironment(environmentID string, req models.UpdateEnvironmentRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.UpdateEnvironment(environmentID, req)
	})
}

func (l *Layer0TestClient) UpdateService(serviceID string, req models.UpdateServiceRequest) string {
	return l.jobHelper(func() (string, error) {
		return l.Client.UpdateService(serviceID, req)
	})
}
