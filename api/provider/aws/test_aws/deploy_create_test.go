package test_aws

import (
	"encoding/json"
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

func TestDeployCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	defer provider.SetEntityIDGenerator("dpl_id")()

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

	// define container defaults
	container := &ecs.ContainerDefinition{}
	container.SetName("test_name")
	logConfig := &ecs.LogConfiguration{
		LogDriver: aws.String("awslogs"),
		Options: map[string]*string{
			"awslogs-group":         aws.String("test_group"),
			"awslogs-region":        aws.String("test_region"),
			"awslogs-stream-prefix": aws.String("test_prefix"),
		},
	}
	container.SetLogConfiguration(logConfig)
	containers := []*ecs.ContainerDefinition{}
	containers = append(containers, container)

	// define request
	reqDeployFile := &ecs.TaskDefinition{}
	reqDeployFile.SetContainerDefinitions(containers)
	deployFile, _ := json.Marshal(reqDeployFile)

	req := models.CreateDeployRequest{
		DeployName: "dpl_id",
		DeployFile: deployFile,
	}

	registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{}
	registerTaskDefinitionInput.SetFamily("l0-test-dpl_id")
	registerTaskDefinitionInput.SetTaskRoleArn("")
	registerTaskDefinitionInput.SetContainerDefinitions(containers)
	registerTaskDefinitionInput.SetVolumes(nil)
	registerTaskDefinitionInput.SetPlacementConstraints(nil)

	taskDefinitionOutput := &ecs.TaskDefinition{}
	registerTaskDefinitionOutput := &ecs.RegisterTaskDefinitionOutput{}
	registerTaskDefinitionOutput.SetTaskDefinition(taskDefinitionOutput)

	mockAWS.ECS.EXPECT().
		RegisterTaskDefinition(registerTaskDefinitionInput).
		Return(registerTaskDefinitionOutput, nil)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", result)

	expectedTags := models.Tags{
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

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
