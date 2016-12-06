package ecsbackend

import (
	"fmt"
	aws_ecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/mock_ecsbackend"
	"gitlab.imshealth.com/xfra/layer0/api/backend/mock_backend"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs/mock_cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2/mock_ec2"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs/mock_ecs"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

type MockECSServiceManager struct {
	ECS            *mock_ecs.MockProvider
	EC2            *mock_ec2.MockProvider
	CloudWatchLogs *mock_cloudwatchlogs.MockProvider
	ClusterScaler  *mock_ecsbackend.MockClusterScaler
	Backend        *mock_backend.MockBackend
}

func NewMockECSServiceManager(ctrl *gomock.Controller) *MockECSServiceManager {
	return &MockECSServiceManager{
		ECS:            mock_ecs.NewMockProvider(ctrl),
		EC2:            mock_ec2.NewMockProvider(ctrl),
		CloudWatchLogs: mock_cloudwatchlogs.NewMockProvider(ctrl),
		ClusterScaler:  mock_ecsbackend.NewMockClusterScaler(ctrl),
		Backend:        mock_backend.NewMockBackend(ctrl),
	}
}

func (this *MockECSServiceManager) Service() *ECSServiceManager {
	return NewECSServiceManager(this.ECS, this.EC2, this.CloudWatchLogs, this.ClusterScaler, this.Backend)
}

func TestGetService(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.DescribeService with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				mockService.ECS.EXPECT().
					DescribeService(environmentID.String(), serviceID.String()).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.GetService("envid", "svcid")
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted ids",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				mockService.ECS.EXPECT().
					DescribeService(gomock.Any(), gomock.Any()).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)

				service, err := manager.GetService("envid", "svcid")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(service.EnvironmentID, "envid")
				reporter.AssertEqual(service.ServiceID, "svcid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate ecs.DescribeService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				mockService.ECS.EXPECT().
					DescribeService(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)

				if _, err := manager.GetService("envid", "svcid"); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListServices(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.Helper_DescribeServices with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()
				service := ecs.NewService(clusterARN, serviceID.String())

				mockService.ECS.EXPECT().
					Helper_DescribeServices(id.PREFIX).
					Return([]*ecs.Service{service}, nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.ListServices()
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted ids",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()
				service := ecs.NewService(clusterARN, serviceID.String())

				mockService.ECS.EXPECT().
					Helper_DescribeServices(gomock.Any()).
					Return([]*ecs.Service{service}, nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)

				services, err := manager.ListServices()
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(len(services), 1)
				reporter.AssertEqual(services[0].EnvironmentID, "envid")
				reporter.AssertEqual(services[0].ServiceID, "svcid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate ecs.ListServices error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				mockService.ECS.EXPECT().
					Helper_DescribeServices(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)

				if _, err := manager.ListServices(); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteService(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.UpdateService, DescribeService, and DeleteService with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				mockService.ECS.EXPECT().
					UpdateService(environmentID.String(), serviceID.String(), nil, int64p(0)).
					Return(nil)

				mockService.ECS.EXPECT().
					DescribeService(environmentID.String(), serviceID.String()).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil)

				mockService.ECS.EXPECT().
					DeleteService(environmentID.String(), serviceID.String()).
					Return(nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.DeleteService("envid", "svcid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate unexpected aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSServiceManager {
					mockService := NewMockECSServiceManager(ctrl)

					environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
					clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
					serviceID := id.L0ServiceID("svcid").ECSServiceID()

					mockService.ECS.EXPECT().
						UpdateService(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockService.ECS.EXPECT().
						DescribeService(gomock.Any(), gomock.Any()).
						Return(ecs.NewService(clusterARN, serviceID.String()), g.Error()).
						AnyTimes()

					mockService.ECS.EXPECT().
						DeleteService(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					return mockService.Service()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSServiceManager)

				for i := 0; i < 3; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some eror"))

					manager := setup(g)
					if err := manager.DeleteService("envid", "svcid"); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateService(t *testing.T) {
	defer id.StubIDGeneration("svcid")()

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				deployID := id.L0DeployID("dplyid.1").ECSDeployID()
				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						Revision: int64p(1),
						Family:   stringp(deployID.FamilyName()),
					},
				}

				mockService.ECS.EXPECT().
					DescribeTaskDefinition(deployID.TaskDefinition()).
					Return(task, nil).
					AnyTimes()

				mockService.ClusterScaler.EXPECT().
					TriggerScalingAlgorithm(environmentID, &deployID, 1).
					Return(0, false, nil)

				mockService.ECS.EXPECT().CreateService(
					environmentID.String(),
					serviceID.String(),
					deployID.TaskDefinition(),
					int64(1),
					nil,
					nil).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.CreateService("svc_name", "envid", "dplyid.1", "")
			},
		},
		testutils.TestCase{
			Name: "Should propagate unexpected errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSServiceManager {
					mockService := NewMockECSServiceManager(ctrl)

					//deployID := id.L0DeployID("dplyid.1").ECSDeployID()
					environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
					clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
					serviceID := id.L0ServiceID("svcid").ECSServiceID()

					mockService.ClusterScaler.EXPECT().
						TriggerScalingAlgorithm(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(0, false, g.Error()).
						AnyTimes()

					mockService.ECS.EXPECT().CreateService(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any()).
						Return(ecs.NewService(clusterARN, serviceID.String()), g.Error()).
						AnyTimes()

					return mockService.Service()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSServiceManager)

				for i := 0; i < 2; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some eror"))

					manager := setup(g)
					if _, err := manager.CreateService("svc_name", "envid", "dplyid.1", ""); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateService(t *testing.T) {
	defer id.StubIDGeneration("svcid")()

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				deployID := id.L0DeployID("dplyid.1").ECSDeployID()
				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						Revision: int64p(1),
						Family:   stringp(deployID.FamilyName()),
					},
				}

				mockService.ECS.EXPECT().
					DescribeTaskDefinition(deployID.TaskDefinition()).
					Return(task, nil).
					AnyTimes()

				mockService.ClusterScaler.EXPECT().
					TriggerScalingAlgorithm(environmentID, &deployID, 1).
					Return(0, false, nil)

				mockService.ECS.EXPECT().UpdateService(
					environmentID.String(),
					serviceID.String(),
					stringp(deployID.TaskDefinition()),
					nil).
					Return(nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.updateService("envid", "svcid", "dplyid.1")
			},
		},
		testutils.TestCase{
			Name: "Should propagate unexpected errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSServiceManager {
					mockService := NewMockECSServiceManager(ctrl)

					mockService.ClusterScaler.EXPECT().
						TriggerScalingAlgorithm(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(0, false, g.Error()).
						AnyTimes()

					mockService.ECS.EXPECT().UpdateService(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					return mockService.Service()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSServiceManager)

				for i := 0; i < 2; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some eror"))

					manager := setup(g)
					if err := manager.updateService("envid", "svcid", "dplid.1"); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestScaleService(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call aws with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()
				deployID := id.L0DeployID("dplyid.1").ECSDeployID()
				service := ecs.NewService(clusterARN, serviceID.String())
				service.TaskDefinition = stringp(deployID.TaskDefinition())

				// 2nd call is at the end when we call GetService
				mockService.ECS.EXPECT().
					DescribeService(environmentID.String(), serviceID.String()).
					Return(service, nil).
					Times(2)

				mockService.ClusterScaler.EXPECT().
					TriggerScalingAlgorithm(environmentID, &deployID, 2).
					Return(0, false, nil)

				mockService.ECS.EXPECT().
					UpdateService(environmentID.String(), serviceID.String(), nil, int64p(2)).
					Return(nil)

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.ScaleService("envid", "svcid", 2)
			},
		},
		testutils.TestCase{
			Name: "Should propagate unexpected errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSServiceManager {
					mockService := NewMockECSServiceManager(ctrl)

					environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
					clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
					serviceID := id.L0ServiceID("svcid").ECSServiceID()
					deployID := id.L0DeployID("dplyid.1").ECSDeployID()
					service := ecs.NewService(clusterARN, serviceID.String())
					service.TaskDefinition = stringp(deployID.TaskDefinition())

					mockService.ECS.EXPECT().
						DescribeService(gomock.Any(), gomock.Any()).
						Return(service, g.Error()).
						AnyTimes()

					mockService.ClusterScaler.EXPECT().
						TriggerScalingAlgorithm(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(0, false, g.Error()).
						AnyTimes()

					mockService.ECS.EXPECT().
						UpdateService(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					return mockService.Service()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSServiceManager)

				for i := 0; i < 3; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some eror"))

					manager := setup(g)
					if _, err := manager.ScaleService("envid", "svc_id", 2); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestGetServiceLogs(t *testing.T) {
	tmp := GetLogs
	defer func() { GetLogs = tmp }()

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call GetLogs with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				mockService.ECS.EXPECT().
					DescribeService(environmentID.String(), serviceID.String()).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil).AnyTimes()

				// ensure we actually call GetLogs
				recorder := testutils.NewRecorder(ctrl)
				recorder.EXPECT().Call("")

				GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, tail int) ([]*models.LogFile, error) {
					recorder.Call("")
					reporter.AssertEqual(tail, 100)
					return nil, nil
				}

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)
				manager.GetServiceLogs("envid", "svcid", 100)
			},
		},
		testutils.TestCase{
			Name: "Should propagate GetLogs error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockService := NewMockECSServiceManager(ctrl)

				environmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
				serviceID := id.L0ServiceID("svcid").ECSServiceID()

				mockService.ECS.EXPECT().
					DescribeService(environmentID.String(), serviceID.String()).
					Return(ecs.NewService(clusterARN, serviceID.String()), nil).AnyTimes()

				// ensure we actually call GetLogs
				recorder := testutils.NewRecorder(ctrl)
				recorder.EXPECT().Call("")

				GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, tail int) ([]*models.LogFile, error) {
					recorder.Call("")
					return nil, fmt.Errorf("some error")
				}

				return mockService.Service()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSServiceManager)

				if _, err := manager.GetServiceLogs("envid", "svcid", 100); err == nil {
					reporter.Fatalf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
