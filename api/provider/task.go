package provider

import "github.com/quintilesims/layer0/common/models"

// todo: all create methods should return a string
type TaskProvider interface {
	Create(req models.CreateTaskRequest) (string, error)
	Delete(taskID string) error
	List() ([]models.TaskSummary, error)
	Read(taskID string) (*models.Task, error)
}
