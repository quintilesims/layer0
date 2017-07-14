package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type TaskHandler struct {
	TaskLogic logic.TaskLogic
	JobLogic  logic.JobLogic
}

func NewTaskHandler(taskLogic logic.TaskLogic, jobLogic logic.JobLogic) *TaskHandler {
	return &TaskHandler{
		TaskLogic: taskLogic,
		JobLogic:  jobLogic,
	}
}

func (this *TaskHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/task").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the task").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.ListTasks).
		Doc("List all tasks").
		Returns(200, "OK", []models.Task{}))

	service.Route(service.GET("/{id}").
		Filter(basicAuthenticate).
		To(this.GetTask).
		Doc("Return a task").
		Param(id).
		Writes(models.Task{}))

	service.Route(service.DELETE("/{id}").
		Filter(basicAuthenticate).
		To(this.DeleteTask).
		Doc("Stop and remove a task").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(this.CreateTask).
		Doc("Create a task").
		Reads(models.CreateTaskRequest{}).
		Returns(http.StatusCreated, "Created", models.Task{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Task{}))

	service.Route(service.GET("/{id}/logs").
		Filter(basicAuthenticate).
		To(this.GetTaskLogs).
		Doc("Return recent task logs").
		Param(service.PathParameter("id", "identifier of the task").DataType("string")).
		Param(service.QueryParameter("tail", "number of lines from the end to return").DataType("string")).
		Writes([]models.LogFile{}))

	return service
}

func (this *TaskHandler) ListTasks(request *restful.Request, response *restful.Response) {
	tasks, err := this.TaskLogic.ListTasks()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(tasks)
}

func (this *TaskHandler) DeleteTask(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := this.TaskLogic.DeleteTask(id); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(``)
}

func (this *TaskHandler) GetTask(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	task, err := this.TaskLogic.GetTask(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(task)
}

func (this *TaskHandler) CreateTask(request *restful.Request, response *restful.Response) {
	var req models.CreateTaskRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	job, err := this.JobLogic.CreateJob(types.CreateTaskJob, req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	WriteJobResponse(response, job.JobID)
}

func (this *TaskHandler) GetTaskLogs(request *restful.Request, response *restful.Response) {
	taskID := request.PathParameter("id")
	if taskID == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.InvalidTaskID, err)
		return
	}

	var tail int
	if param := request.QueryParameter("tail"); param != "" {
		t, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			BadRequest(response, errors.InvalidJSON, err)
			return
		}

		tail = int(t)
	}

	logs, err := this.TaskLogic.GetTaskLogs(taskID, tail)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(logs)
}
