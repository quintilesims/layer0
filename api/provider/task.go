package provider

import "github.com/quintilesims/layer0/common/models"

type TaskProvider interface {
	// todo: all creates should return a string
	Create(req models.CreateTaskRequest) (string, error)
	Read(taskID string) (*models.Task, error)
	List() ([]models.TaskSummary, error)
	Delete(taskID string) error
}
