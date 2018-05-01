package retry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	var calls int
	fn := func() (shouldRetry bool, err error) {
		calls++
		return calls < 5, nil
	}

	if err := Retry(fn); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, calls)
}

func TestRetryError(t *testing.T) {
	fn := func() (shouldRetry bool, err error) {
		return false, fmt.Errorf("some error")
	}

	if err := Retry(fn); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestRetryOptions(t *testing.T) {
	newOption := func() (Option, *int) {
		var calls int
		option := func() error {
			calls++
			return nil
		}

		return option, &calls
	}

	var calls int
	fn := func() (shouldRetry bool, err error) {
		calls++
		return calls < 5, nil
	}

	optionA, callsA := newOption()
	optionB, callsB := newOption()

	if err := Retry(fn, optionA, optionB); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, *callsA)
	assert.Equal(t, 5, *callsB)
}

func TestRetryOptionError(t *testing.T) {
	option := func() error {
		return fmt.Errorf("some error")
	}

	fn := func() (shouldRetry bool, err error) {
		return false, nil
	}

	if err := Retry(fn, option); err == nil {
		t.Fatal("Error was nil!")
	}
}
