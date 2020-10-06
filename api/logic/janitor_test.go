package logic

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
)

func TestJanitorPulse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

	jobs := []*models.Job{
		{
			JobID:       "old_job",
			TimeCreated: time.Now().Add(-(JOB_LIFETIME * 2)),
		},
		{
			JobID:       "young_job",
			TimeCreated: time.Now(),
		},
	}

	jobLogicMock.EXPECT().
		ListJobs().
		Return(jobs, nil)

	jobLogicMock.EXPECT().
		Delete("old_job").
		Return(nil)

	janitor := NewJobJanitor(jobLogicMock)
	if err := janitor.pulse(); err != nil {
		t.Fatal(err)
	}
}
