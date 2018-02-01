package client

import (
	"fmt"
	"net/url"

	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/rclient"
)

func (c *APIClient) CreateTask(req models.CreateTaskRequest) (string, error) {
	var resp models.CreateEntityResponse
	if err := c.client.Post("/task", req, &resp); err != nil {
		return "", err
	}

	return resp.EntityID, nil
}

func (c *APIClient) DeleteTask(taskID string) error {
	path := fmt.Sprintf("/task/%s", taskID)
	if err := c.client.Delete(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ListTasks() ([]models.TaskSummary, error) {
	var tasks []models.TaskSummary
	if err := c.client.Get("/task", &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *APIClient) ReadTask(taskID string) (*models.Task, error) {
	var task *models.Task
	path := fmt.Sprintf("/task/%s", taskID)
	if err := c.client.Get(path, &task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *APIClient) ReadTaskLogs(taskID string, query url.Values) ([]models.LogFile, error) {
	var logs []models.LogFile
	path := fmt.Sprintf("/task/%s/logs", taskID)
	if err := c.client.Get(path, &logs, rclient.Query(query)); err != nil {
		return nil, err
	}

	return logs, nil
}
