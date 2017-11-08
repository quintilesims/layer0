package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeployRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

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
			Value:      "dpl_version",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("arn:aws:ecs:region:012345678910:task/arn")

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Deploy{
		DeployFile: []byte("null"),
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		Version:    "dpl_version",
	}

	assert.Equal(t, expected, result)
}
