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

func (d *Daemon) Run() error {
	acquired, err := d.lock.Acquire(d.lockID)
	if err != nil {
		return err
	}

	if !acquired {
		log.Printf("[DEBUG] [%sDaemon]: Lock %s already acquired", d.Name, d.lockID)
		return nil
	}

	defer func() {
		if err := d.lock.Release(d.lockID); err != nil {
			log.Printf("[ERROR] [%sDaemon]: Failed to release lock %s: %v", d.Name, d.lockID, err)
		}
	}()

	log.Printf("[DEBUG] [%sDaemon]: Starting Run", d.Name)
	return d.fn()
}

func (d *Daemon) RunEvery(period time.Duration) *time.Ticker {
	ticker := time.NewTicker(period)
	go func() {
		for range ticker.C {
			if err := d.Run(); err != nil {
				log.Printf("[ERROR] [%sDaemon]: %v", d.Name, err)
			}
		}
	}()

	return ticker
}
