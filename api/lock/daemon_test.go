package lock

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaemonFN(t *testing.T) {
	lock := newTestLock(t, 0)
	expiry := time.Second
	daemonFN := NewDaemonFN(lock, expiry)

	for i := 0; i < 5; i++ {
		lockID := strconv.Itoa(i)
		if _, err := lock.Acquire(lockID); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(expiry)
	if _, err := lock.Acquire("keep"); err != nil {
		t.Fatal(err)
	}

	if err := daemonFN(); err != nil {
		t.Fatal(err)
	}

	locks, err := lock.List()
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, locks, 1)
	assert.Equal(t, "keep", locks[0].LockID)
}
