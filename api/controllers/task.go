package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

type TaskController struct {
	TaskProvider provider.TaskProvider
	JobScheduler job.Scheduler
}

func NewTaskController(t provider.TaskProvider, j job.Scheduler) *TaskController {
	return &TaskController{
		TaskProvider: t,
		JobScheduler: j,
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

	return scheduleJob(t.JobScheduler, job.CreateTaskJob, req)
}

func (t *TaskController) DeleteTask(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	return scheduleJob(t.JobScheduler, job.DeleteTaskJob, id)
}

func (t *TaskController) GetTask(c *fireball.Context) (fireball.Response, error) {
	id := c.PathVariables["id"]
	model, err := t.TaskProvider.Read(id)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, model)
}

func (t *TaskController) ListTasks(c *fireball.Context) (fireball.Response, error) {
	summaries, err := t.TaskProvider.List()
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, summaries)

}