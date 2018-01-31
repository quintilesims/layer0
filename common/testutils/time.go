package testutils

import (
	"testing"
	"time"
)

func WaitFor(t *testing.T, interval, timeout time.Duration, conditionSatisfied func() bool) {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(interval) {
		if conditionSatisfied() {
			return
		}
	}

	t.Fatalf("Timout reached after %v", timeout)
}
