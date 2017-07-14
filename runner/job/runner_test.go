package job

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/db/job_store/mock_job_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
)

func getStubbedLogic(ctrl *gomock.Controller) *logic.Logic {
	mockJobStore := mock_job_store.NewMockJobStore(ctrl)
	mockJobStore.EXPECT().
		UpdateJobStatus(gomock.Any(), gomock.Any()).
		AnyTimes()

	return logic.NewLogic(nil, mockJobStore, nil, nil)
}

func stepWithError() Step {
	return Step{
		Name:    "step with error",
		Timeout: time.Second * 1,
		Action:  func(chan bool, *JobContext) error { return fmt.Errorf("some error") },
	}
}

func TestRunnerLoad(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should load correct job",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				model := &models.Job{
					JobID:   "some_job_id",
					JobType: int64(types.DeleteEnvironmentJob),
				}

				mockJobStore.EXPECT().SelectByID("some_job_id").
					Return(model, nil)

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Load()
			},
		},
		{
			Name: "Should propagate JobStore.SelectByID error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				mockJobStore.EXPECT().SelectByID(gomock.Any()).
					Return(nil, fmt.Errorf("some error")).
					AnyTimes()

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)

				if err := runner.Load(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		{
			Name: "Should retry after failed job load",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				model := &models.Job{
					JobID:   "some_job_id",
					JobType: int64(types.DeleteEnvironmentJob),
				}

				gomock.InOrder(
					mockJobStore.EXPECT().SelectByID("some_job_id").Return(nil, fmt.Errorf("some error")),
					mockJobStore.EXPECT().SelectByID("some_job_id").Return(model, nil),
				)

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Load()
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestRunnerRun_StepExecution(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should run steps in correct order",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")
				recorder := testutils.NewRecorder(ctrl)

				gomock.InOrder(
					recorder.EXPECT().Call("step1"),
					recorder.EXPECT().Call("step2"),
				)

				runner.Steps = []Step{
					{
						Name:    "step1",
						Timeout: time.Second * 1,
						Action: func(chan bool, *JobContext) error {
							recorder.Call("step1")
							return nil
						},
					},
					{
						Name:    "step2",
						Timeout: time.Second * 1,
						Action: func(chan bool, *JobContext) error {
							recorder.Call("step2")
							return nil
						},
					},
				}

				return runner
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		{
			Name: "Should return step error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")

				runner.Steps = []Step{stepWithError()}
				return runner
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)

				if err := runner.Run(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		{
			Name: "Should close quit channel after timeout",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")

				runner.Steps = []Step{
					{
						Name:    "timeout step",
						Timeout: time.Nanosecond * 0,
						Action: func(quit chan bool, c *JobContext) error {
							select {
							case <-quit:
								return nil
							case <-time.After(time.Second * 1):
								t.Errorf("quit channel was not closed")
								return nil
							}
						},
					},
				}

				return runner
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestRunnerRun_JobStateManagement(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should mark status to InProgress at start of Run",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				gomock.InOrder(
					mockJobStore.EXPECT().UpdateJobStatus("some_job_id", types.InProgress),
					mockJobStore.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.InProgress)).AnyTimes(),
				)

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		{
			Name: "Should mark status to Completed at end of Run without errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				gomock.InOrder(
					mockJobStore.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.Completed)).AnyTimes(),
					mockJobStore.EXPECT().UpdateJobStatus("some_job_id", types.Completed),
				)

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		{
			Name: "Should mark status to Error at the end of Run with errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobStore := mock_job_store.NewMockJobStore(ctrl)

				gomock.InOrder(
					mockJobStore.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.Error)).AnyTimes(),
					mockJobStore.EXPECT().UpdateJobStatus("some_job_id", types.Error),
				)

				mockLogic := logic.NewLogic(nil, mockJobStore, nil, nil)
				runner := NewJobRunner(mockLogic, "some_job_id")

				runner.Steps = []Step{stepWithError()}
				return runner
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
	}

	testutils.RunTests(t, testCases)
}
