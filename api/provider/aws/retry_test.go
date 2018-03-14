package aws

import (
	"fmt"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetry_Timeout(t *testing.T) {
	var ch = make(chan error, 1)
	lastKnownError := fmt.Errorf("Tick")

	fn := func() error {
		return lastKnownError
	}

	go retry(20*time.Millisecond, 5*time.Millisecond, ch, fn)
	assert.Equal(t, <-ch, errors.New(errors.FailedRequestTimeout, lastKnownError))
}

func TestRetry_NoTimeout(t *testing.T) {
	var ch chan error = make(chan error, 1)
	fn := func() error {
		return nil
	}

	go retry(20*time.Millisecond, 5*time.Millisecond, ch, fn)
	if err := <-ch; err != nil {
		t.Fatal(err)
	}
}

func TestRetry_GoRoutines(t *testing.T) {
	ch := make(chan error, 2)
	defer close(ch)

	fn := func() error {
		return fmt.Errorf("Error")
	}

	fn2 := func() error {
		return fmt.Errorf("Error 2")
	}

	go retry(20*time.Millisecond, 5*time.Millisecond, ch, fn)
	go retry(20*time.Millisecond, 5*time.Millisecond, ch, fn2)

	expectedErrors := []error{
		errors.New(errors.FailedRequestTimeout, fmt.Errorf("Error")),
		errors.New(errors.FailedRequestTimeout, fmt.Errorf("Error 2")),
	}

	for i := 0; i < 2; i++ {
		err := <-ch
		assert.Contains(t, expectedErrors, err)
	}
}
