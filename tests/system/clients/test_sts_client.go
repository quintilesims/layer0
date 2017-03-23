package clients

import (
	"github.com/quintilesims/sts/client"
	"github.com/quintilesims/sts/models"
	"testing"
)

type STSTestClient struct {
	T      *testing.T
	Client *client.APIClient
}

func NewSTSTestClient(t *testing.T, url string) *STSTestClient {
	stsClient, err := client.New(url)
	if err != nil {
		t.Fatal(err)
	}

	return &STSTestClient{
		T:      t,
		Client: stsClient,
	}
}

func (l *STSTestClient) GetHealth() *models.Health {
	health, err := l.Client.GetHealth()
	if err != nil {
		l.T.Fatal(err)
	}

	return health
}

func (l *STSTestClient) SetHealth(mode string) *models.Health {
	health, err := l.Client.SetHealth(mode)
	if err != nil {
		l.T.Fatal(err)
	}

	return health
}
