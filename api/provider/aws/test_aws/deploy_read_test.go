package test_aws

import (
	"encoding/json"
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
			Value:      "arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// Set Container Definition Defaults
	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("taskDefinition")
	containerDefinitions := []*ecs.ContainerDefinition{containerDefinition}

	// Set Task Definition Defaults
	taskDefinition := &ecs.TaskDefinition{}
	taskDefinition.SetContainerDefinitions(containerDefinitions)

	// Set up ECS mock inputs and outputs
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1")

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	deployFile, err := json.Marshal(taskDefinitionOutput.TaskDefinition)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Deploy{
		DeployFile: deployFile,
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		Version:    "dpl_version",
	}

	assert.Equal(t, expected, result)
}
