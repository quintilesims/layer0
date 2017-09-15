package job

import (
	"log"
)

type Worker struct {
	ID     int
	Store  Store
	Queue  chan string
	Runner Runner
}

func NewWorker(id int, store Store, queue chan string, runner Runner) *Worker {
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
			case jobID := <-w.Queue:
				ok, err := w.Store.AcquireJob(jobID)
				if err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Unexpected error when acquiring job %s: %v", w.ID, jobID, err)
					continue
				}

				if !ok {
					log.Printf("[DEBUG] [JobWorker %d]: Job %s is already acquired", w.ID, jobID)
					continue
				}

				job, err := w.Store.SelectByID(jobID)
				if err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Failed to select job %s: %v", w.ID, jobID, err)
					continue
				}

				log.Printf("[INFO] [JobWorker %d]: Starting job %s", w.ID, jobID)
				if err := w.Runner.Run(*job); err != nil {
					log.Printf("[ERROR] [JobWorker %d]: Failed to run job %s: %v", w.ID, jobID, err)
					w.Store.SetJobError(jobID, err)
					continue
				}

				log.Printf("[INFO] [JobWorker %d]: Finished job %s", w.ID, jobID)
				w.Store.SetJobStatus(jobID, Completed)
			case <-quit:
				log.Printf("[DEBUG] [JobWorker %d]: Quit signalled", w.ID)
				return
			}
		}
	}()

	return func() { quit <- true }
}
