package job

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type JobRunner func(job models.Job) error

type Dispatcher struct {
	runner JobRunner
	store  JobStore
}

func NewDispatcher(runner JobRunner, store JobStore) *Dispatcher {
	return &Dispatcher{
		runner: runner,
		store:  store,
	}
}

func (d *Dispatcher) RunEvery(period time.Duration) *time.Ticker {
	ticker := time.NewTicker(period)
	go func() {
		for range ticker.C {
			d.Run()
		}
	}()

	return ticker
}

func (d *Dispatcher) Run() {
	// todo: get jobs from the store where job.Status == Pending
	var jobs []models.Job

	// todo: use a worker queue to limit the number of jobs
	// running at one time
	for _, job := range jobs {
		go d.runJob(job)
	}
}

func (d *Dispatcher) runJob(job models.Job) {
	if err := d.runner(job); err != nil {
		// todo: set JobStatus to Error
		// todo: set JobError to err
		log.Printf("[ERROR] [JobRunner] Failed to run job %s: %v", job.JobID, err)
	}
}
