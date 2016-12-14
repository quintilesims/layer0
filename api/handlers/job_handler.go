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
		To(this.SelectAll).
		Doc("List all Jobs").
		Returns(200, "OK", []models.Job{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(this.SelectByID).
		Doc("Return a single Job").
		Param(id).
		Writes(models.Job{}))

	service.Route(service.DELETE("/{id}").
		Filter(basicAuthenticate).
		To(this.Delete).
		Doc("Stop and remove a job").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (this *JobHandler) SelectAll(request *restful.Request, response *restful.Response) {
	jobs, err := this.JobLogic.SelectAll()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(jobs)
}

func (this *JobHandler) SelectByID(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := this.JobLogic.SelectByID(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(job)
}

func (this *JobHandler) Delete(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := this.JobLogic.Delete(id); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(``)
}
