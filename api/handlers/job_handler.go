package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"net/http"
)

type JobHandler struct {
	JobLogic logic.JobLogic
}

func NewJobHandler(jobLogic logic.JobLogic) *JobHandler {
	return &JobHandler{
		JobLogic: jobLogic,
	}
}

func (this *JobHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/job").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the job").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.ListJobs).
		Doc("List all Jobs").
		Returns(200, "OK", []models.Job{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(this.GetJob).
		Doc("Return a single Job").
		Param(id).
		Writes(models.Job{}))

	service.Route(service.DELETE("/{id}").
		Filter(basicAuthenticate).
		To(this.DeleteJob).
		Doc("Stop and remove a job").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (this *JobHandler) ListJobs(request *restful.Request, response *restful.Response) {
	jobs, err := this.JobLogic.ListJobs()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(jobs)
}

func (this *JobHandler) GetJob(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := this.JobLogic.GetJob(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(job)
}

func (this *JobHandler) DeleteJob(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := this.JobLogic.DeleteJob(id); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(``)
}
