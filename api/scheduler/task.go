package scheduler

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/common/lock"
	"github.com/quintilesims/layer0/common/models"
)

type ECSTaskScheduler struct {
	lock lock.Lock
}

func NewECSTaskScheduler(l lock.Lock) *ECSTaskScheduler {
	return &ECSTaskScheduler{
		lock: l,
	}
}

func (s *ECSTaskScheduler) ScheduleTask(req models.CreateTaskRequest) (string, error) {
	// add entry into dynamodb, return the unique id of the entry
	return "", nil
}

func (s *ECSTaskScheduler) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			if err := s.run(); err != nil {
				log.Printf("[ERROR] [Task Scheduler] Encountered error during run: %v", err)
				continue
			}
		}
	}()

	return ticker
}

func (s *ECSTaskScheduler) run() error {
	if err := s.lock.Acquire(); err != nil {
		if lock.IsAcquiredError(err) {
			log.Printf("[INFO] [Task Scheduler] Lock is already acquired")
			return nil
		}

		return err
	}
	defer s.lock.Release()

	// for each entry in dynamodb, call ecs.StartTask(...)
	// if task starts successfully, remove it from the scheduler db
	return nil
}
