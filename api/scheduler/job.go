package scheduler

import "github.com/quintilesims/layer0/common/models"

type ScheduleJobFunc func(req models.CreateJobRequest) (string, error)

func (s ScheduleJobFunc) ScheduleJob(req models.CreateJobRequest) (string, error) {
	return s(req)
}
