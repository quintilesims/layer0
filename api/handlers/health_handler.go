package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
)

type HealthHandler struct {
	HealthLogic logic.HealthLogic
}

func NewHealthHandler(healthLogic logic.HealthLogic) *HealthHandler {
	return &HealthHandler{
		HealthLogic: healthLogic,
	}
}

func (this HealthHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/health").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/").
		To(this.GetHealth).
		Doc("Returns Health of API Server"))

	return service
}

func (this *HealthHandler) GetHealth(request *restful.Request, response *restful.Response) {
	response.WriteAsJson("")
}
