package daemon

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/lock/mock_lock"
	"github.com/stretchr/testify/assert"
)

func TestDaemonRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock_lock.NewMockLock(ctrl)
	mockLock.EXPECT().
		Acquire("lock_id").
		Return(true, nil)

	mockLock.EXPECT().
		Release("lock_id").
		Return(nil)

	var called bool
	daemon := NewDaemon("", "lock_id", mockLock, func() error {
		called = true
		return nil
	})

	daemon.Run()
	assert.True(t, called)
}

func TestDaemonHonorsLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock_lock.NewMockLock(ctrl)
	mockLock.EXPECT().
		Acquire("lock_id").
		Return(false, nil)

	daemon := NewDaemon("", "lock_id", mockLock, func() error {
		t.Fatalf("The Daemon's function was called")
		return nil
	})

	daemon.Run()
}

func TestDaemonRunEvery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLock := mock_lock.NewMockLock(ctrl)
	mockLock.EXPECT().
		Acquire(gomock.Any()).
		Return(true, nil).
		AnyTimes()

	mockLock.EXPECT().
		Release(gomock.Any()).
		Return(nil).
		AnyTimes()

	c := make(chan bool)
	daemon := NewDaemon("", "", mockLock, func() error {
		c <- true
		return nil
	})

	ticker := daemon.RunEvery(time.Nanosecond)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		<-c
	}
}
