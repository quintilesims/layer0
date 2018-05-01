package retry

import (
	"fmt"
	"time"
)

type Option func() error

func WithDelay(d time.Duration) Option {
	return func() error {
		time.Sleep(d)
		return nil
	}
}

func WithTimeout(d time.Duration) Option {
	start := time.Now()
	return func() error {
		if time.Since(start) > d {
			return fmt.Errorf("Timeout after %s", d.String())
		}

		return nil
	}
}

func WithMaxAttempts(max int) Option {
	var attempts int
	return func() error {
		attempts++
		if attempts > max {
			return fmt.Errorf("Maximum retry attempts reached")
		}

		return nil
	}
}
