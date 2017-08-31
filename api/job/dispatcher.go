package job

import (
	"context"
	"log"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type JobRunner func(c context.Context, store JobStore, job *models.Job) error

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
	// get jobs from the store where job.Status == Pending
	var jobs []*models.Job

	for _, job := range jobs {
		// use a semver so we don't run > max jobs
		// attempt to acquire a lock on the job

		go func() {
			// set job.Status = InProgress

			if err := d.runJob(job); err != nil {
				log.Printf("[ERROR] [Job Dispatcher] %v", err)
				// set job.Status = Error
				// set job.Error = Error
			}
		}()
	}
}

func (d *Dispatcher) runJob(job *models.Job) error {
	// todo: use config timeout
	timeout := time.Minute * 1

	c, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	d.runner(c, nil, job)
	// todo: how to determine errors?
	// do we always retry until timeout and have the runner log errors along the way?
	return nil
}
