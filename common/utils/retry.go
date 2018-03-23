package utils

import (
	"time"

	"github.com/quintilesims/layer0/common/errors"
)

func Retry(timeout, delay time.Duration, fn func() (bool, error)) error {
	ticker := time.NewTicker(delay)
	errc := make(chan error)
	defer func() {
		ticker.Stop()
		close(errc)
	}()

	go func() { errc <- retry(ticker, fn) }()

	select {
	case err := <-errc:
		return err
	case <-time.After(timeout):
		return errors.New(errors.FailedRequestTimeout, nil)
	}
}

func retry(ticker *time.Ticker, fn func() (bool, error)) error {
	for ; true; <-ticker.C {
		shouldContinue, err := fn()
		if err != nil {
			return err
		}

		if !shouldContinue {
			return err
		}
	}

	return nil
}
