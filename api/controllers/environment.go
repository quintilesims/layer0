package controllers

import (
	restful "github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/api/scheduler"
	"github.com/quintilesims/layer0/common/job"
	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentController struct {
	Provider     entity.Provider
	JobScheduler scheduler.JobScheduler
}

func NewEnvironmentController(p entity.Provider, j scheduler.JobScheduler) *EnvironmentController {
	return &EnvironmentController{
		Provider:     p,
		JobScheduler: j,
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

	environmentModel, err := environment.Model()
	if err != nil {
		// todo: handle error
	}

	response.WriteAsJson(environmentModel)
}

func (e *EnvironmentController) DeleteEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		// todo: handle error
	}

	job := models.CreateJobRequest{
		JobType: job.DeleteEnvironmentJob,
		Request: id,
	}

	jobID, err := e.JobScheduler.ScheduleJob(job)
	if err != nil {
		// todo: handle error
	}

	WriteJobResponse(response, jobID)
}

func (e *EnvironmentController) GetEnvironment(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		// todo: handle error
	}

	environment := e.Provider.GetEnvironment(id)
	environmentModel, err := environment.Model()
	if err != nil {
		// todo: handle error
	}

	response.WriteAsJson(environmentModel)
}

func (e *EnvironmentController) ListEnvironments(request *restful.Request, response *restful.Response) {
	environmentIDs, err := e.Provider.ListEnvironmentIDs()
	if err != nil {
		// todo: handle error
	}

	environmentSummaries := make([]*models.EnvironmentSummary, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		environment := e.Provider.GetEnvironment(environmentID)
		environmentSummary, err := environment.Summary()
		if err != nil {
			// todo: handle error
		}

		environmentSummaries[i] = environmentSummary
	}

	response.WriteAsJson(environmentSummaries)
}
