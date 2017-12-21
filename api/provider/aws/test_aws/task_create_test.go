package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := config.NewTestContext(t, nil, map[string]interface{}{
		config.FlagInstance.GetName(): "test",
	})

	defer provider.SetEntityIDGenerator("tsk_id")()

	tags := models.Tags{
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
			Value:      "version",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	containerOverride := []models.ContainerOverride{
		{
			ContainerName:        "container",
			EnvironmentOverrides: map[string]string{"key": "val"},
		},
	}

	req := models.CreateTaskRequest{
		DeployID:           "dpl_id",
		EnvironmentID:      "env_id",
		TaskName:           "tsk_name",
		ContainerOverrides: containerOverride,
	}

	kvp := &ecs.KeyValuePair{}
	kvp.SetName("key")
	kvp.SetValue("val")

	override := &ecs.ContainerOverride{}
	override.SetName("container")
	override.SetEnvironment([]*ecs.KeyValuePair{kvp})

	taskOverride := &ecs.TaskOverride{}
	taskOverride.SetContainerOverrides([]*ecs.ContainerOverride{override})

	runTaskInput := &ecs.RunTaskInput{}
	runTaskInput.SetCluster("l0-test-env_id")
	runTaskInput.SetStartedBy("test")
	runTaskInput.SetTaskDefinition("l0-test-dpl_name:version")
	runTaskInput.SetOverrides(taskOverride)

	task := &ecs.Task{}
	task.SetTaskArn("arn:aws:ecs:region:012345678910:task/arn")

	runTaskOutput := &ecs.RunTaskOutput{}
	runTaskOutput.SetTasks([]*ecs.Task{task})

	mockAWS.ECS.EXPECT().
		RunTask(runTaskInput).
		Return(runTaskOutput, nil)

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, c)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "tsk_id", result)

	expectedTags := models.Tags{
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
