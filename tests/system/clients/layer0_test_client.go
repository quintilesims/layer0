package clients

import (
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/common/models"
	"testing"
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

func (l *Layer0TestClient) CreateTask(taskName, environmentID, deployID string, copies int, overrides []models.ContainerOverride) string {
	jobID, err := l.Client.CreateTask(taskName, environmentID, deployID, copies, overrides)
	if err != nil {
		l.T.Fatal(err)
	}

	return jobID
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

func (l *Layer0TestClient) ListTasks() []*models.TaskSummary {
	tasks, err := l.Client.ListTasks()
	if err != nil {
		l.T.Fatal(err)
	}

	return tasks
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
