package logic

import (
	"fmt"

	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type TaskLogic interface {
	CreateTask(models.CreateTaskRequest) (string, error)
	ListTasks() ([]*models.TaskSummary, error)
	GetTask(string) (*models.Task, error)
	GetEnvironmentTasks(environmentID string) ([]*models.Task, error)
	DeleteTask(string) error
	GetTaskLogs(string, string, string, int) ([]*models.LogFile, error)
}

type L0TaskLogic struct {
	Logic
}

func NewL0TaskLogic(logic Logic) *L0TaskLogic {
	return &L0TaskLogic{
		Logic: logic,
	}
}

func (this *L0TaskLogic) ListTasks() ([]*models.TaskSummary, error) {
	taskARNs, err := this.Backend.ListTasks()
	if err != nil {
		return nil, err
	}

	return this.makeTaskSummaryModels(taskARNs)
}

func (this *L0TaskLogic) GetTask(taskID string) (*models.Task, error) {
	environmentID, err := this.lookupTaskEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}

	taskARN, err := this.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	taskModel, err := this.Backend.GetTask(environmentID, taskARN)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.InvalidTaskID {
			return nil, errors.Newf(errors.InvalidTaskID, "Task %s does not exist", taskID)
		}

		return nil, err
	}

	if err := this.populateModel(taskID, taskModel); err != nil {
		return nil, err
	}

	return taskModel, nil
}

func (this *L0TaskLogic) GetEnvironmentTasks(environmentID string) ([]*models.Task, error) {
	taskARNModels, err := this.Backend.GetEnvironmentTasks(environmentID)
	if err != nil {
		return nil, err
	}

	taskModels := []*models.Task{}
	for taskARN, taskModel := range taskARNModels {
		taskID, err := this.getTaskARNFromID(taskARN)
		if err != nil {
			return nil, err
		}

		if err := this.populateModel(taskID, taskModel); err != nil {
			return nil, err
		}

		taskModels = append(taskModels, taskModel)
	}

	return taskModels, nil
}

func (this *L0TaskLogic) DeleteTask(taskID string) error {
	environmentID, err := this.lookupTaskEnvironmentID(taskID)
	if err != nil {
		return err
	}

	taskARN, err := this.lookupTaskARN(taskID)
	if err != nil {
		return err
	}

	if err := this.Backend.DeleteTask(environmentID, taskARN); err != nil {
		return err
	}

	if err := this.deleteEntityTags("task", taskID); err != nil {
		return err
	}

	return nil
}

func (this *L0TaskLogic) CreateTask(req models.CreateTaskRequest) (string, error) {
	if req.EnvironmentID == "" {
		return "", errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.DeployID == "" {
		return "", errors.Newf(errors.MissingParameter, "DeployID not specified")
	}

	if req.TaskName == "" {
		return "", errors.Newf(errors.MissingParameter, "TaskName not specified")
	}

	taskARN, err := this.Backend.CreateTask(req.EnvironmentID, req.DeployID, req.ContainerOverrides)
	if err != nil {
		return "", err
	}

	taskID := id.GenerateHashedEntityID(req.TaskName)
	tags := []models.Tag{
		{EntityID: taskID, EntityType: "task", Key: "name", Value: req.TaskName},
		{EntityID: taskID, EntityType: "task", Key: "environment_id", Value: req.EnvironmentID},
		{EntityID: taskID, EntityType: "task", Key: "deploy_id", Value: req.DeployID},
		{EntityID: taskID, EntityType: "task", Key: "arn", Value: taskARN},
	}

	for _, tag := range tags {
		if err := this.TagStore.Insert(tag); err != nil {
			return "", err
		}
	}

	return taskID, nil
}

func (this *L0TaskLogic) GetTaskLogs(taskID, start, end string, tail int) ([]*models.LogFile, error) {
	environmentID, err := this.lookupTaskEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}

	taskARN, err := this.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	logs, err := this.Backend.GetTaskLogs(environmentID, taskARN, start, end, tail)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (t *L0TaskLogic) getTaskARNFromID(taskARN string) (string, error) {
	tags, err := t.TagStore.SelectByType("task")
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("arn").WithValue(taskARN).First(); ok {
		return tag.EntityID, nil
	}

	return "", fmt.Errorf("Failed to find task id for ARN %s", taskARN)
}

func (this *L0TaskLogic) lookupTaskEnvironmentID(taskID string) (string, error) {
	tags, err := this.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		return tag.Value, nil
	}

	return "", errors.Newf(errors.TaskDoesNotExist, "Failed to find environment for task %s", taskID)
}

func (t *L0TaskLogic) lookupTaskARN(taskID string) (string, error) {
	tags, err := t.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.Newf(errors.TaskDoesNotExist, "Task '%s' does not exist", taskID)
	}

	if tag, ok := tags.WithKey("arn").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Failed to find ARN for task '%s'", taskID)
}

func (t *L0TaskLogic) makeTaskSummaryModels(taskARNs []string) ([]*models.TaskSummary, error) {
	environmentTags, err := t.TagStore.SelectByType("environment")
	if err != nil {
		return nil, err
	}

	taskTags, err := t.TagStore.SelectByType("task")
	if err != nil {
		return nil, err
	}

	taskARNMatches := map[string]bool{}
	for _, taskARN := range taskARNs {
		taskARNMatches[taskARN] = true
	}

	taskModels := make([]*models.TaskSummary, 0, len(taskARNs))
	for _, tag := range taskTags.WithKey("arn") {
		if taskARNMatches[tag.Value] {
			model := &models.TaskSummary{
				TaskID: tag.EntityID,
			}

			if tag, ok := taskTags.WithID(model.TaskID).WithKey("name").First(); ok {
				model.TaskName = tag.Value
			}

			if tag, ok := taskTags.WithID(model.TaskID).WithKey("environment_id").First(); ok {
				model.EnvironmentID = tag.Value

				if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
					model.EnvironmentName = t.Value
				}
			}

			taskModels = append(taskModels, model)
		}
	}

	return taskModels, nil
}

func (this *L0TaskLogic) populateModel(taskID string, model *models.Task) error {
	model.TaskID = taskID

	tags, err := this.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		model.EnvironmentID = tag.Value
	}

	if tag, ok := tags.WithKey("deploy_id").First(); ok {
		model.DeployID = tag.Value
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.TaskName = tag.Value
	}

	if model.EnvironmentID != "" {
		tags, err := this.TagStore.SelectByTypeAndID("environment", model.EnvironmentID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.EnvironmentName = tag.Value
		}
	}

	if model.DeployID != "" {
		tags, err := this.TagStore.SelectByTypeAndID("deploy", model.DeployID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.DeployName = tag.Value
		}

		if tag, ok := tags.WithKey("version").First(); ok {
			model.DeployVersion = tag.Value
		}
	}

	return nil
}
