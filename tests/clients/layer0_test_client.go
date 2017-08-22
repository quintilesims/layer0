package clients

import (
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/common/models"
)

type Tester interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Layer0TestClient struct {
	T      Tester
	Client *client.APIClient
}

func NewLayer0TestClient(t Tester, endpoint, token string) *Layer0TestClient {
	return &Layer0TestClient{
		T: t,
		Client: client.NewAPIClient(client.Config{
			Endpoint: endpoint,
			Token:    token,
		}),
	}
}

func (l *Layer0TestClient) CreateTask(taskName, environmentID, deployID string, copies int, overrides []models.ContainerOverride) string {
	jobID, err := l.Client.CreateTask(taskName, environmentID, deployID, copies, overrides)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateEnvironment(name string) *models.Environment {
	environment, err := l.Client.CreateEnvironment(name, "m3.medium", 0, nil, "linux", "")
	if err != nil {
		l.T.Fatal(err)
	}

	return environment
}

func (l *Layer0TestClient) CreateDeploy(name string, content []byte) *models.Deploy {
	deploy, err := l.Client.CreateDeploy(name, content)
	if err != nil {
		l.T.Fatal(err)
	}

	return deploy
}

func (l *Layer0TestClient) CreateLoadBalancer(name, environmentID string) *models.LoadBalancer {
	hc := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           10,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	ports := []models.Port{{HostPort: 80, ContainerPort: 80, Protocol: "http"}}

	loadBalancer, err := l.Client.CreateLoadBalancer(name, environmentID, hc, ports, true)
	if err != nil {
		l.T.Fatal(err)
	}

	return loadBalancer
}

func (l *Layer0TestClient) CreateService(name, environmentID, deployID, loadBalancerID string) *models.Service {
	service, err := l.Client.CreateService(name, environmentID, deployID, loadBalancerID)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) GetService(id string) *models.Service {
	service, err := l.Client.GetService(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) GetTask(id string) *models.Task {
	task, err := l.Client.GetTask(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return task
}

func (l *Layer0TestClient) GetEnvironment(id string) *models.Environment {
	environment, err := l.Client.GetEnvironment(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return environment
}

func (l *Layer0TestClient) ListEnvironments() []*models.EnvironmentSummary {
	environments, err := l.Client.ListEnvironments()
	if err != nil {
		l.T.Fatal(err)
	}

	return environments
}

func (l *Layer0TestClient) ListLoadBalancers() []*models.LoadBalancerSummary {
	loadbalancers, err := l.Client.ListLoadBalancers()
	if err != nil {
		l.T.Fatal(err)
	}

	return loadbalancers
}

func (l *Layer0TestClient) ListDeploys() []*models.DeploySummary {
	deploys, err := l.Client.ListDeploys()
	if err != nil {
		l.T.Fatal(err)
	}

	return deploys
}

func (l *Layer0TestClient) ListServices() []*models.ServiceSummary {
	services, err := l.Client.ListServices()
	if err != nil {
		l.T.Fatal(err)
	}

	return services
}

func (l *Layer0TestClient) ListTasks() []*models.TaskSummary {
	tasks, err := l.Client.ListTasks()
	if err != nil {
		l.T.Fatal(err)
	}

	return tasks
}

func (l *Layer0TestClient) ListJobs() []*models.Job {
	jobs, err := l.Client.ListJobs()
	if err != nil {
		l.T.Fatal(err)
	}
	return jobs
}

func (l *Layer0TestClient) ScaleService(id string, scale int) *models.Service {
	service, err := l.Client.ScaleService(id, scale)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) CreateLink(id1, id2 string) {
	if err := l.Client.CreateLink(id1, id2); err != nil {
		l.T.Fatal(err)
	}
}

func (l *Layer0TestClient) DeleteLink(id1, id2 string) {
	if err := l.Client.DeleteLink(id1, id2); err != nil {
		l.T.Fatal(err)
	}
}
