package clients

import (
	"fmt"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/sts/client"
	"github.com/quintilesims/sts/models"
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
	logrus.Debugf("Waiting for sts service to be healthy")
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 10) {
		if _, err := s.Client.GetHealth(); err != nil {
			logrus.Debug(err)
		}
	}

	fmt.Errorf("[ERROR] Timeout reached after %v", timeout)
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

func (s *STSTestClient) RunCommand(args ...string) (string, error) {
	command, err := s.Client.CreateCommand(args[0], args...)
	if err != nil {
		return "", err
	}

	return command.Output, nil
}
