package controllers

import (
	"encoding/json"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

func createJob(tagStore tag.Store, jobStore job.Store, jobType job.JobType, req interface{}) (fireball.Response, error) {
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

	jobID, err := jobStore.Insert(jobType, requestStr)
	if err != nil {
		return nil, err
	}

	t := models.Tag{
		EntityID:   jobID,
		EntityType: "job",
		Key:        "name",
		Value:      jobID,
	}

	if err := tagStore.Insert(t); err != nil {
		return nil, err
	}

	resp := models.CreateJobResponse{
		JobID: jobID,
	}

	return fireball.NewJSONResponse(200, resp)
}
