package retry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	var calls int
	fn := func() (shouldRetry bool) {
		calls++
		return calls < 5
	}

	Retry(fn)
	assert.Equal(t, 5, calls)
}

func TestRetryCalled(t *testing.T) {
	var called bool
	fn := func() (shouldRetry bool) {
		called = false
		return called
	}

	Retry(fn)

	assert.Equal(t, false, called)
}

func TestRetryOptions(t *testing.T) {
	newOption := func() (Option, *int) {
		var calls int
		option := func() bool {
			calls++
			return true
		}

		return option, &calls
	}

	var calls int
	fn := func() (shouldRetry bool) {
		calls++
		return calls < 5
	}

	optionA, callsA := newOption()
	optionB, callsB := newOption()

	Retry(fn, optionA, optionB)

	assert.Equal(t, 5, *callsA)
	assert.Equal(t, 5, *callsB)
}

func TestRetryOptionCalled(t *testing.T) {
	var called bool

	option := func() bool {
		called = true
		return false
	}

	var calls int
	fn := func() (shouldRetry bool) {
		calls++
		return true
	}

	Retry(fn, option)

	assert.Equal(t, 0, calls)
	assert.Equal(t, true, called)
}
