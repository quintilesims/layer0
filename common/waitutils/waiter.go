package waitutils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"time"
)

// If timeout is specified, waiter will honor it
// If retries == 0, waiter will wait forever
type Waiter struct {
	Name    string
	Timeout time.Duration
	Retries int
	Delay   time.Duration
	Check   func() (bool, error)
	Clock   Clock
}

func (this Waiter) Wait() error {
	start := time.Now()

	shouldContinue := func(i int) bool {
		if this.Timeout != 0 {
			if this.Clock.Since(start) > this.Timeout {
				return false
			}
		}

		if this.Retries == 0 {
			return true
		}

		return i < this.Retries
	}

	for i := 0; shouldContinue(i); i++ {
		if ok, err := this.Check(); err != nil {
			return err
		} else if ok {
			return nil
		}

		retryStr := fmt.Sprintf("%v", this.Retries)
		if this.Retries < 0 {
			retryStr = "<infinite>"
		}

		log.Debugf("Wait %s iteration %d of %s", this.Name, i+1, retryStr)
		this.Clock.Sleep(this.Delay)
	}

	return fmt.Errorf("Wait for `%s` timeout", this.Name)
}
