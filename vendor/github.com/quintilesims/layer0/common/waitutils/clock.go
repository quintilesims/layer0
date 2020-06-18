package waitutils

import (
	"time"
)

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
	Sleep(time.Duration)
}

type RealClock struct {
}

func (RealClock) Now() time.Time {
	return time.Now()
}

func (RealClock) Sleep(s time.Duration) {
	time.Sleep(s)
}

func (RealClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}
