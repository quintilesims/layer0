package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
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
}

func NewTaskController(t provider.TaskProvider, j job.Store) *TaskController {
	return &TaskController{
		TaskProvider: t,
		JobStore:     j,
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

	return createJob(t.JobStore, job.CreateTaskJob, req)
}

func (t *TaskController) DeleteTask(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return createJob(t.JobStore, job.DeleteTaskJob, id)
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
	query := c.Request.URL.Query()

	var tail int
	if v := query.Get("tail"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.Newf(errors.InvalidRequest, "Tail must be an integer")
		}

		tail = t
	}

	parseTime := func(v string) (time.Time, error) {
		if v == "" {
			return time.Time{}, nil
		}

		return time.Parse(TIME_LAYOUT, v)
	}

	start, err := parseTime(query.Get("start"))
	if err != nil {
		return nil, errors.Newf(errors.InvalidRequest, "Invalid time: start must be in format YYYY-MM-DD HH:MM")
	}

	end, err := parseTime(query.Get("end"))
	if err != nil {
		return nil, errors.Newf(errors.InvalidRequest, "Invalid time: end must be in format YYYY-MM-DD HH:MM")
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
