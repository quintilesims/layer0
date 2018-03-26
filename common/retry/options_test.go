package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDelay(t *testing.T) {
	start := time.Now()
	if err := WithDelay(time.Millisecond)(); err != nil {
		t.Fatal(err)
	}

	assert.True(t, time.Since(start) > time.Millisecond)
}

func TestWithTimeout(t *testing.T) {
	option := WithTimeout(time.Millisecond)
	if err := option(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond)
	if err := option(); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestWithMaxAttempts(t *testing.T) {
	option := WithMaxAttempts(1)
	if err := option(); err != nil {
		t.Fatal(err)
	}

	if err := option(); err == nil {
		t.Fatal("Error was nil!")
	}
}
