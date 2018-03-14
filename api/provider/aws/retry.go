package aws

import (
	"time"

	"github.com/quintilesims/layer0/common/errors"
)

func retry(timeout, tick time.Duration, ch chan<- error, fn func() error) {
	var lastError error
	after := time.After(timeout)
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				lastError = err
				break
			}
			ch <- nil
			return
		case <-after:
			ch <- errors.New(errors.FailedRequestTimeout, lastError)
			return
		}
	}
}
