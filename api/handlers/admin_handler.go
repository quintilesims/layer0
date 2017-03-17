package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"net/http"
	"strings"
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

	id := service.PathParameter("id", "identifier of the environment").
		DataType("string")

	service.Route(service.GET("/version").
		Filter(basicAuthenticate).
		To(this.GetVersion).
		Doc("Returns Current API version"))

	service.Route(service.GET("/health").
		To(this.GetHealth).
		Doc("Returns Health of API Server"))

	service.Route(service.PUT("/scale/{id}").
		Filter(basicAuthenticate).
		To(this.RunResourceManager).
		Reads("").
		Param(id).
		Doc("Run resource manager on an environment"))

	service.Route(service.GET("/config").
		To(this.GetConfig).
		Doc("Returns Configuration of the API Server").
		Writes(models.APIConfig{}))

	service.Route(service.POST("/sql").
		Filter(basicAuthenticate).
		To(this.UpdateSQL).
		Reads(models.SQLVersion{}).
		Doc("Configures sql settings"))

	return service
}

func (this *AdminHandler) GetVersion(request *restful.Request, response *restful.Response) {
	version := config.APIVersion()
	response.WriteAsJson(version)
}

func (this *AdminHandler) GetConfig(request *restful.Request, response *restful.Response) {
	publicSubnets := []string{}
	for _, subnet := range strings.Split(config.AWSPublicSubnets(), ",") {
		publicSubnets = append(publicSubnets, subnet)
	}

	privateSubnets := []string{}
	for _, subnet := range strings.Split(config.AWSPrivateSubnets(), ",") {
		privateSubnets = append(privateSubnets, subnet)
	}

	model := models.APIConfig{
		Prefix:         config.Prefix(),
		VPCID:          config.AWSVPCID(),
		PublicSubnets:  publicSubnets,
		PrivateSubnets: privateSubnets,
	}

	response.WriteAsJson(model)
}

func (this *AdminHandler) GetHealth(request *restful.Request, response *restful.Response) {
	// todo: does right sizer have health anymore?
	response.WriteAsJson("")
}

func (this *AdminHandler) RunResourceManager(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	info, err := this.AdminLogic.RunResourceManager(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(info)
}

func (this *AdminHandler) UpdateSQL(request *restful.Request, response *restful.Response) {
	if err := this.AdminLogic.UpdateSQL(); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
