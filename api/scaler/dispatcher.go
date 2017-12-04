package scaler

import (
	"log"
	"sync"
	"time"

	"github.com/quintilesims/layer0/api/provider"
)

const (
	SCALE_GRACE_PERIOD = time.Second * 1
)

var (
	timeMultiplier time.Duration = 1
)

type Dispatcher struct {
	environmentProvider provider.EnvironmentProvider
	scaler              Scaler
	schedule            map[string]*time.Timer
	lock                *sync.Mutex
}

func NewDispatcher(e provider.EnvironmentProvider, s Scaler) *Dispatcher {
	return &Dispatcher{
		environmentProvider: e,
		scaler:              s,
		schedule:            map[string]*time.Timer{},
		lock:                &sync.Mutex{},
	}
}

func (d *Dispatcher) ScheduleRun(environmentID string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if timer, ok := d.schedule[environmentID]; ok {
		log.Printf("[DEBUG] [ScalerDispatcher] Pushing back scheduled run for environment %s", environmentID)
		timer.Stop()
	}

	log.Printf("[DEBUG] [ScalerDispatcher] Scaling environment %s in %v", environmentID, SCALE_GRACE_PERIOD)
	d.schedule[environmentID] = time.AfterFunc(SCALE_GRACE_PERIOD*timeMultiplier, func() {
		// remove this run from the schedule
		d.lock.Lock()
		if _, ok := d.schedule[environmentID]; ok {
			delete(d.schedule, environmentID)
		}
		d.lock.Unlock()

		log.Printf("[DEBUG] [ScalerDispatcher] Scaling environment %s", environmentID)
		if err := d.scaler.Scale(environmentID); err != nil {
			log.Printf("[ERROR] [ScalerDispatcher] Failed to scale environment %s: %v", environmentID, err)
			return
		}

		log.Printf("[DEBUG] [ScalerDispatcher] Finished scaling environment %s", environmentID)
	})
}

func (d *Dispatcher) RunEvery(period time.Duration) *time.Ticker {
	ticker := time.NewTicker(period)
	go func() {
		for range ticker.C {
			d.RunAll()
		}
	}()

	return ticker
}

func (d *Dispatcher) RunAll() {
	log.Printf("[DEBUG] [ScalerDispatcher] Scaling all environments")

	environments, err := d.environmentProvider.List()
	if err != nil {
		log.Printf("[ERROR] [ScalerDispatcher] Failed to list environments: %v", err)
		return
	}

	for _, environment := range environments {
		d.ScheduleRun(environment.EnvironmentID)
	}
}
