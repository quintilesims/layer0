package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/api/logic"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/types"
	"net/http"
)

type EnvironmentHandler struct {
	EnvironmentLogic logic.EnvironmentLogic
	JobLogic         logic.JobLogic
}

func NewEnvironmentHandler(environmentLogic logic.EnvironmentLogic, jobLogic logic.JobLogic) *EnvironmentHandler {
	return &EnvironmentHandler{
		EnvironmentLogic: environmentLogic,
		JobLogic:         jobLogic,
	}
}

func (this *EnvironmentHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/environment").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the environment").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.ListEnvironments).
		Doc("List all Environments").
		Returns(200, "OK", []models.Environment{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(this.GetEnvironment).
		Doc("Return a single Environment").
		Param(id).
		Writes(models.Environment{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(this.CreateEnvironment).
		Doc("Create a new Environment").
		Reads(models.CreateEnvironmentRequest{}).
		Returns(http.StatusCreated, "Created", models.Environment{}).
		Writes(models.Environment{}))

	service.Route(service.PUT("{id}").
		Filter(basicAuthenticate).
		To(this.UpdateEnvironment).
		Reads(models.UpdateEnvironmentRequest{}).
		Param(id).
		Doc("Update environment").
		Writes(models.Environment{}))

	service.Route(service.DELETE("{id}").
		Filter(basicAuthenticate).
		To(this.DeleteEnvironment).
		Doc("Delete an Environment").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (this *EnvironmentHandler) ListEnvironments(request *restful.Request, response *restful.Response) {
	environments, err := this.EnvironmentLogic.ListEnvironments()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environments)
}

func (this *EnvironmentHandler) GetEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	environment, err := this.EnvironmentLogic.GetEnvironment(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}

func (this *EnvironmentHandler) DeleteEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Paramter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := this.JobLogic.CreateJob(types.DeleteEnvironmentJob, id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	WriteJobResponse(response, job.JobID)
}

func (this *EnvironmentHandler) CreateEnvironment(request *restful.Request, response *restful.Response) {
	var req models.CreateEnvironmentRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	ok, err := this.EnvironmentLogic.CanCreateEnvironment(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	if !ok {
		err := fmt.Errorf("Environment with name '%s' already exists", req.EnvironmentName)
		BadRequest(response, errors.InvalidEnvironmentID, err)
		return
	}

	environment, err := this.EnvironmentLogic.CreateEnvironment(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}

func (this *EnvironmentHandler) UpdateEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.UpdateEnvironmentRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	environment, err := this.EnvironmentLogic.UpdateEnvironment(id, req.MinClusterCount)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}
