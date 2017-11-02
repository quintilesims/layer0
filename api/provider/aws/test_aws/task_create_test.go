package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name1",
		},
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn1",
		},
		{
			EntityID:   "tsk_id1",
			EntityType: "task",
			Key:        "id",
			Value:      "env_id1",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk_name2",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "arn",
			Value:      "arn2",
		},
		{
			EntityID:   "tsk_id2",
			EntityType: "task",
			Key:        "id",
			Value:      "env_id2",
		},
		{
			EntityID:   "env_id1",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name1",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// ensure we filter out clusters that don't belong to our instance
	// TODo: Make Helper?
	listClusterPagesFN := func(input *ecs.ListClustersInput, fn func(output *ecs.ListClustersOutput, lastPage bool) bool) error {
		clusterARNs := []*string{
			aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id1"),
			aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id2"),
			aws.String("arn:aws:ecs:region:012345678910:cluster/l0-bad-env_id1"),
			aws.String("arn:aws:ecs:region:012345678910:cluster/bad2"),
		}

		output := &ecs.ListClustersOutput{}
		output.SetClusterArns(clusterARNs)
		fn(output, true)

		return nil
	}

	mockAWS.ECS.EXPECT().
		ListClustersPages(gomock.Any(), gomock.Any()).
		Do(listClusterPagesFN).
		Return(nil)

	// ensure we filter out tasks that don't belong to our instance
	listTaskPagesFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
		taskARNs := []*string{
			aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id1"),
			aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id2"),
			aws.String("arn:aws:ecs:region:012345678910:task/l0-bad-tsk_id1"),
			aws.String("arn:aws:ecs:region:012345678910:task/bad2"),
		}

		output := &ecs.ListTasksOutput{}
		output.SetTaskArns(taskARNs)
		fn(output, true)

		return nil
	}

	mockAWS.ECS.EXPECT().
		ListTasksPages(gomock.Any(), gomock.Any()).
		Do(listTaskPagesFN).
		Return(nil)

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.TaskSummary{
		{
			TaskID:          "tsk_id1",
			TaskName:        "tsk_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			TaskID:          "tsk_id2",
			TaskName:        "tsk_name2",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}

	assert.Equal(t, expected, result)
}