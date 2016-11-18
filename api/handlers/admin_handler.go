package handlers

import (
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/api/logic"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"net/http"
)

type AdminHandler struct {
	AdminLogic logic.AdminLogic
}

func NewAdminHandler(adminLogic logic.AdminLogic) *AdminHandler {
	return &AdminHandler{
		AdminLogic: adminLogic,
	}
}

func (this AdminHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/admin").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/sql").
		Filter(basicAuthenticate).
		To(this.GetSQL).
		Doc("Returns Current SQL status").
		Writes(models.SQLVersion{}))

	service.Route(service.GET("/version").
		Filter(basicAuthenticate).
		To(this.GetVersion).
		Doc("Returns Current API version"))

	service.Route(service.GET("/health").
		To(this.GetHealth).
		Doc("Returns Health of API Server"))

	service.Route(service.POST("/sql").
		Filter(basicAuthenticate).
		To(this.UpdateSQL).
		Reads(models.SQLVersion{}).
		Doc("Configures sql settings"))

	return service
}

func (this *AdminHandler) GetSQL(request *restful.Request, response *restful.Response) {
	version, err := this.AdminLogic.GetSQLStatus()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(version)
}

func (this *AdminHandler) GetVersion(request *restful.Request, response *restful.Response) {
	version := config.APIVersion()
	response.WriteAsJson(version)
}

func (this *AdminHandler) GetHealth(request *restful.Request, response *restful.Response) {
	message, err := this.AdminLogic.GetHealth()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(message)
}

func (this *AdminHandler) UpdateSQL(request *restful.Request, response *restful.Response) {
	if err := this.AdminLogic.UpdateSQL(); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
