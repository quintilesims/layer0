package ecsbackend

import (
	"fmt"
	aws_ecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/mock_ecsbackend"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs/mock_ecs"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

type MockECSDeployManager struct {
	ECS           *mock_ecs.MockProvider
	ClusterScaler *mock_ecsbackend.MockClusterScaler
}

func NewMockECSDeployManager(ctrl *gomock.Controller) *MockECSDeployManager {
	return &MockECSDeployManager{
		ECS:           mock_ecs.NewMockProvider(ctrl),
		ClusterScaler: mock_ecsbackend.NewMockClusterScaler(ctrl),
	}
}

func (this *MockECSDeployManager) Deploy() *ECSDeployManager {
	return NewECSDeployManager(this.ECS, this.ClusterScaler)
}

func TestGetDeploy(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.DescribeTaskDefinition with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				deployID := id.L0DeployID("some_id").ECSDeployID()

				mockDeploy.ECS.EXPECT().
					DescribeTaskDefinition(deployID.TaskDefinition()).
					Return(nil, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.GetDeploy("some_id")
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted deploy id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				// note that GetDeploy doesn't generate the id from the returned TaskDefinition
				// this call is simply used to check for existence
				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						Revision: int64p(1),
						Family:   stringp("l0-prefix-family"),
					},
				}

				mockDeploy.ECS.EXPECT().
					DescribeTaskDefinition(gomock.Any()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				deploy, err := manager.GetDeploy("some_id")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(deploy.DeployID, "some_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate ecs.DescribeTaskDefinition error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					DescribeTaskDefinition(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				if _, err := manager.GetDeploy("some_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListDeploys(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.Helper_ListTaskDefinitions with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					Helper_ListTaskDefinitions(id.PREFIX).
					Return(nil, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.ListDeploys()
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted deploy ids and versions",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				deployAlpha := fmt.Sprintf("%salpha:4", id.PREFIX)
				deployBeta := fmt.Sprintf("%sbeta:18", id.PREFIX)

				deployARNs := []*string{
					stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + deployAlpha),
					stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + deployBeta),
				}

				mockDeploy.ECS.EXPECT().
					Helper_ListTaskDefinitions(gomock.Any()).
					Return(deployARNs, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				deploys, err := manager.ListDeploys()
				if err != nil {
					reporter.Fatal(err)
				}

				deployAlpha := &models.Deploy{
					DeployID: "alpha.4",
					Version:  "4",
				}

				deployBeta := &models.Deploy{
					DeployID: "beta.18",
					Version:  "18",
				}

				reporter.AssertEqual(len(deploys), 2)
				reporter.AssertInSlice(deployAlpha, deploys)
				reporter.AssertInSlice(deployBeta, deploys)
			},
		},
		testutils.TestCase{
			Name: "Should propagate iam.ListDeploys error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					Helper_ListTaskDefinitions(id.PREFIX).
					Return(nil, fmt.Errorf("some_error"))

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				if _, err := manager.ListDeploys(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteDeploy(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.DeleteTaskDefinition with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				ecsDeployID := id.L0DeployID("some_id").ECSDeployID()

				mockDeploy.ECS.EXPECT().
					DeleteTaskDefinition(ecsDeployID.TaskDefinition()).
					Return(nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				deploy := target.(*ECSDeployManager)
				deploy.DeleteDeploy("some_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate ecs.DeleteTaskDefinition error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					DeleteTaskDefinition(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				deploy := target.(*ECSDeployManager)

				if err := deploy.DeleteDeploy("some_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateDeploy(t *testing.T) {
	dockerrun := []byte(`
		{
			"ContainerDefinitions": [
				{	
					"name": "test",
					"image": "d.ims.io/xfra/test",
					"essential": true,
					"memory": 128
				}
			],
			"Volumes": [
				{
					"name": "test",
					"host": {
      						"sourcePath": "some_path"
    					}

				}
			],
			"Family": "",
			"NetworkMode": "host",
			"TaskRoleARN": "some_role"	
		}
	`)

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call ecs.RegisterTaskDefinition with proper id param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				deployID := id.L0DeployID("some_name").ECSDeployID()
				taskDefinition := fmt.Sprintf("%ssome_name:1", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(deployID.String(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrun)
			},
		},
		testutils.TestCase{
			Name: "Should marshal dockerrun correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				taskDefinition := fmt.Sprintf("%ssome_name:1", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(family, taskRoleARN, network string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume) {
						reporter.AssertEqual(network, "host")
						reporter.AssertEqual(taskRoleARN, "some_role")
						reporter.AssertEqual(len(containers), 1)
						reporter.AssertEqual(*containers[0].Name, "test")
						reporter.AssertEqual(*volumes[0].Name, "test")

					}).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrun)
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted id and version",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				deployID := id.L0DeployID("some_name").ECSDeployID()
				taskDefinition := fmt.Sprintf("%ssome_name:2", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(deployID.String(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				deploy, err := manager.CreateDeploy("some_name", dockerrun)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(deploy.DeployID, "some_name.2")
				reporter.AssertEqual(deploy.Version, "2")
			},
		},
		testutils.TestCase{
			Name: "Should error if deployName contains '.'",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)
				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				if _, err := manager.CreateDeploy("some.name", dockerrun); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate ecs.RegisterTaskDefinition error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				if _, err := manager.CreateDeploy("some_name", dockerrun); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
