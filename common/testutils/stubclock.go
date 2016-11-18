package testutils

import (
	"time"
)

type StubClock struct {
	InnerTime time.Time
}

func (this *StubClock) Now() time.Time {
	this.InnerTime = this.InnerTime.Add(time.Millisecond * 20)
	return this.InnerTime
}

func (this *StubClock) Sleep(s time.Duration) {
	this.InnerTime = this.InnerTime.Add(s)
}

func (this *StubClock) Since(t time.Time) time.Duration {
	return this.Now().Sub(t)
}
