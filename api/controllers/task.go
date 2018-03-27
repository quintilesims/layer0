package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type TaskController struct {
	TaskProvider provider.TaskProvider
}

func NewTaskController(t provider.TaskProvider) *TaskController {
	return &TaskController{
		TaskProvider: t,
	}
}

func (t *TaskController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/task",
			Handlers: fireball.Handlers{
				"GET":  t.listTasks,
				"POST": t.createTask,
			},
		},
		{
			Path: "/task/:id",
			Handlers: fireball.Handlers{
				"GET":    t.readTask,
				"DELETE": t.deleteTask,
			},
		},
		{
			Path: "/task/:id/logs",
			Handlers: fireball.Handlers{
				"GET": t.readTaskLogs,
			},
		},
	}
}

func (t *TaskController) createTask(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	taskID, err := t.TaskProvider.Create(req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateEntityResponse{EntityID: taskID})
}

func (t *TaskController) deleteTask(c *fireball.Context) (fireball.Response, error) {
	taskID := c.PathVariables["id"]
	if err := t.TaskProvider.Delete(taskID); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (t *TaskController) listTasks(c *fireball.Context) (fireball.Response, error) {
	tasks, err := t.TaskProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, tasks)
}

func (t *TaskController) readTask(c *fireball.Context) (fireball.Response, error) {
	taskID := c.PathVariables["id"]
	task, err := t.TaskProvider.Read(taskID)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, task)
}

func (t *TaskController) readTaskLogs(c *fireball.Context) (fireball.Response, error) {
	taskID := c.PathVariables["id"]
	tail, start, end, err := client.ParseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	logs, err := t.TaskProvider.Logs(taskID, tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}
