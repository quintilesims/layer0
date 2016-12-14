package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

func TestJobJanitorPulse(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should delete only old jobs",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
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

				return NewJobJanitor(jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				janitor := target.(*JobJanitor)
				janitor.pulse()
			},
		},
		testutils.TestCase{
			Name: "Should propagate Delete error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				jobs := []*models.Job{
					&models.Job{
						JobID:       "some_job",
						TimeCreated: time.Now().Add(-(JOB_LIFETIME * 2)),
					},
				}

				jobLogicMock.EXPECT().
					SelectAll().
					Return(jobs, nil)

				jobLogicMock.EXPECT().
					Delete(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return NewJobJanitor(jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				janitor := target.(*JobJanitor)

				if err := janitor.pulse(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate SelectAll error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				jobLogicMock.EXPECT().
					SelectAll().
					Return(nil, fmt.Errorf("some error"))

				return NewJobJanitor(jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				janitor := target.(*JobJanitor)

				if err := janitor.pulse(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
