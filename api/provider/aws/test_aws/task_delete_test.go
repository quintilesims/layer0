package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

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
			Value:      "arn1",
		},
		{
			EntityID:   "tsk_id",
			EntityType: "task",
			Key:        "id",
			Value:      "env_id",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// // ensure we filter out tasks that don't belong to our instance
	// listTaskPagesFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
	// 	taskARNs := []*string{
	// 		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id1"),
	// 		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id2"),
	// 		aws.String("arn:aws:ecs:region:012345678910:task/l0-bad-tsk_id1"),
	// 		aws.String("arn:aws:ecs:region:012345678910:task/bad2"),
	// 	}

	// 	output := &ecs.ListTasksOutput{}
	// 	output.SetTaskArns(taskARNs)
	// 	fn(output, true)

	// 	return nil
	// }

	stopTaskInput := &ecs.StopTaskInput{}
	stopTaskInput.SetCluster("l0-test-env_id")
	stopTaskInput.SetTask("l0-test-tsk_id")

	mockAWS.ECS.EXPECT().
		StopTask(stopTaskInput).
		Return(nil, nil)

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
	err := target.Delete("tsk_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestDeleteTaskIdempotence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	mockAWS.ECS.EXPECT().
		StopTask(gomock.Any()).
		Return(nil, awserr.New("TaskDoesNotExist", "", nil))

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("tsk_id"); err != nil {
		t.Fatal(err)
	}
}
