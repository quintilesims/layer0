package logic

import (
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/types"
	"github.com/quintilesims/layer0/common/waitutils"
)

var jobLogger = logutils.NewStackTraceLogger("Job Janitor")

type JobJanitor struct {
	jobLogic JobLogic
	Clock    waitutils.Clock
}

func NewJobJanitor(jobLogic JobLogic) *JobJanitor {
	return &JobJanitor{
		jobLogic: jobLogic,
		Clock:    waitutils.RealClock{},
	}
}

func (this *JobJanitor) Run() {
	go func() {
		for {
			jobLogger.Info("Starting cleanup")
			this.pulse()
			jobLogger.Infof("Finished cleanup")
			this.Clock.Sleep(JANITOR_SLEEP_DURATION)
		}
	}()
}

func (this *JobJanitor) pulse() error {
	jobs, err := this.jobLogic.ListJobs()
	if err != nil {
		jobLogger.Errorf("Failed to list jobs: %v", err)
		return err
	}

	errs := []error{}
	for _, job := range jobs {
		timeSinceCreated := this.Clock.Since(job.TimeCreated)
		if job.JobStatus != int64(types.InProgress) {
			if timeSinceCreated > JOB_LIFETIME || job.JobStatus == int64(types.Completed) {
				jobLogger.Infof("Deleting job '%s'", job.JobID)

				if err := this.jobLogic.Delete(job.JobID); err != nil {
					jobLogger.Errorf("Failed to delete job '%s': %v", job.JobID, err)
					errs = append(errs, err)
				} else {
					jobLogger.Infof("Finished deleting job '%s'", job.JobID)
				}
			}
		}
	}

	return errors.MultiError(errs)
}
