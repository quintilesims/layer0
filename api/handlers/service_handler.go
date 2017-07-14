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

type ServiceHandler struct {
	ServiceLogic logic.ServiceLogic
	JobLogic     logic.JobLogic
}

func NewServiceHandler(serviceLogic logic.ServiceLogic, jobLogic logic.JobLogic) *ServiceHandler {
	return &ServiceHandler{
		ServiceLogic: serviceLogic,
		JobLogic:     jobLogic,
	}
}

func (this *ServiceHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/service").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the service").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.ListServices).
		Doc("List all services").
		Returns(200, "OK", []models.Service{}))

	service.Route(service.GET("/{id}").
		Filter(basicAuthenticate).
		To(this.GetService).
		Doc("Return a service").
		Param(id).
		Writes(models.Service{}))

	service.Route(service.DELETE("/{id}").
		Filter(basicAuthenticate).
		To(this.DeleteService).
		Doc("Stop and remove a service").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(this.CreateService).
		Doc("Create a service").
		Reads(models.CreateServiceRequest{}).
		Returns(http.StatusCreated, "Created", models.Service{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Service{}))

	service.Route(service.PUT("/{id}/scale").
		Filter(basicAuthenticate).
		To(this.ScaleService).
		Doc("Scale a service").
		Reads(models.ScaleServiceRequest{}).
		Param(id).
		Returns(http.StatusAccepted, "Scaling", models.Service{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Service{}))

	service.Route(service.PUT("/{id}/deploy").
		Filter(basicAuthenticate).
		To(this.UpdateService).
		Doc("Run a new deploy on a service").
		Reads(models.UpdateServiceRequest{}).
		Param(id).
		Returns(http.StatusAccepted, "Scaling", models.Service{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Service{}))

	service.Route(service.GET("/{id}/logs").
		Filter(basicAuthenticate).
		To(this.GetServiceLogs).
		Doc("Return recent service logs").
		Param(service.PathParameter("id", "identifier of the service").DataType("string")).
		Param(service.QueryParameter("tail", "number of lines from the end to return").DataType("string")).
		Param(service.QueryParameter("start", "The start of the time range to fetch logs (format MM/DD HH:MM)").DataType("string")).
		Param(service.QueryParameter("end", "The end of the time range to fetch logs (format MM/DD HH:MM)").DataType("string")).
		Writes([]models.LogFile{}))

	return service
}

func (this *ServiceHandler) ListServices(request *restful.Request, response *restful.Response) {
	services, err := this.ServiceLogic.ListServices()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(services)
}

func (this *ServiceHandler) DeleteService(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required.")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := this.JobLogic.CreateJob(types.DeleteServiceJob, id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	WriteJobResponse(response, job.JobID)
}

func (this *ServiceHandler) GetService(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	service, err := this.ServiceLogic.GetService(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(service)
}

func (this *ServiceHandler) CreateService(request *restful.Request, response *restful.Response) {
	var req models.CreateServiceRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	service, err := this.ServiceLogic.CreateService(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(service)
}

func (this *ServiceHandler) ScaleService(request *restful.Request, response *restful.Response) {
	serviceID := request.PathParameter("id")
	if serviceID == "" {
		err := fmt.Errorf("Parameter 'id' is required.")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.ScaleServiceRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	service, err := this.ServiceLogic.ScaleService(serviceID, int(req.DesiredCount))
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(service)
}

func (this *ServiceHandler) UpdateService(request *restful.Request, response *restful.Response) {
	serviceID := request.PathParameter("id")
	if serviceID == "" {
		err := fmt.Errorf("Parameter 'id' is required.")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.UpdateServiceRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	service, err := this.ServiceLogic.UpdateService(serviceID, req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(service)
}

func (this *ServiceHandler) GetServiceLogs(request *restful.Request, response *restful.Response) {
	serviceID := request.PathParameter("id")
	if serviceID == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.InvalidServiceID, err)
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

	logs, err := this.ServiceLogic.GetServiceLogs(serviceID, request.QueryParameter("start"), request.QueryParameter("end"), tail)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(logs)
}
