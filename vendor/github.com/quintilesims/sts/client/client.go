package client

import (
	"fmt"
	"github.com/quintilesims/sts/models"
	"github.com/zpatrick/rclient"
)

type APIClient struct {
	client *rclient.RestClient
}

func New(url string) (*APIClient, error) {
	client, err := rclient.NewRestClient(url)
	if err != nil {
		return nil, err
	}

	apiClient := &APIClient{
		client: client,
	}

	return apiClient, nil
}

func (a *APIClient) CreateCommand(name string, args ...string) (*models.Command, error) {
	req := models.CreateCommandRequest{
		Name: name,
		Args: args,
	}

	var command *models.Command
	if err := a.client.Post("/command", req, &command); err != nil {
		return nil, err
	}

	return command, nil
}

func (a *APIClient) GetCommand(name string) (*models.Command, error) {
	var command *models.Command
	if err := a.client.Get(fmt.Sprintf("/command/%s", name), &command); err != nil {
		return nil, err
	}

	return command, nil
}

func (a *APIClient) GetHealth() (*models.Health, error) {
	var health *models.Health
	if err := a.client.Get("/health", &health); err != nil {
		return nil, err
	}

	return health, nil
}

func (a *APIClient) ListCommands() ([]*models.Command, error) {
	var commands []*models.Command
	if err := a.client.Get("/command", &commands); err != nil {
		return nil, err
	}

	return commands, nil
}

func (a *APIClient) SetHealth(mode string) (*models.Health, error) {
	req := models.SetHealthRequest{
		Mode: mode,
	}

	var health *models.Health
	if err := a.client.Post("/health", req, &health); err != nil {
		return nil, err
	}

	return health, nil
}
