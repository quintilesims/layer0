package job

// Job Design:
// Each JobType has list of Steps that get executed in sequence.
// The JobRunner will continue to execute each step until an error occurs or it runs out of steps.
// If an error occurs, the JobRunner will traverse the list of steps from the current index back to the start,
// running the rollback function (if it is set) for each step along the way.
// Regardless of any actions the rollback(s) take, the job will always be marked in the Error state.
// Because of this, any retry logic should be performed in step.Action, not step.Rollback.

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/errors"
	"sync"
	"time"
)

type Action func(chan bool, JobContext) error
type Rollback func(JobContext) (JobContext, []Step, error)

type Step struct {
	Name     string
	Timeout  time.Duration
	Action   Action
	Rollback Rollback
}

func Fold(actions ...Action) Action {
	return func(quit chan bool, context JobContext) error {
		var wg sync.WaitGroup
		var errs []error

		wg.Add(len(actions))

		for _, fn := range actions {
			go func(fn Action) {
				defer wg.Done()
				if err := fn(quit, context); err != nil {
					errs = append(errs, err)
				}
			}(fn)
		}

		wg.Wait()
		return errors.MultiError(errs)
	}
}

func runAndRetry(quit chan bool, interval time.Duration, fn func() error) error {
	for {
		select {
		default:
			if err := fn(); err != nil {
				// todo: track errors and return them when quit is called
				log.Warning(err)
				time.Sleep(interval)
				continue
			}

			return nil
		case <-quit:
			return fmt.Errorf("Quit signalled")
		}
	}
}
