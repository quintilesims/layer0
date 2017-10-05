package provider

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type TaskProvider interface {
	Create(req models.CreateTaskRequest) (string, error)
	Delete(taskID string) error
	List() ([]models.TaskSummary, error)
	Logs(taskID string, tail int, start, end time.Time) ([]models.LogFile, error)
	Read(taskID string) (*models.Task, error)
}
