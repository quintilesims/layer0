package job

import (
	"log"

	"github.com/quintilesims/layer0/api/lock"
	"github.com/quintilesims/layer0/common/models"
)

type Worker struct {
	ID     int
	Store  Store
	Queue  chan string
	Runner Runner
	Lock   lock.Lock
}

func NewWorker(id int, store Store, queue chan string, runner Runner, lock lock.Lock) *Worker {
	return &Worker{
		ID:     id,
		Store:  store,
		Queue:  queue,
		Runner: runner,
		Lock:   lock,
	}
}

func (w *Worker) Start() func() {
	quit := make(chan bool)
	go func() {
		log.Printf("[DEBUG] [Worker] [%d]: Start signalled\n", w.ID)
		for {
			select {
			case jobID := <-w.Queue:
				acquired, err := w.Lock.Acquire(jobID)
				if err != nil {
					log.Printf("[ERROR] [Worker] [%d]: Failed to acquire lock %s: %v", w.ID, jobID, err)
					continue
				}

				if !acquired {
					log.Printf("[DEBUG] [Worker [%d]: Job lock %s is already acquired", w.ID, jobID)
					continue
				}

				job, err := w.Store.SelectByID(jobID)
				if err != nil {
					log.Printf("[ERROR] [Worker] [%d]: Failed to select job %s: %v", w.ID, jobID, err)
					continue
				}

				log.Printf("[INFO] [Worker] [%d]: Starting job %s", w.ID, jobID)
				result, err := w.Runner.Run(*job)
				if err != nil {
					log.Printf("[ERROR] [Worker] [%d]: Failed to run job %s: %v", w.ID, jobID, err)
					w.Store.SetJobError(jobID, err)
					continue
				}

				if result != "" {
					if err := w.Store.SetJobResult(jobID, result); err != nil {
						log.Printf("[ERROR] [Worker] [%d]: Failed to set job result for job %s: %v", w.ID, jobID, err)
						w.Store.SetJobError(jobID, err)
						continue
					}
				}

				if err := w.Store.SetJobStatus(jobID, models.CompletedJobStatus); err != nil {
					log.Printf("[ERROR] [Worker] [%d]: Failed to set job status for job %s: %v", w.ID, jobID, err)
					continue
				}

				if err := w.Lock.Release(jobID); err != nil {
					log.Printf("[ERROR] [Worker] [%d]: Failed to release lock %s: %v", w.ID, jobID, err)
					continue
				}

				log.Printf("[INFO] [Worker] [%d]: Finished job %s", w.ID, jobID)
			case <-quit:
				log.Printf("[DEBUG] [Worker] [%d]: Quit signalled", w.ID)
				return
			}
		}
	}()

	return func() { quit <- true }
}
