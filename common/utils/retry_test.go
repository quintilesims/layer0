package utils

import (
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetryTimeout(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	}

	Retry(time.Millisecond, time.Millisecond, fn)
	assert.WithinDuration(t, start.Add(time.Millisecond), time.Now(), time.Millisecond)
}

func TestRetryNoTimeout(t *testing.T) {
	fn := func() (bool, error) {
		return false, nil
	}

	if err := Retry(time.Millisecond*100, time.Millisecond*1, fn); err != nil {
		t.Fatal(err)
	}
}

func TestRetryWaitTimeLongerThanTimeout(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	}

	err := Retry(time.Millisecond*5, time.Millisecond*1, fn)
	assert.Equal(t, err, errors.New(errors.FailedRequestTimeout, nil))
	assert.WithinDuration(t, time.Now(), start.Add(time.Millisecond*5), time.Millisecond)
}

func TestRetryCount(t *testing.T) {
	retries := 0
	fn := func() (bool, error) {
		retries++
		return true, nil
	}

	Retry(time.Millisecond*5, time.Millisecond, fn)
	assert.Equal(t, retries, 6)
}

func TestRetryFNRanBeforeDelay(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		assert.WithinDuration(t, time.Now(), start, time.Millisecond*10)
		return true, nil
	}

	Retry(time.Millisecond*5, time.Millisecond, fn)
}
