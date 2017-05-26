package job

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"time"
)

const (
	JOB_LOAD_ATTEMPTS       = 10
	JOB_LOAD_SLEEP_INTERVAL = time.Second * 5
)

var timeMultiplier time.Duration = 1

type JobRunner struct {
	Logic   *logic.Logic
	Context *JobContext
	Steps   []Step
	jobID   string
}

func NewJobRunner(logic *logic.Logic, jobID string) *JobRunner {
	return &JobRunner{
		jobID: jobID,
		Logic: logic,
	}
}

func (j *JobRunner) MarkStatus(status types.JobStatus) error {
	return j.Logic.JobStore.UpdateJobStatus(j.jobID, status)
}

func (j *JobRunner) Load() error {
	log.Infof("Loading job '%s'", j.jobID)

	job, err := j.tryLoadJob()
	if err != nil {
		return err
	}

	switch types.JobType(job.JobType) {
	case types.DeleteEnvironmentJob:
		j.Steps = DeleteEnvironmentSteps
	case types.DeleteLoadBalancerJob:
		j.Steps = DeleteLoadBalancerSteps
	case types.DeleteServiceJob:
		j.Steps = DeleteServiceSteps
	case types.DeleteTaskJob:
		j.Steps = DeleteTaskSteps
	case types.CreateTaskJob:
		j.Steps = CreateTaskSteps
	default:
		return fmt.Errorf("Unknown job type '%v'!", job.JobType)
	}

	j.Context = NewJobContext(j.jobID, j.Logic, job.Request)
	return nil
}

func (j *JobRunner) tryLoadJob() (*models.Job, error) {
	var job *models.Job
	var err error

	for i := 0; i < JOB_LOAD_ATTEMPTS; i++ {
		job, err = j.Logic.JobStore.SelectByID(j.jobID)
		if err != nil {
			log.Warningf("Failed to load job %s (attempt %d/%d): %v", j.jobID, i+1, JOB_LOAD_ATTEMPTS, err)
			time.Sleep(JOB_LOAD_SLEEP_INTERVAL * timeMultiplier)
			continue
		}

		return job, nil
	}

	return nil, err
}

func (j *JobRunner) Run() error {
	if err := j.MarkStatus(types.InProgress); err != nil {
		return err
	}

	for _, step := range j.Steps {
		log.Infof("Running step '%s'", step.Name)

		if err := j.runStep(step, j.Context); err != nil {
			log.Errorf("Error on step '%s': %v", step.Name, err)

			if err := j.MarkStatus(types.Error); err != nil {
				log.Errorf("Failed to mark job status to Error: %v", err)
			}

			return fmt.Errorf("Error on step '%s': %v", step.Name, err)
		}
	}

	return j.MarkStatus(types.Completed)
}

func (j *JobRunner) runStep(step Step, context *JobContext) error {
	var err error
	quitc := make(chan bool)
	stepc := make(chan error)
	go func() { stepc <- step.Action(quitc, context) }()

	select {
	case err = <-stepc:
	case <-time.After(step.Timeout):
		close(quitc)
		<-stepc
		err = fmt.Errorf("Timeout reached after %v", step.Timeout)
	}

	return err
}
