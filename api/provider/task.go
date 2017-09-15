package provider

import "github.com/quintilesims/layer0/common/models"

type TaskProvider interface {
	Create(req models.CreateTaskRequest) (*models.Task, error)
	Delete(taskID string) error
	List() ([]models.TaskSummary, error)
	Read(taskID string) (*models.Task, error)
}
