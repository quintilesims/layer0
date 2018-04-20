package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDelay(t *testing.T) {
	start := time.Now()
	if !WithDelay(time.Millisecond)() {
		t.Fatal("Setup failed")
	}

	assert.True(t, time.Since(start) > time.Millisecond)
}

func TestWithTimeout(t *testing.T) {
	option := WithTimeout(time.Millisecond)

	time.Sleep(2 * time.Millisecond)
	if option() {
		t.Fatal("Setup failed")
	}
}

func TestWithMaxAttempts(t *testing.T) {
	option := WithMaxAttempts(1)
	if !option() {
		t.Fatal("Setup failed")
	}
}
