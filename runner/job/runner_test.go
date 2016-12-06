package job

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/data/mock_data"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"testing"
	"time"
)

func getStubbedLogic(ctrl *gomock.Controller) *logic.Logic {
	mockJobData := mock_data.NewMockJobData(ctrl)
	mockJobData.EXPECT().
		UpdateJobStatus(gomock.Any(), gomock.Any()).
		AnyTimes()

	return logic.NewLogic(nil, nil, mockJobData, nil)
}

func stepWithError() Step {
	return Step{
		Name:    "step with error",
		Timeout: time.Second * 1,
		Action:  func(chan bool, JobContext) error { return fmt.Errorf("some error") },
	}
}

func TestRunnerLoad(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should load correct job",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobData := mock_data.NewMockJobData(ctrl)

				model := &models.Job{
					JobID:   "some_job_id",
					JobType: int64(types.DeleteEnvironmentJob),
				}

				mockJobData.EXPECT().GetJob("some_job_id").
					Return(model, nil)

				mockLogic := logic.NewLogic(nil, nil, mockJobData, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Load()
			},
		},
		testutils.TestCase{
			Name: "Should propagate JobData.GetJob error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobData := mock_data.NewMockJobData(ctrl)

				mockJobData.EXPECT().GetJob(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockLogic := logic.NewLogic(nil, nil, mockJobData, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)

				if err := runner.Load(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestRunnerRun_StepExecution(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
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
					Step{
						Name:    "step1",
						Timeout: time.Second * 1,
						Action: func(chan bool, JobContext) error {
							recorder.Call("step1")
							return nil
						},
					},
					Step{
						Name:    "step2",
						Timeout: time.Second * 1,
						Action: func(chan bool, JobContext) error {
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
		testutils.TestCase{
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
		testutils.TestCase{
			Name: "Should close quit channel after timeout",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")

				runner.Steps = []Step{
					Step{
						Name:    "timeout step",
						Timeout: time.Nanosecond * 0,
						Action: func(quit chan bool, c JobContext) error {
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

func TestRunnerRun_RollbackExecution(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call step.Rollback when error occurs",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")

				recorder := testutils.NewRecorder(ctrl)
				recorder.EXPECT().Call(gomock.Any())

				step := stepWithError()
				step.Rollback = func(JobContext) (JobContext, []Step, error) {
					recorder.Call("")
					return nil, nil, nil
				}

				runner.Steps = []Step{step}
				return runner
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		testutils.TestCase{
			Name: "Should call step.Rollback in correct order",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				logic := getStubbedLogic(ctrl)
				runner := NewJobRunner(logic, "")

				recorder := testutils.NewRecorder(ctrl)

				gomock.InOrder(
					recorder.EXPECT().Call("step2"),
					recorder.EXPECT().Call("step1"),
				)

				runner.Steps = []Step{
					Step{
						Name:    "step1",
						Timeout: time.Second * 1,
						Action:  func(chan bool, JobContext) error { return nil },
						Rollback: func(JobContext) (JobContext, []Step, error) {
							recorder.Call("step1")
							return nil, nil, nil
						},
					},
					Step{
						Name:    "step2",
						Timeout: time.Second * 1,
						Action:  func(chan bool, JobContext) error { return fmt.Errorf("some error") },
						Rollback: func(JobContext) (JobContext, []Step, error) {
							recorder.Call("step2")
							return nil, nil, nil
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
		testutils.TestCase{
			Name: "Should mark status to InProgress at start of Run",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobData := mock_data.NewMockJobData(ctrl)

				gomock.InOrder(
					mockJobData.EXPECT().UpdateJobStatus("some_job_id", types.InProgress),
					mockJobData.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.InProgress)).AnyTimes(),
				)

				mockLogic := logic.NewLogic(nil, nil, mockJobData, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		testutils.TestCase{
			Name: "Should mark status to Completed at end of Run without errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobData := mock_data.NewMockJobData(ctrl)

				gomock.InOrder(
					mockJobData.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.Completed)).AnyTimes(),
					mockJobData.EXPECT().UpdateJobStatus("some_job_id", types.Completed),
				)

				mockLogic := logic.NewLogic(nil, nil, mockJobData, nil)
				return NewJobRunner(mockLogic, "some_job_id")
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				runner := target.(*JobRunner)
				runner.Run()
			},
		},
		testutils.TestCase{
			Name: "Should mark status to Error at the end of Run with errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockJobData := mock_data.NewMockJobData(ctrl)

				gomock.InOrder(
					mockJobData.EXPECT().UpdateJobStatus(gomock.Any(), gomock.Not(types.Error)).AnyTimes(),
					mockJobData.EXPECT().UpdateJobStatus("some_job_id", types.Error),
				)

				mockLogic := logic.NewLogic(nil, nil, mockJobData, nil)
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
