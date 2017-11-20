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

	mockJobStore := mock_job.NewMockStore(ctrl)
	janitor := NewJanitor(mockJobStore, time.Hour*24)

	expiry := -24 * time.Hour
	now := time.Now()

	jobs := []*models.Job{
		{
			JobID:   "delete",
			Created: now.Add(expiry - time.Hour),
		},
		{
			JobID:   "keep",
			Created: now.Add(expiry + time.Hour),
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
