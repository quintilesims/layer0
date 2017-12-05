package lock

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/stretchr/testify/assert"
)

func newTestLock(t *testing.T, expiry time.Duration) *DynamoLock {
	session := config.GetTestAWSSession()
	table := os.Getenv(config.ENVVAR_TEST_AWS_DYNAMO_LOCK_TABLE)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.ENVVAR_TEST_AWS_DYNAMO_LOCK_TABLE)
	}

	lock := NewDynamoLock(session, table, expiry)
	return lock
}

func TestDynamoLock_acquireAfterRelease(t *testing.T) {
	lock := newTestLock(t, time.Hour)
	if err := lock.Release("test"); err != nil {
		t.Fatal(err)
	}

	acquired, err := lock.Acquire("test")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, acquired)
}

func TestDynamoLock_acquireAfterExpiry(t *testing.T) {
	expiry := time.Nanosecond
	lock := newTestLock(t, expiry)
	if _, err := lock.Acquire("test"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(expiry + 1)
	acquired, err := lock.Acquire("test")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, acquired)
}

func TestDynamoLock_acquireFailureOnContention(t *testing.T) {
	lock := newTestLock(t, time.Hour)
	if _, err := lock.Acquire("test"); err != nil {
		t.Fatal(err)
	}

	acquired, err := lock.Acquire("test")
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, acquired)
}

func TestDynamoLock_release(t *testing.T) {
	lock := newTestLock(t, time.Hour)
	if _, err := lock.Acquire("test"); err != nil {
		t.Fatal(err)
	}

	if err := lock.Release("test"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoLock_releaseWhenDoesNotExist(t *testing.T) {
	lock := newTestLock(t, time.Hour)
	if err := lock.Release("test"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoLock_locksAreDiscrete(t *testing.T) {
	for i := 0; i < 10; i++ {
		lockID := strconv.Itoa(i)
		t.Run(lockID, func(t *testing.T) {
			lock := newTestLock(t, time.Hour)
			if err := lock.Release(lockID); err != nil {
				t.Fatal(err)
			}

			acquired, err := lock.Acquire(lockID)
			if err != nil {
				t.Fatal(err)
			}

			assert.True(t, acquired)

			// release after an async delay - if locks are not
			// discrete, then the next acquisition(s) would fail
			time.AfterFunc(time.Second, func() {
				if err := lock.Release(lockID); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}
