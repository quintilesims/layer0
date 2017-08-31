package controllers

import (
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
)

func scheduleJob(scheduler job.Scheduler, jobType job.JobType, req interface{}) (fireball.Response, error) {
	job := models.ScheduleJobRequest{
		JobType: jobType.String(),
		Request: req,
	}

	jobID, err := scheduler.Schedule(job)
	if err != nil {
		return nil, err
	}

	resp := models.ScheduleJobResponse{
		JobID: jobID,
	}

	return fireball.NewJSONResponse(200, resp)
}
