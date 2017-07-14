package logic

import (
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type TaskLogic interface {
	CreateTask(models.CreateTaskRequest) (*models.Task, error)
	ListTasks() ([]*models.TaskSummary, error)
	GetTask(string) (*models.Task, error)
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
	tasks, err := this.Backend.ListTasks()
	if err != nil {
		return nil, err
	}

	summaries := make([]*models.TaskSummary, len(tasks))
	for i, task := range tasks {
		if err := this.populateModel(task); err != nil {
			return nil, err
		}

		summaries[i] = &models.TaskSummary{
			TaskID:          task.TaskID,
			TaskName:        task.TaskName,
			EnvironmentID:   task.EnvironmentID,
			EnvironmentName: task.EnvironmentName,
		}
	}

	return summaries, nil
}

func (this *L0TaskLogic) GetTask(taskID string) (*models.Task, error) {
	environmentID, err := this.getEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}

	task, err := this.Backend.GetTask(environmentID, taskID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.InvalidTaskID {
			return nil, errors.Newf(errors.InvalidTaskID, "Task %s does not exist", taskID)
		}

		return nil, err
	}

	if err := this.populateModel(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (this *L0TaskLogic) DeleteTask(taskID string) error {
	environmentID, err := this.getEnvironmentID(taskID)
	if err != nil {
		return err
	}

	if err := this.Backend.DeleteTask(environmentID, taskID); err != nil {
		return err
	}

	if err := this.deleteEntityTags("task", taskID); err != nil {
		return err
	}

	return nil
}

func (this *L0TaskLogic) CreateTask(req models.CreateTaskRequest) (*models.Task, error) {
	if req.EnvironmentID == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.DeployID == "" {
		return nil, errors.Newf(errors.MissingParameter, "DeployID not specified")
	}

	if req.TaskName == "" {
		return nil, errors.Newf(errors.MissingParameter, "TaskName not specified")
	}

	task, err := this.Backend.CreateTask(
		req.EnvironmentID,
		req.TaskName,
		req.DeployID,
		req.ContainerOverrides)
	if err != nil {
		return nil, err
	}

	taskID := task.TaskID
	if err := this.TagStore.Insert(models.Tag{EntityID: taskID, EntityType: "task", Key: "name", Value: req.TaskName}); err != nil {
		return task, err
	}

	environmentID := req.EnvironmentID
	if err := this.TagStore.Insert(models.Tag{EntityID: taskID, EntityType: "task", Key: "environment_id", Value: environmentID}); err != nil {
		return task, err
	}

	deployID := req.DeployID
	if err := this.TagStore.Insert(models.Tag{EntityID: taskID, EntityType: "task", Key: "deploy_id", Value: deployID}); err != nil {
		return task, err
	}

	if err := this.populateModel(task); err != nil {
		return task, err
	}

	return task, nil
}

func (this *L0TaskLogic) GetTaskLogs(taskID, start, end string, tail int) ([]*models.LogFile, error) {
	environmentID, err := this.getEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}

	logs, err := this.Backend.GetTaskLogs(environmentID, taskID, start, end, tail)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (this *L0TaskLogic) getEnvironmentID(taskID string) (string, error) {
	tags, err := this.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		return tag.Value, nil
	}

	tasks, err := this.ListTasks()
	if err != nil {
		return "", err
	}

	for _, task := range tasks {
		if task.TaskID == taskID {
			return task.EnvironmentID, nil
		}
	}

	return "", errors.Newf(errors.InvalidTaskID, "Task %s does not exist", taskID)
}

func (this *L0TaskLogic) populateModel(model *models.Task) error {
	tags, err := this.TagStore.SelectByTypeAndID("task", model.TaskID)
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
