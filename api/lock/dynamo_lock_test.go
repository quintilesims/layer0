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
	table := os.Getenv(config.FlagTestAWSLockTable.EnvVar)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.FlagTestAWSLockTable.EnvVar)
	}

	lock := NewDynamoLock(session, table, expiry)
	if err := lock.clear(); err != nil {
		t.Fatal(err)
	}

	return lock
}

func TestDynamoLock_acquireAfterRelease(t *testing.T) {
	lock := newTestLock(t, time.Hour)
	if _, err := lock.Acquire("test"); err != nil {
		t.Fatal(err)
	}

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
	lock := newTestLock(t, time.Hour)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(lockID string) {
			acquired, err := lock.Acquire(lockID)
			if err != nil {
				t.Fatal(err)
			}

			assert.True(t, acquired)
			<-done
		}(strconv.Itoa(i))
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
