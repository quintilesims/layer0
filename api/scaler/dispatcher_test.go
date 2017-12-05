package scaler

import (
	"strconv"
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
	dispatcher := NewDispatcher(mockJobStore, time.Nanosecond)

	for i := 0; i < 10; i++ {
		environmentID := strconv.Itoa(i)

		mockJobStore.EXPECT().
			Insert(models.ScaleEnvironmentJob, environmentID).
			Return("", nil)

		dispatcher.Dispatch(environmentID)
	}

	time.Sleep(time.Millisecond)
}

func TestDispatcherDispatchWithGracePeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobStore := mock_job.NewMockStore(ctrl)
	dispatcher := NewDispatcher(mockJobStore, time.Millisecond)

	mockJobStore.EXPECT().
		Insert(models.ScaleEnvironmentJob, "env_id").
		Return("", nil)

	for i := 0; i < 100; i++ {
		dispatcher.Dispatch("env_id")
	}

	time.Sleep(time.Millisecond * 2)
}
