package scheduler

import (
	"github.com/quintilesims/layer0/common/models"
)

type JobScheduler interface {
	ScheduleJob(req models.CreateJobRequest) (string, error)
}

type TaskScheduler interface {
	ScheduleTask(req models.CreateTaskRequest) (string, error)
}

type ServiceScheduler interface {
	ScheduleService(req models.CreateServiceRequest) (string, error)
}
