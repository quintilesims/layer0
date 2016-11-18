package logic

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/api/logic/mock_logic"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"gitlab.imshealth.com/xfra/layer0/common/types"
	"testing"
	"time"
)

func TestGetJob(t *testing.T) {
	model := &models.Job{
		JobID:       "some_job_id",
		JobStatus:   1,
		JobType:     1,
		Request:     "some_request",
		TimeCreated: time.Now(),
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should pass JobID to data, return correct model",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					GetJob("some_job_id").
					Return(model, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				output, err := job.GetJob("some_job_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(model, output)
			},
		},
		testutils.TestCase{
			Name: "Should propagate data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					GetJob("some_job_id").
					Return(nil, fmt.Errorf("some error"))

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				if _, err := job.GetJob("some_job_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListJobs(t *testing.T) {
	models := []*models.Job{
		&models.Job{
			JobID:       "some_job_id1",
			JobStatus:   1,
			JobType:     1,
			Request:     "some_request1",
			TimeCreated: time.Now(),
		},
		&models.Job{
			JobID:       "some_job_id2",
			JobStatus:   2,
			JobType:     2,
			Request:     "some_request2",
			TimeCreated: time.Now(),
		},
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should return correct models",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					ListJobs().
					Return(models, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				output, err := job.ListJobs()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(models, output)
			},
		},
		testutils.TestCase{
			Name: "Should propagate data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					ListJobs().
					Return(nil, fmt.Errorf("some error"))

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				if _, err := job.ListJobs(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateJob(t *testing.T) {
	deployModel := &models.Deploy{
		DeployID: "some_deploy_id",
	}

	taskModel := &models.Task{
		TaskID: "some_task_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should pass correct params to jobData.CreateJob",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				var checkJob = func(job *models.Job) {
					reporter.AssertEqual(job.JobStatus, int64(types.Pending))
					reporter.AssertEqual(job.JobType, int64(types.DeleteEnvironmentJob))
					reporter.AssertEqual(job.Request, "some_request")
				}

				mockLogic.Job.EXPECT().
					CreateJob(gomock.Any()).
					Do(checkJob).
					Return(nil)

				mockLogic.Tag.EXPECT().
					Make(gomock.Any())

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(deployModel, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Return(taskModel, nil)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)
				if _, err := job.CreateJob(types.DeleteEnvironmentJob, "some_request"); err != nil {
					reporter.Error(err)
				}
			},
		},
		testutils.TestCase{
			Name: "Should pass correct params to deployLogic.CreateDeploy",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					CreateJob(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					Make(gomock.Any())

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Return(taskModel, nil)

				var checkDeploy = func(deploy models.CreateDeployRequest) {
					var dockerrun Dockerrun

					if err := json.Unmarshal(deploy.Dockerrun, &dockerrun); err != nil {
						reporter.Error(err)
					}

					var exists bool
					for _, env := range dockerrun.Containers[0].Environment {
						if env.Name == "LAYER0_JOB_ID" && env.Value != "" {
							exists = true
						}
					}

					reporter.AssertEqual(exists, true)
				}

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Do(checkDeploy).
					Return(deployModel, nil)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)
				job.CreateJob(types.DeleteEnvironmentJob, "some_request")
			},
		},
		testutils.TestCase{
			Name: "Should pass correct params to taskLogic.CreateTask",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					CreateJob(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					Make(gomock.Any())

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(deployModel, nil)

				var checkTask = func(task models.CreateTaskRequest) {
					reporter.AssertEqual(task.DeployID, deployModel.DeployID)
					reporter.AssertEqual(task.EnvironmentID, config.API_SERVICE_ID)
					reporter.AssertEqual(task.Copies, int64(1))
				}

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Do(checkTask).
					Return(taskModel, nil)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)
				job.CreateJob(types.DeleteEnvironmentJob, "some_request")
			},
		},
		testutils.TestCase{
			Name: "Generated job should have proper task id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					CreateJob(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					Make(gomock.Any())

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(deployModel, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Return(taskModel, nil)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				output, err := job.CreateJob(types.DeleteEnvironmentJob, "some_request")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(output.TaskID, taskModel.TaskID)
			},
		},
		testutils.TestCase{
			Name: "Should propagate jobData.CreateJob error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.Job.EXPECT().
					CreateJob(gomock.Any()).
					Return(fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(deployModel, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Return(taskModel, nil)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				if _, err := job.CreateJob(types.DeleteEnvironmentJob, "some_request"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate deployLogic.CreateDeploy error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				mockTask := mock_logic.NewMockTaskLogic(ctrl)

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				if _, err := job.CreateJob(types.DeleteEnvironmentJob, "some_request"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate taskLogic.CreateTask error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(deployModel, nil)

				mockTask := mock_logic.NewMockTaskLogic(ctrl)
				mockTask.EXPECT().
					CreateTask(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewL0JobLogic(mockLogic.Logic(), mockTask, mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				job := target.(*L0JobLogic)

				if _, err := job.CreateJob(types.DeleteEnvironmentJob, "some_request"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
