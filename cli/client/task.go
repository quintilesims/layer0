package client

import (
	"fmt"
	"github.com/quintilesims/layer0/common/models"
)

func (c *APIClient) CreateTask(
	name string,
	environmentID string,
	deployID string,
	copies int,
	overrides []models.ContainerOverride,
) (string, error) {
	req := models.CreateTaskRequest{
		TaskName:           name,
		EnvironmentID:      environmentID,
		DeployID:           deployID,
		Copies:             copies,
		ContainerOverrides: overrides,
	}

	jobID, err := c.ExecuteWithJob(c.Sling("task/").Post("").BodyJSON(req))
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (c *APIClient) DeleteTask(id string) error {
	var response *string
	if err := c.Execute(c.Sling("task/").Delete(id), &response); err != nil {
		return err
	}

	return nil
}

func (c *APIClient) GetTask(id string) (*models.Task, error) {
	var task *models.Task
	if err := c.Execute(c.Sling("task/").Get(id), &task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *APIClient) GetTaskLogs(id string, tail int) ([]*models.LogFile, error) {
	url := id + "/logs"
	if tail > 0 {
		url = fmt.Sprintf("%s?tail=%d", url, tail)
	}

	var logFiles []*models.LogFile
	if err := c.Execute(c.Sling("task/").Get(url), &logFiles); err != nil {
		return nil, err
	}

	return logFiles, nil
}

func (c *APIClient) ListTasks() ([]*models.TaskSummary, error) {
	var tasks []*models.TaskSummary
	if err := c.Execute(c.Sling("task/").Get(""), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
