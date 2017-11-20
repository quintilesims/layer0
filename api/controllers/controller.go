package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

func createJob(store job.Store, jobType models.JobType, req interface{}) (fireball.Response, error) {
	var requestStr string
	switch v := req.(type) {
	case string:
		requestStr = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		requestStr = string(bytes)
	}

	jobID, err := store.Insert(jobType, requestStr)
	if err != nil {
		return nil, err
	}

	resp := models.CreateJobResponse{
		JobID: jobID,
	}

	return fireball.NewJSONResponse(200, resp)
}
