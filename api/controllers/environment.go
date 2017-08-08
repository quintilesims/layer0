package controllers

import (
	restful "github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/entity"
	"github.com/zpatrick/forge/common/models"
)

type EnvironmentController struct {
	Provider entity.Provider
}

func NewEnvironmentController(p entity.Provider) *EnvironmentController {
	return &EnvironmentController{
		Provider: p,
	}
}

func (e *EnvironmentController) CreateEnvironment(request *restful.Request, response *restful.Response) {
	var req models.CreateEnvironmentRequest
	if err := request.ReadEntity(&req); err != nil {
		// todo: handle error
	}

	// todo: call req.Validate in provider layer
	environment := e.Provider.GetEnvironment("")
	if err := environment.Create(req); err != nil {
		// todo: handle error
	}

	environmentModel, err := environment.Read()
	if err != nil {
		// todo: handle error
	}

	response.WriteAsJson(environmentModel)
}

func (e *EnvironmentController) DeleteEnvironment(request *restful.Request, response *restful.Response) {
}

func (e *EnvironmentController) GetEnvironment(request *restful.Request, response *restful.Response) {
}

func (e *EnvironmentController) ListEnvironments(request *restful.Request, response *restful.Response) {
}
