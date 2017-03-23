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

func (l *Layer0TestClient) GetService(id string) *models.Service {
	service, err := l.Client.GetService(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}

func (l *Layer0TestClient) GetEnvironment(id string) *models.Environment {
	environment, err := l.Client.GetEnvironment(id)
	if err != nil {
		l.T.Fatal(err)
	}

	return environment
}

func (l *Layer0TestClient) ScaleService(id string, scale int) *models.Service {
	service, err := l.Client.ScaleService(id, scale)
	if err != nil {
		l.T.Fatal(err)
	}

	return service
}
