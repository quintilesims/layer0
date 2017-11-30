package job

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/common/models"
)

func TestJanitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expiry := time.Hour * 12
	mockJobStore := mock_job.NewMockStore(ctrl)
	janitor := NewJanitor(mockJobStore, expiry)
	now := time.Now()

	// There is no Sub function that returns Time in time package,
	// so we must negate the expiry
	jobs := []*models.Job{
		{
			JobID:   "delete",
			Created: now.Add(-(expiry + time.Second)), // 1 Second over expiry
		},
		{
			JobID:   "keep",
			Created: now.Add(-(expiry - time.Second)), // 1 Second before expiry
		},
	}

	mockJobStore.EXPECT().
		SelectAll().
		Return(jobs, nil)

	mockJobStore.EXPECT().
		Delete("delete").
		Return(nil)

	if err := janitor.Run(); err != nil {
		t.Fatal(err)
	}

}
