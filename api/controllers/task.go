package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

const (
	// 'YYYY-MM-DD HH:MM' time layout as described by https://golang.org/src/time/format.go
	TIME_LAYOUT = "2006-01-02 15:04"
)

type TaskController struct {
	TaskProvider provider.TaskProvider
	JobStore     job.Store
	TagStore     tag.Store
}

func NewTaskController(p provider.TaskProvider, j job.Store, t tag.Store) *TaskController {
	return &TaskController{
		TaskProvider: p,
		JobStore:     j,
		TagStore:     t,
	}
}

func (t *TaskController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/task",
			Handlers: fireball.Handlers{
				"GET":  t.ListTasks,
				"POST": t.CreateTask,
			},
		},
		{
			Path: "/task/:id",
			Handlers: fireball.Handlers{
				"GET":    t.GetTask,
				"DELETE": t.DeleteTask,
			},
		},
		{
			Path: "/task/:id/logs",
			Handlers: fireball.Handlers{
				"GET": t.GetTaskLogs,
			},
		},
	}
}

func (t *TaskController) CreateTask(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	return createJob(t.TagStore, t.JobStore, models.CreateTaskJob, req)
}

func (t *TaskController) DeleteTask(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return createJob(t.TagStore, t.JobStore, models.DeleteTaskJob, id)
}

func (t *TaskController) GetTask(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := t.TaskProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (t *TaskController) GetTaskLogs(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]

	tail, start, end, err := parseLoggingQuery(c.Request.URL.Query())
	if err != nil {
		return nil, err
	}

	logs, err := t.TaskProvider.Logs(id, tail, start, end)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, logs)
}

func (t *TaskController) ListTasks(c *fireball.Context) (fireball.Response, error) {
	summaries, err := t.TaskProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)
}
