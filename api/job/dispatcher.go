package job

import (
	"log"
	"time"
)

const (
	// dynamodb allows a burst-read/write every 5 minutes
	// matching that value here to try and avoid hitting that limit
	DISPATCHER_PERIOD = time.Minute * 5
)

func RunWorkersAndDispatcher(numWorkers int, store Store, runner Runner) *time.Ticker {
	queue := make(chan string)
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

	store.SetInsertHook(func(jobID string) {
		go func() { queue <- jobID }()
	})

	return ticker
}

type Dispatcher struct {
	store Store
	queue chan<- string
}

func NewDispatcher(store Store, queue chan<- string) *Dispatcher {
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
			d.queue <- job.JobID
		}
	}

	return nil
}
