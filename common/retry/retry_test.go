package retry

import (
	"fmt"
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

func TestRetryError(t *testing.T) {
	var err error
	fn := func() (shouldRetry bool) {
		err = fmt.Errorf("some error")
		return false
	}

	Retry(fn)

	if err == nil {
		t.Fatal("Error was nil!")
	}

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

func TestRetryOptionError(t *testing.T) {
	option := func() bool {
		return false
	}

	fn := func() (shouldRetry bool) {
		return false
	}

	Retry(fn, option)
}
