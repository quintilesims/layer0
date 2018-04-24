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

	if !option() {
		t.Fatal("Setup failed")
	}

	time.Sleep(2 * time.Millisecond)
}

func TestWithMaxAttempts(t *testing.T) {
	option := WithMaxAttempts(1)
	if !option() {
		t.Fatal("Setup failed")
	}

	if option() {
		t.Fatal("Max attempt reached")
	}
}
