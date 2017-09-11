package provider

import "github.com/quintilesims/layer0/common/models"

type TaskProvider interface {
	Create(req models.CreateTaskRequest) (*models.Task, error)
	Read(taskID string) (*models.Task, error)
	List() ([]models.TaskSummary, error)
	Delete(taskID string) error
	Update(taskID string) error
}
