package retry

import (
	"time"
)

type Option func() bool

func WithDelay(d time.Duration) Option {
	return func() bool {
		time.Sleep(d)
		return true
	}
}

func WithTimeout(d time.Duration) Option {
	start := time.Now()
	return func() bool {
		if time.Since(start) > d {
			return false
		}

		return true
	}
}

func WithMaxAttempts(max int) Option {
	var attempts int
	return func() bool {
		attempts++
		if attempts > max {
			return false
		}

		return true
	}
}
