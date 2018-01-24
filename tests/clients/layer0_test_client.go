package clients

import (
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

type Layer0TestClient struct {
	T      testutils.Tester
	Client *client.APIClient
}

func NewLayer0TestClient(t testutils.Tester, endpoint, token string) *Layer0TestClient {
	return &Layer0TestClient{
		T: t,
		Client: client.NewAPIClient(client.Config{
			Endpoint: endpoint,
			Token:    token,
		}),
	}
}

func (l *Layer0TestClient) CreateTask(taskName, environmentID, deployID string, copies int, overrides []models.ContainerOverride) string {
	req := models.CreateTaskRequest{
		DeployID:      deployID,
		EnvironmentID: environmentID,
		TaskName:      taskName,
	}

	jobID, err := l.Client.CreateTask(req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) CreateEnvironment(name string) string {
	req := models.CreateEnvironmentRequest{
		EnvironmentName: name,
		InstanceType:    "m3.medium",
		MinScale:        0,
		OperatingSystem: "linux",
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
	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: name,
		EnvironmentID:    environmentID,
		IsPublic:         true,
		Ports:            []models.Port{{HostPort: 80, ContainerPort: 80, Protocol: "http"}},
		HealthCheck: models.HealthCheck{
			Target:             "TCP:80",
			Interval:           10,
			Timeout:            5,
			HealthyThreshold:   2,
			UnhealthyThreshold: 2,
		},
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

func (l *Layer0TestClient) ReadEnvironment(id string) *models.Environment {
	environment, err := l.Client.ReadEnvironment(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return environment
}

func (l *Layer0TestClient) ReadLoadBalancer(id string) *models.LoadBalancer {
	loadBalancer, err := l.Client.ReadLoadBalancer(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return loadBalancer
}

func (l *Layer0TestClient) ReadDeploy(id string) *models.Deploy {
	deploy, err := l.Client.ReadDeploy(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return deploy
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

func (l *Layer0TestClient) UpdateEnvironment(environmentID string, minScale, maxScale int, links []string) string {
	req := models.UpdateEnvironmentRequest{
		MinScale: &minScale,
		MaxScale: &maxScale,
		Links:    &links,
	}

	jobID, err := l.Client.UpdateEnvironment(environmentID, req)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
}

func (l *Layer0TestClient) UpdateLoadBalancer(loadBalancerID string) string {
	req := models.UpdateLoadBalancerRequest{
		Ports: &[]models.Port{{HostPort: 81, ContainerPort: 81, Protocol: "tcp"}},
		HealthCheck: &models.HealthCheck{
			Target:             "TCP:81",
			Interval:           15,
			Timeout:            10,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
	}

	jobID, err := l.Client.UpdateLoadBalancer(loadBalancerID, req)
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
