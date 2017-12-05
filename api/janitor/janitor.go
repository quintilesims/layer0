package janitor

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/api/lock"
)

type Janitor struct {
	Name   string
	lock   lock.Lock
	lockID string
	fn     func() error
}

func NewJanitor(name, lockID string, lock lock.Lock, fn func() error) *Janitor {
	return &Janitor{
		Name:   name,
		lock:   lock,
		lockID: lockID,
		fn:     fn,
	}
}

func (j *Janitor) Run() error {
	acquired, err := j.lock.Acquire(j.lockID)
	if err != nil {
		return err
	}

	if !acquired {
		log.Printf("[DEBUG] %s Janitor: Lock already acquired", j.Name)
		return nil
	}

	defer func() {
		if err := j.lock.Release(j.lockID); err != nil {
			log.Printf("[ERROR] %s Janitor: Failed to release lock: %v", j.Name, err)
		}
	}()

	log.Printf("[DEBUG] %s Janitor: Starting Run", j.Name)
	return j.fn()
}

func (j *Janitor) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			if err := j.Run(); err != nil {
				log.Printf("[ERROR] %s Janitor: %v", j.Name, err)
			}
		}
	}()

	return ticker
}
