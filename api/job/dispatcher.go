package job

import (
	"github.com/quintilesims/layer0/common/models"
)

// todo: better name
func RunWorkersAndDispatcher(numWorkers int, store Store, runner Runner) *Dispatcher {
	queue := make(chan models.Job)
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, queue, runner)
		worker.Start()
	}

	return NewDispatcher(store, queue)
}

type Dispatcher struct {
	store Store
	queue chan<- models.Job
}

func NewDispatcher(store Store, queue chan<- models.Job) *Dispatcher {
	return &Dispatcher{
		store: store,
		queue: queue,
	}
}

func (d *Dispatcher) Run() error {
	jobs, err := d.store.SelectAll()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if Status(job.Status) == Pending {
			// todo: a lot of time could pass while waiting for the queue to open up
			// the worker should attempt to acquire a lock before running the job
			d.queue <- *job
		}
	}

	return nil
}
