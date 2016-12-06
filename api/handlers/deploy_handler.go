package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"net/http"
)

type DeployHandler struct {
	DeployLogic logic.DeployLogic
}

func NewDeployHandler(deployLogic logic.DeployLogic) *DeployHandler {
	return &DeployHandler{
		DeployLogic: deployLogic,
	}
}

func (this *DeployHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/deploy").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the deploy").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.ListDeploys).
		Doc("List all Deploys").
		Returns(200, "OK", []models.Deploy{}).
		Returns(400, "OK", []models.Deploy{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(this.GetDeploy).
		Doc("Return a single Deploy").
		Param(id).
		Writes(models.Deploy{}))

	service.Route(service.DELETE("/{id}").
		Filter(basicAuthenticate).
		To(this.DeleteDeploy).
		Doc("Delete a deploy").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.POST("/create").
		Filter(basicAuthenticate).
		To(this.CreateDeploy).
		Doc("Create a new Deploy").
		Returns(http.StatusCreated, "Created", models.Deploy{}).
		Reads(models.CreateDeployRequest{}))

	return service
}

func (this *DeployHandler) ListDeploys(request *restful.Request, response *restful.Response) {
	deploys, err := this.DeployLogic.ListDeploys()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(deploys)
}

func (this *DeployHandler) GetDeploy(request *restful.Request, response *restful.Response) {
	deployID := request.PathParameter("id")
	if deployID == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	deploy, err := this.DeployLogic.GetDeploy(deployID)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(deploy)
}

func (this *DeployHandler) DeleteDeploy(request *restful.Request, response *restful.Response) {
	deployID := request.PathParameter("id")
	if deployID == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := this.DeployLogic.DeleteDeploy(deployID); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(``)
}

func (this *DeployHandler) CreateDeploy(request *restful.Request, response *restful.Response) {
	var req models.CreateDeployRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	deploy, err := this.DeployLogic.CreateDeploy(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(deploy)
}
