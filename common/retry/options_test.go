package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDelay(t *testing.T) {
	start := time.Now()
	if !WithDelay(time.Millisecond)() {
		t.Fatal("err")
	}

	assert.True(t, time.Since(start) > time.Millisecond)
}

func TestWithTimeout(t *testing.T) {
	option := WithTimeout(time.Millisecond)
	if !option() {
		t.Fatal("err")
	}

	time.Sleep(time.Millisecond)
	if !option() {
		t.Fatal("Error was nil!")
	}
}

func TestWithMaxAttempts(t *testing.T) {
	option := WithMaxAttempts(1)
	if !option() {
		t.Fatal("err")
	}

	if !option() {
		t.Fatal("Error was nil!")
	}
}
