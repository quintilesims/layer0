package job

import "github.com/quintilesims/layer0/common/models"

type Scheduler interface {
	Schedule(req models.ScheduleJobRequest) (string, error)
	Unschedule(jobID string) error
	List() ([]*models.Job, error)
	Read(jobID string) (*models.Job, error)
}
