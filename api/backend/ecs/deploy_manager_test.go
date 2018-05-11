package ecsbackend

import (
	"fmt"
	"os"
	"testing"

	aws_ecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/ecs/mock_ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

type MockECSDeployManager struct {
	ECS *mock_ecs.MockProvider
}

func NewMockECSDeployManager(ctrl *gomock.Controller) *MockECSDeployManager {
	return &MockECSDeployManager{
		ECS: mock_ecs.NewMockProvider(ctrl),
	}
}

func (this *MockECSDeployManager) Deploy() *ECSDeployManager {
	return NewECSDeployManager(this.ECS)
}

func getTaskDefARN(name string, version int) *string {
	str := fmt.Sprintf("arn:aws:ecs:region:account_id:task-definition/%s%s:%d", id.PREFIX, name, version)
	return &str
}

func TestGetDeploy(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should call ecs.DescribeTaskDefinition with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				deployID := id.L0DeployID("some_id.1").ECSDeployID()

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: getTaskDefARN("some_id", 1),
						Revision:          int64p(1),
						Family:            stringp("l0-prefix-family"),
					},
				}

				mockDeploy.ECS.EXPECT().
					DescribeTaskDefinition(deployID.TaskDefinition()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.GetDeploy("some_id.1")
			},
		},
		{
			Name: "Should return layer0-formatted deploy id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: getTaskDefARN("some_id", 1),
						Revision:          int64p(1),
						Family:            stringp("l0-prefix-family"),
					},
				}

				mockDeploy.ECS.EXPECT().
					DescribeTaskDefinition(gomock.Any()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				deploy, err := manager.GetDeploy("some_id.1")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(deploy.DeployID, "some_id.1")
			},
		},
		{
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
		{
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
		{
			Name: "Should return layer0-formatted deploy ids and versions",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				arns := []*string{
					getTaskDefARN("a", 1),
					getTaskDefARN("b", 2),
				}

				mockDeploy.ECS.EXPECT().
					Helper_ListTaskDefinitions(gomock.Any()).
					Return(arns, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)

				deploys, err := manager.ListDeploys()
				if err != nil {
					reporter.Fatal(err)
				}

				expected := []*models.Deploy{
					{DeployID: "a.1"},
					{DeployID: "b.2"},
				}

				reporter.AssertEqual(deploys, expected)
			},
		},
		{
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
		{
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
		{
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
					"image": "quintilesims/test",
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

	dockerrunWithPCs := []byte(`
		{
			"ContainerDefinitions": [
				{
					"name": "test",
					"image": "quintilesims/test",
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
			"TaskRoleARN": "some_role",
			"PlacementConstraints": [
				{
					"expression": "attribute:ecs.instance-type =~ m3.*",
					"type": "memberOf"
				}
			]
		}
	`)

	testCases := []testutils.TestCase{
		{
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
					RegisterTaskDefinition(deployID.String(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrun)
			},
		},
		{
			Name: "Should marshal dockerrun without placement constraints correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				taskDefinition := fmt.Sprintf("%ssome_name:1", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(family, taskRoleARN, network string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume, placementConstraints []*ecs.PlacementConstraint) {
						reporter.AssertEqual(network, "host")
						reporter.AssertEqual(taskRoleARN, "some_role")
						reporter.AssertEqual(len(containers), 1)
						reporter.AssertEqual(*containers[0].Name, "test")
						reporter.AssertEqual(*volumes[0].Name, "test")
						reporter.AssertEqual(len(placementConstraints), 0)
						reporter.AssertEqual(*containers[0].Image, "12345.dkr.ecr.us-west-2.amazonaws.com/test")
					}).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrun)
			},
		},
		{
			Name: "Should marshal dockerrun with placement constraints correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				taskDefinition := fmt.Sprintf("%ssome_name:1", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(family, taskRoleARN, network string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume, placementConstraints []*ecs.PlacementConstraint) {
						reporter.AssertEqual(network, "host")
						reporter.AssertEqual(taskRoleARN, "some_role")
						reporter.AssertEqual(len(containers), 1)
						reporter.AssertEqual(*containers[0].Name, "test")
						reporter.AssertEqual(*volumes[0].Name, "test")
						reporter.AssertEqual(len(placementConstraints), 1)
						reporter.AssertEqual(*containers[0].Image, "12345.dkr.ecr.us-west-2.amazonaws.com/test")
					}).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrunWithPCs)
			},
		},
		{
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
					RegisterTaskDefinition(deployID.String(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
		{
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
		{
			Name: "Should propagate ecs.RegisterTaskDefinition error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockDeploy := NewMockECSDeployManager(ctrl)

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
		{
			Name: "Image prefix should remain same if docker repo override is NOT specified",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				os.Unsetenv(config.DOCKER_REPO_OVERRIDE)
				mockDeploy := NewMockECSDeployManager(ctrl)

				taskDefinition := fmt.Sprintf("%ssome_name:1", id.PREFIX)

				task := &ecs.TaskDefinition{
					&aws_ecs.TaskDefinition{
						TaskDefinitionArn: stringp("arn:aws:ecs:us-west-2:12345678:task-definition/" + taskDefinition),
					},
				}

				mockDeploy.ECS.EXPECT().
					RegisterTaskDefinition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(family, taskRoleARN, network string, containers []*ecs.ContainerDefinition, volumes []*ecs.Volume, placementConstraints []*ecs.PlacementConstraint) {
						reporter.AssertEqual(*containers[0].Image, "quintilesims/test")
					}).
					Return(task, nil)

				return mockDeploy.Deploy()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSDeployManager)
				manager.CreateDeploy("some_name", dockerrun)
			},
		},
	}

	testutils.RunTests(t, testCases)
}
