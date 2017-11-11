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

	// define container defaults
	logConfig := &ecs.LogConfiguration{
		LogDriver: aws.String("awslogs"),
		Options: map[string]*string{
			"awslogs-group":         aws.String("test_group"),
			"awslogs-region":        aws.String("test_region"),
			"awslogs-stream-prefix": aws.String("test_prefix"),
		},
	}

	cntr1 := &ecs.ContainerDefinition{}
	cntr1.SetName("cntr_name_1")
	cntr1.SetLogConfiguration(logConfig)

	cntr2 := &ecs.ContainerDefinition{}
	cntr2.SetName("cntr_name_2")
	cntr2.SetLogConfiguration(logConfig)

	containers := []*ecs.ContainerDefinition{cntr1, cntr2}

	// define request
	reqDeployFile := &ecs.TaskDefinition{}
	reqDeployFile.SetContainerDefinitions(containers)
	deployFile, err := json.Marshal(reqDeployFile)
	if err != nil {
		t.Fatal(err)
	}

	req := models.CreateDeployRequest{
		DeployName: "dpl_id",
		DeployFile: deployFile,
	}

	registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{}
	registerTaskDefinitionInput.SetFamily("l0-test-dpl_id")
	registerTaskDefinitionInput.SetContainerDefinitions(containers)
	registerTaskDefinitionInput.SetTaskRoleArn("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1")

	taskDefinitionOutput := &ecs.TaskDefinition{}
	taskDefinitionOutput.SetFamily("l0-test-dpl_id")
	taskDefinitionOutput.SetTaskRoleArn("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1")
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
			Value:      "arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
