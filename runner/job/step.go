package job

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/errors"
	"sync"
	"time"
)

type Action func(chan bool, *JobContext) error

type Step struct {
	Name    string
	Timeout time.Duration
	Action  Action
}

// Fold takes a slice of Actions and runs them async
func Fold(actions ...Action) Action {
	return func(quit chan bool, context *JobContext) error {
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
