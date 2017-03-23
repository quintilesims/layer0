package job

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/types"
	"time"
)

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

	model, err := j.Logic.JobStore.SelectByID(j.jobID)
	if err != nil {
		return err
	}

	switch types.JobType(model.JobType) {
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
		return fmt.Errorf("Unknown job type '%v'!", model.JobType)
	}

	j.Context = NewJobContext(j.jobID, j.Logic, model.Request)
	return nil
}

func (j *JobRunner) Run() error {
	if err := j.MarkStatus(types.InProgress); err != nil {
		return err
	}

	for i, step := range j.Steps {
		log.Infof("Running step '%s'", step.Name)

		if err := j.runStep(step, j.Context); err != nil {
			log.Errorf("Error on step '%s': %v", step.Name, err)

			j.rollback(i)

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

func (j *JobRunner) rollback(from int) {
	log.Infof("Starting Rollback")

	for i := from; i >= 0; i-- {
		step := j.Steps[i]

		if step.Rollback == nil {
			continue
		}

		context, rollbackSteps, err := step.Rollback(j.Context)
		if err != nil {
			log.Errorf("Failed to get rollback steps for '%s': %v", step.Name, err)
			continue
		}

		for _, rollbackStep := range rollbackSteps {
			log.Infof("Running rollback step '%s'", rollbackStep.Name)

			if err := j.runStep(rollbackStep, context); err != nil {
				log.Errorf("Error during rollback step '%s': %v", rollbackStep.Name, err)
			}
		}
	}
}
