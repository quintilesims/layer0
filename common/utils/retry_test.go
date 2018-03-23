package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetryTimeout(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		time.Sleep((time.Duration(rand.Intn(2))) * time.Millisecond)
		return true, nil
	}

	if err := Retry(time.Millisecond*1, time.Millisecond*1, fn); err != nil {
		assert.WithinDuration(t, start.Add(time.Millisecond*1), time.Now(), time.Millisecond)
	}
}

func TestRetryNoTimeout(t *testing.T) {
	retries := 0
	fn := func() (bool, error) {
		retries++
		time.Sleep(1 * time.Millisecond)

		if retries == 5 {
			return false, nil
		}

		return true, nil
	}

	if err := Retry(time.Millisecond*100, time.Millisecond*1, fn); err != nil {
		t.Fatal(err)
	}
}

func TestRetryFuncTimeRandom(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		time.Sleep((time.Duration(rand.Intn(2))) * time.Millisecond)
		return true, nil
	}

	if err := Retry(time.Millisecond*5, time.Millisecond*1, fn); err != nil {
		assert.WithinDuration(t, time.Now(), start.Add(time.Millisecond*5), time.Millisecond)
	}
}

func TestRetryWaitTimeLongerThanTimeout(t *testing.T) {
	start := time.Now()
	retries := 0
	fn := func() (bool, error) {
		retries++
		time.Sleep(10 * time.Millisecond)
		return true, nil
	}

	if err := Retry(time.Millisecond*5, time.Millisecond*1, fn); err != nil {
		assert.Equal(t, retries, 1)
		assert.Equal(t, err, errors.New(errors.FailedRequestTimeout, nil))
		assert.WithinDuration(t, time.Now(), start.Add(time.Millisecond*5), time.Millisecond)
	}
}

func TestRetryCount(t *testing.T) {
	retries := 0
	fn := func() (bool, error) {
		retries++
		time.Sleep(1 * time.Millisecond)
		return true, nil
	}

	if err := Retry(time.Millisecond*5, time.Millisecond*1, fn); err != nil {
		assert.Equal(t, retries, 5)
		assert.Equal(t, err, errors.New(errors.FailedRequestTimeout, nil))
	}
}

func TestRetryFNRanBeforeDelay(t *testing.T) {
	start := time.Now()
	fn := func() (bool, error) {
		assert.WithinDuration(t, time.Now(), start, time.Millisecond)
		time.Sleep(1 * time.Second)
		return true, nil
	}

	Retry(time.Millisecond*5, time.Millisecond*1, fn)
}
