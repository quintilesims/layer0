package controllers

import (
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

func createJob(store job.Store, jobType job.JobType, request interface{}) (fireball.Response, error) {
	req := models.CreateJobRequest{
		JobType: jobType.String(),
		Request: request,
	}

	jobID, err := store.Insert(req)
	if err != nil {
		return nil, err
	}

	resp := models.CreateJobResponse{
		JobID: jobID,
	}

	return fireball.NewJSONResponse(200, resp)
}
