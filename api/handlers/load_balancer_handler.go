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

type LoadBalancerHandler struct {
	LoadBalancerLogic logic.LoadBalancerLogic
	JobLogic          logic.JobLogic
}

func NewLoadBalancerHandler(loadBalancerLogic logic.LoadBalancerLogic, jobLogic logic.JobLogic) *LoadBalancerHandler {
	return &LoadBalancerHandler{
		LoadBalancerLogic: loadBalancerLogic,
		JobLogic:          jobLogic,
	}
}

func (l *LoadBalancerHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/loadbalancer").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the load balancer").
		DataType("string")

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(l.ListLoadBalancers).
		Doc("List all LoadBalancers").
		Returns(200, "OK", []models.LoadBalancer{}))

	service.Route(service.GET("{id}").
		Filter(basicAuthenticate).
		To(l.GetLoadBalancer).
		Doc("Return a single LoadBalancer").
		Param(id).
		Writes(models.LoadBalancer{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(l.CreateLoadBalancer).
		Doc("Create a new LoadBalancer").
		Reads(models.CreateLoadBalancerRequest{}).
		Returns(http.StatusCreated, "Created", models.LoadBalancer{}).
		Writes(models.LoadBalancer{}))

	service.Route(service.DELETE("{id}").
		Filter(basicAuthenticate).
		To(l.DeleteLoadBalancer).
		Doc("Delete a LoadBalancer").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.PUT("{id}/ports").
		Filter(basicAuthenticate).
		To(l.UpdateLoadBalancerPorts).
		Reads(models.UpdateLoadBalancerPortsRequest{}).
		Param(id).
		Doc("Update load balancer ports").
		Writes(models.LoadBalancer{}))

	service.Route(service.PUT("{id}/healthcheck").
		Filter(basicAuthenticate).
		To(l.UpdateLoadBalancerHealthCheck).
		Reads(models.UpdateLoadBalancerHealthCheckRequest{}).
		Param(id).
		Doc("Update load balancer health check").
		Writes(models.LoadBalancer{}))

	return service
}

func (l *LoadBalancerHandler) ListLoadBalancers(request *restful.Request, response *restful.Response) {
	loadbalancers, err := l.LoadBalancerLogic.ListLoadBalancers()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(loadbalancers)
}

func (l *LoadBalancerHandler) GetLoadBalancer(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	loadbalancer, err := l.LoadBalancerLogic.GetLoadBalancer(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(loadbalancer)
}

func (l *LoadBalancerHandler) DeleteLoadBalancer(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	job, err := l.JobLogic.CreateJob(types.DeleteLoadBalancerJob, id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	WriteJobResponse(response, job.JobID)
}

func (l *LoadBalancerHandler) CreateLoadBalancer(request *restful.Request, response *restful.Response) {
	var req models.CreateLoadBalancerRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	loadBalancer, err := l.LoadBalancerLogic.CreateLoadBalancer(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(loadBalancer)
}

func (l *LoadBalancerHandler) UpdateLoadBalancerPorts(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.UpdateLoadBalancerPortsRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	loadBalancer, err := l.LoadBalancerLogic.UpdateLoadBalancerPorts(id, req.Ports)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(loadBalancer)
}

func (l *LoadBalancerHandler) UpdateLoadBalancerHealthCheck(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var req models.UpdateLoadBalancerHealthCheckRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	loadBalancer, err := l.LoadBalancerLogic.UpdateLoadBalancerHealthCheck(id, req.HealthCheck)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(loadBalancer)
}
