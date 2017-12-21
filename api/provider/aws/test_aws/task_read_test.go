package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := config.NewTestContext(t, nil, map[string]interface{}{
		config.FlagInstance.GetName(): "test",
	})

	tags := models.Tags{
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl_name",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "deployVersion",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeTaskInput := &ecs.DescribeTasksInput{}
	describeTaskInput.SetCluster("l0-test-env_id")
	describeTaskInput.SetTasks([]*string{aws.String("arn:aws:ecs:region:012345678910:task/arn")})

	containerECS := &ecs.Container{}
	containerECS.SetName("container")
	containerECS.SetLastStatus("status")
	containerECS.SetExitCode(1)

	task := &ecs.Task{}
	task.SetTaskArn("arn:aws:ecs:region:012345678910:task/arn")
	task.SetTaskDefinitionArn("arn:aws:ecs:region:account:task-definition/dpl_id:deployVersion")
	task.SetLastStatus(ecs.DesiredStatusRunning)
	task.SetContainers([]*ecs.Container{containerECS})

	describeTaskOutput := &ecs.DescribeTasksOutput{}
	describeTaskOutput.SetTasks([]*ecs.Task{task})

	mockAWS.ECS.EXPECT().
		DescribeTasks(describeTaskInput).
		Return(describeTaskOutput, nil)

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, c)
	result, err := target.Read("tsk_id")
	if err != nil {
		t.Fatal(err)
	}

	container := models.Container{
		ContainerName: "container",
		Status:        "status",
		ExitCode:      1,
	}

	expected := &models.Task{
		TaskID:        "tsk_id",
		TaskName:      "tsk_name",
		EnvironmentID: "env_id",
		DeployID:      "dpl_id",
		DeployName:    "dpl_name",
		DeployVersion: "deployVersion",
		Status:        "RUNNING",
		Containers:    []models.Container{container},
	}

	assert.Equal(t, expected, result)
}
