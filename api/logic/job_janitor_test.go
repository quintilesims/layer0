package logic

import (
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
	"testing"
	"time"
)

func TestJobJanitorPulse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

	jobs := []*models.Job{
		&models.Job{
			JobID:       "old_job",
			TimeCreated: time.Now().Add(-(JOB_LIFETIME * 2)),
		},
		&models.Job{
			JobID:       "young_job",
			TimeCreated: time.Now(),
		},
	}

	jobLogicMock.EXPECT().
		SelectAll().
		Return(jobs, nil)

	jobLogicMock.EXPECT().
		Delete("old_job").
		Return(nil)

	janitor := NewJobJanitor(jobLogicMock)
	if err := janitor.pulse(); err != nil {
		t.Fatal(err)
	}
}
