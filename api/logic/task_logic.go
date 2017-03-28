package logic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type TaskLogic interface {
	CreateTask(models.CreateTaskRequest) (*models.Task, error)
	ListTasks() ([]*models.TaskSummary, error)
	GetTask(string) (*models.Task, error)
	DeleteTask(string) error
	GetTaskLogs(string, int) ([]*models.LogFile, error)
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

	createf := func() (*models.Task, error) {
		return this.Backend.CreateTask(
			req.EnvironmentID,
			req.TaskName,
			req.DeployID,
			int(req.Copies),
			req.ContainerOverrides)
	}

	return this.createTask(req, createf)
}

func (this *L0TaskLogic) createTask(req models.CreateTaskRequest, createf func() (*models.Task, error)) (*models.Task, error) {
	task, err := createf()
	if err != nil {
		switch err := err.(type) {
		case *ecsbackend.PartialCreateTaskFailure:
			return task, this.handlePartialFailure(req, task, err)
		default:
			return nil, err
		}
	}

	return this.createTaskTags(req, task)
}

func (this *L0TaskLogic) handlePartialFailure(req models.CreateTaskRequest, task *models.Task, partialFailure *ecsbackend.PartialCreateTaskFailure) *ecsbackend.PartialCreateTaskFailure {
	// first, try to add tags for the partial successes
	if task != nil {
		if _, err := this.createTaskTags(req, task); err != nil {
			log.Errorf("Failed to create task tags for %s: %v", req.TaskName, err)
		}
	}

	// next, wrap the backend retry function to include the required
	// tag logic
	wrappedError := &ecsbackend.PartialCreateTaskFailure{
		NumFailed: partialFailure.NumFailed,
		Retry: func() (*models.Task, error) {
			return this.createTask(req, partialFailure.Retry)
		},
	}

	return wrappedError
}

func (this *L0TaskLogic) createTaskTags(req models.CreateTaskRequest, task *models.Task) (*models.Task, error) {
	taskID := task.TaskID
	if err := this.upsertTagf(taskID, "task", "name", req.TaskName); err != nil {
		return task, err
	}

	environmentID := req.EnvironmentID
	if err := this.upsertTagf(taskID, "task", "environment_id", environmentID); err != nil {
		return task, err
	}

	deployID := req.DeployID
	if err := this.upsertTagf(taskID, "task", "deploy_id", deployID); err != nil {
		return task, err
	}

	if err := this.populateModel(task); err != nil {
		return task, err
	}

	return task, nil
}

func (this *L0TaskLogic) GetTaskLogs(taskID string, tail int) ([]*models.LogFile, error) {
	environmentID, err := this.getEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}

	logs, err := this.Backend.GetTaskLogs(environmentID, taskID, tail)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (this *L0TaskLogic) getEnvironmentID(taskID string) (string, error) {
	tags, err := this.TagStore.SelectByQuery("task", taskID)
	if err != nil {
		return "", err
	}

	if tag := tags.WithKey("environment_id").First(); tag != nil {
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
	tags, err := this.TagStore.SelectByQuery("task", model.TaskID)
	if err != nil {
		return err
	}

	if tag := tags.WithKey("environment_id").First(); tag != nil {
		model.EnvironmentID = tag.Value
	}

	if tag := tags.WithKey("deploy_id").First(); tag != nil {
		model.DeployID = tag.Value
	}

	if tag := tags.WithKey("name").First(); tag != nil {
		model.TaskName = tag.Value
	}

	if model.EnvironmentID != "" {
		tags, err := this.TagStore.SelectByQuery("environment", model.EnvironmentID)
		if err != nil {
			return err
		}

		if tag := tags.WithKey("name").First(); tag != nil {
			model.EnvironmentName = tag.Value
		}
	}

	if model.DeployID != "" {
		tags, err := this.TagStore.SelectByQuery("deploy", model.DeployID)
		if err != nil {
			return err
		}

		if tag := tags.WithKey("name").First(); tag != nil {
			model.DeployName = tag.Value
		}

		if tag := tags.WithKey("version").First(); tag != nil {
			model.DeployVersion = tag.Value
		}
	}

	return nil
}
