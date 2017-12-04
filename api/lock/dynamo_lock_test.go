package lock

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/stretchr/testify/assert"
)

func newTestLock(t *testing.T, lockID string, expiry time.Duration) *DynamoLock {
	session := config.GetTestAWSSession()
	table := os.Getenv(config.ENVVAR_TEST_AWS_DYNAMO_LOCK_TABLE)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.ENVVAR_TEST_AWS_DYNAMO_LOCK_TABLE)
	}

	lock := NewDynamoLock(session, table, lockID, expiry)
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}

	return lock
}

func TestDynamoLock_acquireWhenDoesNotExist(t *testing.T) {
	lock := newTestLock(t, "test", time.Hour)
	acquired, err := lock.Acquire()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, acquired)
}

func TestDynamoLock_acquireAfterRelease(t *testing.T) {
	lock := newTestLock(t, "test", time.Hour)
	if _, err := lock.Acquire(); err != nil {
		t.Fatal(err)
	}

	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}

	acquired, err := lock.Acquire()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, acquired)
}

func TestDynamoLock_acquireAfterExpiry(t *testing.T) {
	expiry := time.Nanosecond
	lock := newTestLock(t, "test", expiry)
	if _, err := lock.Acquire(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(expiry + 1)
	acquired, err := lock.Acquire()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, acquired)
}

func TestDynamoLock_acquireFailureOnContention(t *testing.T) {
	lock := newTestLock(t, "test", time.Hour)
	if _, err := lock.Acquire(); err != nil {
		t.Fatal(err)
	}

	acquired, err := lock.Acquire()
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, acquired)
}

func TestDynamoLock_release(t *testing.T) {
	lock := newTestLock(t, "test", time.Hour)
	if _, err := lock.Acquire(); err != nil {
		t.Fatal(err)
	}

	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoLock_releaseWhenDoesNotExist(t *testing.T) {
	lock := newTestLock(t, "test", time.Hour)
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoLock_locksAreDiscrete(t *testing.T) {
	for i := 0; i < 10; i++ {
		lockID := strconv.Itoa(i)
		t.Run(lockID, func(t *testing.T) {
			lock := newTestLock(t, lockID, time.Hour)
			acquired, err := lock.Acquire()
			if err != nil {
				t.Fatal(err)
			}

			assert.True(t, acquired)

			// release after an async delay - if locks are not
			// discrete, then the next acquisition(s) would fail
			time.AfterFunc(time.Second, func() {
				if err := lock.Release(); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}

// unlock success on: exists or doesnn't exit
