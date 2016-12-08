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
	Context JobContext
	Steps   []Step
	jobID   string
}

func NewJobRunner(logic *logic.Logic, jobID string) *JobRunner {
	return &JobRunner{
		jobID: jobID,
		Logic: logic,
	}
}

func (this *JobRunner) MarkStatus(status types.JobStatus) error {
	return this.Logic.JobData.UpdateJobStatus(this.jobID, status)
}

func (this *JobRunner) Load() error {
	log.Infof("Loading job '%s'", this.jobID)

	model, err := this.Logic.JobData.GetJob(this.jobID)
	if err != nil {
		return err
	}

	switch types.JobType(model.JobType) {
	case types.DeleteEnvironmentJob:
		this.Steps = DeleteEnvironmentSteps
	case types.DeleteLoadBalancerJob:
		this.Steps = DeleteLoadBalancerSteps
	case types.DeleteServiceJob:
		this.Steps = DeleteServiceSteps
	default:
		return fmt.Errorf("Unknown job type '%v'!", model.JobType)
	}

	this.Context = NewL0JobContext(this.jobID, this.Logic, model.Request)
	return nil
}

func (this *JobRunner) Run() error {
	if err := this.MarkStatus(types.InProgress); err != nil {
		return err
	}

	for i, step := range this.Steps {
		log.Infof("Running step '%s'", step.Name)

		if err := this.runStep(step, this.Context); err != nil {
			log.Errorf("Error on step '%s': %v", step.Name, err)

			this.rollback(i)

			if err := this.MarkStatus(types.Error); err != nil {
				log.Errorf("Failed to mark job status to Error: %v", err)
			}

			return fmt.Errorf("Error on step '%s': %v", step.Name, err)
		}
	}

	return this.MarkStatus(types.Completed)
}

func (this *JobRunner) runStep(step Step, context JobContext) error {
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

func (this *JobRunner) rollback(from int) {
	log.Infof("Starting Rollback")

	for i := from; i >= 0; i-- {
		step := this.Steps[i]

		if step.Rollback == nil {
			continue
		}

		context, rollbackSteps, err := step.Rollback(this.Context)
		if err != nil {
			log.Errorf("Failed to get rollback steps for '%s': %v", step.Name, err)
			continue
		}

		for _, rollbackStep := range rollbackSteps {
			log.Infof("Running rollback step '%s'", rollbackStep.Name)

			if err := this.runStep(rollbackStep, context); err != nil {
				log.Errorf("Error during rollback step '%s': %v", rollbackStep.Name, err)
			}
		}
	}
}
