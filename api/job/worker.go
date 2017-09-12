package job

import (
	"log"

	"github.com/quintilesims/layer0/common/models"
)

type Worker struct {
	ID     int
	Store  Store
	Queue  chan models.Job
	Runner Runner
}

func NewWorker(id int, store Store, queue chan models.Job, runner Runner) *Worker {
	return &Worker{
		ID:     id,
		Store:  store,
		Queue:  queue,
		Runner: runner,
	}
}

func (w *Worker) Start() func() {
	quit := make(chan bool)
	go func() {
		log.Printf("[DEBUG] [JobWorker %d]: Start signalled\n", w.ID)
		for {
			select {
			case job := <-w.Queue:
				ok, err := w.Store.AcquireJob(job.JobID)
				if err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Unexpected error when acquiring job %s: %v", w.ID, job.JobID, err)
					continue
				}

				if !ok {
					log.Printf("[DEBUG] [JobWorker %d]: Job %s is already acquired", w.ID, job.JobID)
					continue
				}

				log.Printf("[INFO] [JobWorker %d]: Starting job %s", w.ID, job.JobID)
				if err := w.Runner.Run(job); err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Failed to run job %s: %v", w.ID, job.JobID, err)
					w.Store.SetJobError(job.JobID, err)
					continue
				}

				log.Printf("[INFO] [JobWorker %d]: Finished job %s", w.ID, job.JobID)
                                w.Store.SetJobStatus(job.JobID, Completed)
			case <-quit:
				log.Printf("[DEBUG] [JobWorker %d]: Quit signalled", w.ID)
				return
			}
		}
	}()

	return func() { quit <- true }
}
