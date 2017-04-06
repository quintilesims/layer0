package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
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

func (e *EnvironmentHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/environment").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the environment").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(e.ListEnvironments).
		Doc("List all Environments").
		Returns(200, "OK", []models.Environment{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(e.GetEnvironment).
		Doc("Return a single Environment").
		Param(id).
		Writes(models.Environment{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(e.CreateEnvironment).
		Doc("Create a new Environment").
		Reads(models.CreateEnvironmentRequest{}).
		Returns(http.StatusCreated, "Created", models.Environment{}).
		Writes(models.Environment{}))

	service.Route(service.PUT("{id}").
		Filter(basicAuthenticate).
		To(e.UpdateEnvironment).
		Reads(models.UpdateEnvironmentRequest{}).
		Param(id).
		Doc("Update environment").
		Writes(models.Environment{}))

	service.Route(service.DELETE("{id}").
		Filter(basicAuthenticate).
		To(e.DeleteEnvironment).
		Doc("Delete an Environment").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.POST("{id}/link").
		Filter(basicAuthenticate).
		To(e.CreateEnvironmentLink).
		Doc("Create an Environment Link").
		Reads(models.CreateEnvironmentLinkRequest{}).
		Param(id).
		Returns(http.StatusNoContent, "Created", nil))

	sourceID := service.PathParameter("source_id", "identifier of the source environment").
		DataType("string")

	destID := service.PathParameter("dest_id", "identifier of the destination environment").
		DataType("string")

	service.Route(service.DELETE("{source_id}/link/{dest_id}").
		Filter(basicAuthenticate).
		To(e.DeleteEnvironmentLink).
		Doc("Delete an Environment Link").
		Param(sourceID).
		Param(destID).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (e *EnvironmentHandler) ListEnvironments(request *restful.Request, response *restful.Response) {
	environments, err := e.EnvironmentLogic.ListEnvironments()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environments)
}

func (e *EnvironmentHandler) GetEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	environment, err := e.EnvironmentLogic.GetEnvironment(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}

func (e *EnvironmentHandler) DeleteEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Paramter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := e.JobLogic.CreateJob(types.DeleteEnvironmentJob, id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	WriteJobResponse(response, job.JobID)
}

func (e *EnvironmentHandler) CreateEnvironment(request *restful.Request, response *restful.Response) {
	var req models.CreateEnvironmentRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	ok, err := e.EnvironmentLogic.CanCreateEnvironment(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	if !ok {
		err := fmt.Errorf("Environment with name '%s' already exists", req.EnvironmentName)
		BadRequest(response, errors.InvalidEnvironmentID, err)
		return
	}

	environment, err := e.EnvironmentLogic.CreateEnvironment(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}

func (e *EnvironmentHandler) UpdateEnvironment(request *restful.Request, response *restful.Response) {
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

	environment, err := e.EnvironmentLogic.UpdateEnvironment(id, req.MinClusterCount)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(environment)
}

func (e *EnvironmentHandler) CreateEnvironmentLink(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.CreateEnvironmentLinkRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := e.EnvironmentLogic.CreateEnvironmentLink(id, req.EnvironmentID); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson("")
}

func (e *EnvironmentHandler) DeleteEnvironmentLink(request *restful.Request, response *restful.Response) {
	sourceID := request.PathParameter("source_id")
	if sourceID == "" {
		err := fmt.Errorf("Parameter 'source_id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	destID := request.PathParameter("dest_id")
	if destID == "" {
		err := fmt.Errorf("Parameter 'dest_id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := e.EnvironmentLogic.DeleteEnvironmentLink(sourceID, destID); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson("")
}
