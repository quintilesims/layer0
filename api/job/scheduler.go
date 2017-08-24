package job

import "github.com/quintilesims/layer0/common/models"

type Scheduler interface {
	Schedule(req models.ScheduleJobRequest) (string, error)
}

type SchedulerFunc func(models.ScheduleJobRequest) (string, error)

func (s SchedulerFunc) Schedule(req models.ScheduleJobRequest) (string, error) {
	return s(req)
}
