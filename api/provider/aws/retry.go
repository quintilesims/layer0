package aws

import (
	"time"

	"github.com/quintilesims/layer0/common/errors"
)

func retry(timeout, tick time.Duration, fn func() error) error {
	after := time.After(timeout)

	for {
		select {
		case <-time.Tick(tick):
			if err := fn(); err == nil {
				return nil
			}
		case <-after:
			return errors.New(errors.FailedRequestTimeout, nil)
		}
	}
}
