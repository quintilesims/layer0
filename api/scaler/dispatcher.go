package scaler

import (
	"log"
	"time"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/models"
)

type Dispatcher struct {
	jobStore    job.Store
	gracePeriod time.Duration
	scheduleOps chan func(map[string]*time.Timer)
}

func NewDispatcher(jobStore job.Store, gracePeriod time.Duration) *Dispatcher {
	d := &Dispatcher{
		jobStore:    jobStore,
		gracePeriod: gracePeriod,
		scheduleOps: make(chan func(map[string]*time.Timer)),
	}

	go d.loop()
	return d
}

func (d *Dispatcher) loop() {
	schedule := map[string]*time.Timer{}
	for op := range d.scheduleOps {
		op(schedule)
	}
}

func (d *Dispatcher) Dispatch(environmentID string) {
	d.scheduleOps <- func(schedule map[string]*time.Timer) {
		timer, ok := schedule[environmentID]
		if ok {
			log.Printf("[DEBUG] [ScalerDispatcher] Scaling environment %s in %v", environmentID, d.gracePeriod)
			timer.Reset(d.gracePeriod)
			return
		}

		schedule[environmentID] = time.AfterFunc(d.gracePeriod, func() {
			d.scheduleOps <- func(schedule map[string]*time.Timer) {
				delete(schedule, environmentID)
			}

			log.Printf("[DEBUG] [ScalerDispatcher] Creating scale job for environment %s", environmentID)
			if _, err := d.jobStore.Insert(models.ScaleEnvironmentJob, environmentID); err != nil {
				log.Printf("[ERROR] [ScalerDispatcer] Failed to create scale job for environment %s: %v", environmentID, err)
			}
		})
	}
}
