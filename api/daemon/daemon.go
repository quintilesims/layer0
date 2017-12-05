package daemon

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/api/lock"
)

type Daemon struct {
	Name   string
	lock   lock.Lock
	lockID string
	fn     func() error
}

func NewDaemon(name, lockID string, lock lock.Lock, fn func() error) *Daemon {
	return &Daemon{
		Name:   name,
		lock:   lock,
		lockID: lockID,
		fn:     fn,
	}
}

func (j *Daemon) Run() error {
	acquired, err := j.lock.Acquire(j.lockID)
	if err != nil {
		return err
	}

	if !acquired {
		log.Printf("[DEBUG] [%sDaemon]: Lock already acquired", j.Name)
		return nil
	}

	defer func() {
		if err := j.lock.Release(j.lockID); err != nil {
			log.Printf("[ERROR] [%sDaemon]: Failed to release lock: %v", j.Name, err)
		}
	}()

	log.Printf("[DEBUG] [%sDaemon]: Starting Run", j.Name)
	return j.fn()
}

func (j *Daemon) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			if err := j.Run(); err != nil {
				log.Printf("[ERROR] [%sDaemon]: %v", j.Name, err)
			}
		}
	}()

	return ticker
}
