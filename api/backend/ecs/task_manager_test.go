package ecsbackend

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_ecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/backend/mock_backend"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs/mock_cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/ecs/mock_ecs"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/stretchr/testify/assert"
)

type MockECSTaskManager struct {
	ECS            *mock_ecs.MockProvider
	CloudWatchLogs *mock_cloudwatchlogs.MockProvider
	Backend        *mock_backend.MockBackend
}

func NewMockECSTaskManager(ctrl *gomock.Controller) *MockECSTaskManager {
	return &MockECSTaskManager{
		ECS:            mock_ecs.NewMockProvider(ctrl),
		CloudWatchLogs: mock_cloudwatchlogs.NewMockProvider(ctrl),
		Backend:        mock_backend.NewMockBackend(ctrl),
	}
}

func (this *MockECSTaskManager) Task() *ECSTaskManager {
	taskManager := NewECSTaskManager(this.ECS, this.CloudWatchLogs, this.Backend)
	return taskManager
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	environmentID := id.L0EnvironmentID("env_id")
	task := &ecs.Task{
		&aws_ecs.Task{
			LastStatus: aws.String("RUNNING"),
		},
	}

	mockTask := NewMockECSTaskManager(ctrl)
	mockTask.ECS.EXPECT().
		DescribeTask(environmentID.ECSEnvironmentID().String(), "task_arn").
		Return(task, nil)

	result, err := mockTask.Task().GetTask("env_id", "task_arn")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(1), result.RunningCount)
	assert.Equal(t, int64(0), result.PendingCount)
	assert.Len(t, result.Copies, 1)
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ecsEnvironmentIDs := []id.ECSEnvironmentID{
		id.ECSEnvironmentID("env_id1"),
		id.ECSEnvironmentID("env_id2"),
	}

	mockTask := NewMockECSTaskManager(ctrl)
	mockTask.Backend.EXPECT().
		ListEnvironments().
		Return(ecsEnvironmentIDs, nil)

	for i, ecsEnvironmentID := range ecsEnvironmentIDs {
		arn := fmt.Sprintf("arn_%d", i)

		mockTask.ECS.EXPECT().
			ListClusterTaskARNs(ecsEnvironmentID.String(), id.PREFIX).
			Return([]string{arn}, nil)
	}

	result, err := mockTask.Task().ListTasks()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"arn_0",
		"arn_1",
	}

	assert.Equal(t, expected, result)
}

func TestDeleteTask(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should use proper params in dependent calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				taskID := id.L0TaskID("tskid").ECSTaskID()

				mockTask.ECS.EXPECT().
					ListTasks(environmentID.String(), nil, gomock.Any(), stringp(taskID.String()), nil).
					Return(nil, nil).
					Times(2)

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)
				manager.DeleteTask("envid", "tskid")
			},
		},
		{
			Name: "Should propagate ecs.ListTasks error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				mockTask.ECS.EXPECT().
					ListTasks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)

				if err := manager.DeleteTask("envid", "tskid"); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateTask(t *testing.T) {
	defer id.StubIDGeneration("tskid")()

	testCases := []testutils.TestCase{
		{
			Name: "Should use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				deployID := id.L0DeployID("dplyid.1").ECSDeployID()
				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						Revision: int64p(1),
						Family:   stringp(deployID.FamilyName()),
					},
				}

				mockTask.ECS.EXPECT().
					DescribeTaskDefinition(deployID.TaskDefinition()).
					Return(task, nil).
					AnyTimes()

				mockTask.ECS.EXPECT().RunTask(
					environmentID.String(),
					deployID.TaskDefinition(),
					int64(1),
					stringp(id.PREFIX),
					[]*ecs.ContainerOverride{},
				).Return(nil, nil, nil)

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)
				manager.CreateTask("envid", "tsk_name", "dplyid.1", nil)
			},
		},
		{
			Name: "Should not create cloudwatch logs group if disableLogging is true",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				mockTask.ECS.EXPECT().RunTask(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil, nil)

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)
				manager.CreateTask("envid", "tsk_name", "dplyid.1", nil)
			},
		},
		{
			Name: "Should add task to scheduler when cluster capacity is low",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				err := awserr.New("", "No Container Instances were found in your cluster", nil)
				mockTask.ECS.EXPECT().RunTask(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil, err)

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)
				manager.CreateTask("envid", "tsk_name", "dplyid.1", nil)
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestGetTaskLogs(t *testing.T) {
	tmp := GetLogs
	defer func() { GetLogs = tmp }()

	testCases := []testutils.TestCase{
		{
			Name: "Should call GetLogs with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				taskID := id.L0TaskID("tskid").ECSTaskID()
				mockTask.ECS.EXPECT().
					ListTasks(environmentID.String(), nil, gomock.Any(), stringp(taskID.String()), nil).
					Return([]*string{}, nil).
					AnyTimes()

				// ensure we actually call GetLogs
				recorder := testutils.NewRecorder(ctrl)
				recorder.EXPECT().Call("")

				GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, start, end string, tail int) ([]*models.LogFile, error) {
					recorder.Call("")
					reporter.AssertEqual(tail, 100)
					reporter.AssertEqual(start, "start")
					reporter.AssertEqual(end, "end")
					return nil, nil
				}

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)
				manager.GetTaskLogs("envid", "tskid", "start", "end", 100)
			},
		},
		{
			Name: "Should propagate GetLogs error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockTask := NewMockECSTaskManager(ctrl)

				mockTask.ECS.EXPECT().
					ListTasks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*string{}, nil).
					AnyTimes()

				// ensure we actually call GetLogs
				recorder := testutils.NewRecorder(ctrl)
				recorder.EXPECT().Call("")

				GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, start, end string, tail int) ([]*models.LogFile, error) {
					recorder.Call("")
					return nil, fmt.Errorf("some error")
				}

				return mockTask.Task()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSTaskManager)

				if _, err := manager.GetTaskLogs("envid", "tskid", "start", "end", 100); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
