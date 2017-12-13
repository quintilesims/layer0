package scaler

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/common/models"
)

func TestDispatcherDispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobStore := mock_job.NewMockStore(ctrl)
	dispatcher := NewDispatcher(mockJobStore, 0)

	done := make(chan bool)
	recordCall := func(...interface{}) { done <- true }

	mockJobStore.EXPECT().
		Insert(models.ScaleEnvironmentJob, "env_id").
		Do(recordCall).
		Return("", nil)

	dispatcher.Dispatch("env_id")
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestDispatcherDispatchWithGracePeriodBuffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gracePeriod := time.Millisecond
	mockJobStore := mock_job.NewMockStore(ctrl)
	dispatcher := NewDispatcher(mockJobStore, gracePeriod)

	done := make(chan bool)
	recordCall := func(...interface{}) { done <- true }

	mockJobStore.EXPECT().
		Insert(models.ScaleEnvironmentJob, "env_id").
		Do(recordCall).
		Return("", nil)

	scheduled := time.Now()
	for i := 0; i < 100; i++ {
		dispatcher.Dispatch("env_id")
	}

	select {
	case <-done:
		if elapsed := time.Since(scheduled); elapsed < gracePeriod {
			t.Fatalf("Only %s has elapsed, expected (at least) %v", elapsed, gracePeriod)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
