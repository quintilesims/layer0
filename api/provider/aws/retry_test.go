package aws

import (
	"fmt"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetry_Timeout(t *testing.T) {
	fn := func() error {
		return fmt.Errorf("Tick")
	}

	if err := retry(20*time.Millisecond, 5*time.Millisecond, fn); err != nil {
		assert.Equal(t, err, errors.New(errors.FailedRequestTimeout, nil))
	}
}

func TestRetry_NoTimeout(t *testing.T) {
	fn := func() error {
		return nil
	}

	if err := retry(20*time.Millisecond, 5*time.Millisecond, fn); err != nil {
		t.Fatal(err)
	}
}
