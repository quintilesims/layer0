package job

import (
	"log"

	"github.com/quintilesims/layer0/common/models"
)

type Runner func(job models.Job) error

type Worker struct {
	ID     int
	Queue  chan models.Job
	Runner Runner
}

func NewWorker(id int, queue chan models.Job, runner Runner) *Worker {
	return &Worker{
		ID:     id,
		Queue:  queue,
		Runner: runner,
	}
}

func (w *Worker) Start() func() {
	quit := make(chan bool)
	go func() {
		log.Printf("[DEBUG] [JobWorker %d]: start signalled", w.ID)
		for {
			select {
			case job := <-w.Queue:
				log.Printf("[INFO] [JobWorker %d]: starting job %s\n", w.ID, job.JobID)
				if err := w.Runner(job); err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Failed to run job %s: %v\n", w.ID, job.JobID, err)
				}
			case <-quit:
				log.Printf("[DEBUG] [JobWorker %d]: quit signalled\n", w.ID)
				return
			}
		}
	}()

	return func() { quit <- true }
}
