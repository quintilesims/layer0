package aws

import (
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/models"
)

type TaskProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
}

func NewTaskProvider(a *awsc.Client, t tag.Store) *TaskProvider {
	return &TaskProvider{
		AWS:      a,
		TagStore: t,
	}
}

func (t *TaskProvider) Create(req models.CreateTaskRequest) (*models.Task, error) {
	return nil, nil
}

func (t *TaskProvider) Read(TaskID string) (*models.Task, error) {
	return nil, nil
}

func (t *TaskProvider) List() ([]models.TaskSummary, error) {
	return nil, nil
}

func (t *TaskProvider) Delete(TaskID string) error {
	return nil
}
