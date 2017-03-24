package clients

import (
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/sts/client"
	"github.com/quintilesims/sts/models"
	"testing"
	"time"
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

func (s *STSTestClient) WaitForHealthy(timeout time.Duration) {
	testutils.WaitFor(s.T, time.Second*10, timeout, func() bool {
		logrus.Printf("Waiting for sts service to be healthy")
		if _, err := s.Client.GetHealth(); err != nil {
			return false
		}

		return true
	})
}

func (s *STSTestClient) GetHealth() *models.Health {
	health, err := s.Client.GetHealth()
	if err != nil {
		s.T.Fatal(err)
	}

	return health
}

func (s *STSTestClient) SetHealth(mode string) *models.Health {
	health, err := s.Client.SetHealth(mode)
	if err != nil {
		s.T.Fatal(err)
	}

	return health
}
