package job

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

const (
	DISPATCHER_PERIOD = time.Second * 5
)

func RunWorkersAndDispatcher(numWorkers int, store Store, runner Runner) *time.Ticker {
	queue := make(chan models.Job)
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, store, queue, runner)
		worker.Start()
	}

	dispatcher := NewDispatcher(store, queue)
	ticker := time.NewTicker(DISPATCHER_PERIOD)
	go func() {
		for range ticker.C {
			log.Printf("[INFO] [JobDispatcher] Starting dispatcher")
			if err := dispatcher.Run(); err != nil {
				log.Printf("[ERROR] [JobDispatcher] Failed to dispatch: %v", err)
			}
		}
	}()

	return ticker
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
