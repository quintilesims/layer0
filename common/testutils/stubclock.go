package testutils

import (
	"time"
	"sync"
)

type StubClock struct {
	Time time.Time
	once sync.Once
}

func (s *StubClock) init(){
	if s.Time.IsZero(){
		s.Time = time.Now()
	}
}

func (s *StubClock) Now() time.Time {
	s.once.Do(s.init)

	s.Time = s.Time.Add(time.Millisecond * 20)
	return s.Time
}

func (s *StubClock) Sleep(d time.Duration) {
	s.once.Do(s.init)

	s.Time = s.Time.Add(d)
}

func (s *StubClock) Since(t time.Time) time.Duration {
	s.once.Do(s.init)

	return s.Now().Sub(t)
}
