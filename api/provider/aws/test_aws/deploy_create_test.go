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

func TestDeployCreate_defaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test_instance").AnyTimes()
	mockConfig.EXPECT().LogGroupName().Return("test_group").AnyTimes()
	mockConfig.EXPECT().Region().Return("test_region").AnyTimes()
	mockConfig.EXPECT().AccountID().Return("012345678910").AnyTimes()
	defer provider.SetEntityIDGenerator("dpl_id")()

	// define input for RegisterTaskDefinition()
	inputLogConfig := &ecs.LogConfiguration{}
	inputLogConfig.SetLogDriver("awslogs")
	inputLogConfig.SetOptions(map[string]*string{
		"awslogs-group":         aws.String("test_group"),
		"awslogs-region":        aws.String("test_region"),
		"awslogs-stream-prefix": aws.String("l0"),
	})

	inputContainer := &ecs.ContainerDefinition{}
	inputContainer.SetLogConfiguration(inputLogConfig)
	inputContainer.SetName("test_name")

	inputContainers := []*ecs.ContainerDefinition{inputContainer}

	registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{}
	registerTaskDefinitionInput.SetContainerDefinitions(inputContainers)
	registerTaskDefinitionInput.SetExecutionRoleArn("arn:aws:iam::012345678910:role/ecsTaskExecutionRole")
	registerTaskDefinitionInput.SetFamily("l0-test_instance-dpl_name")
	registerTaskDefinitionInput.SetRequiresCompatibilities([]*string{aws.String("EC2"), aws.String("FARGATE")})
	registerTaskDefinitionInput.SetTaskRoleArn("")

	// define output from RegisterTaskDefinition()
	outputTaskDef := &ecs.TaskDefinition{}
	outputTaskDef.SetCompatibilities([]*string{aws.String("EC2")})
	outputTaskDef.SetContainerDefinitions(inputContainers)
	outputTaskDef.SetFamily("l0-test-dpl_name")
	outputTaskDef.SetRequiresCompatibilities([]*string{aws.String("EC2"), aws.String("FARGATE")})
	outputTaskDef.SetRevision(1)
	outputTaskDef.SetTaskDefinitionArn("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1")
	outputTaskDef.SetTaskRoleArn("")

	registerTaskDefinitionOutput := &ecs.RegisterTaskDefinitionOutput{TaskDefinition: outputTaskDef}

	mockAWS.ECS.EXPECT().
		RegisterTaskDefinition(registerTaskDefinitionInput).
		Return(registerTaskDefinitionOutput, nil)

	// define request
	reqContainer := &ecs.ContainerDefinition{}
	reqContainer.SetName("test_name")

	reqContainers := []*ecs.ContainerDefinition{}
	reqContainers = append(reqContainers, reqContainer)

	reqDeployFile := &ecs.TaskDefinition{}
	reqDeployFile.SetContainerDefinitions(reqContainers)
	deployFile, err := json.Marshal(reqDeployFile)
	if err != nil {
		t.Fatal(err)
	}

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: deployFile,
	}

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
			Value:      "1",
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

func TestDeployCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().LogGroupName().Return("l0-test").AnyTimes()
	mockConfig.EXPECT().Region().Return("us-west-2").AnyTimes()
	mockConfig.EXPECT().AccountID().Return("012345678910").AnyTimes()
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

	logConfig2 := &ecs.LogConfiguration{
		LogDriver: aws.String("awslogs"),
		Options: map[string]*string{
			"awslogs-group":         aws.String("l0-test"),
			"awslogs-region":        aws.String("us-west-2"),
			"awslogs-stream-prefix": aws.String("l0"),
		},
	}

	cntr1 := &ecs.ContainerDefinition{}
	cntr1.SetName("cntr_name_1")
	cntr1.SetLogConfiguration(logConfig)

	cntr2 := &ecs.ContainerDefinition{}
	cntr2.SetName("cntr_name_2")

	containers := []*ecs.ContainerDefinition{cntr1, cntr2}

	compatibilities := []*string{aws.String("EC2")}

	// define request
	reqDeployFile := &ecs.TaskDefinition{}
	reqDeployFile.SetCpu("256")
	reqDeployFile.SetContainerDefinitions(containers)
	reqDeployFile.SetMemory("1GB")
	reqDeployFile.SetRequiresCompatibilities(compatibilities)
	reqDeployFile.SetTaskRoleArn("arn:aws:iam::012345678910:role/test-role")
	deployFile, err := json.Marshal(reqDeployFile)
	if err != nil {
		t.Fatal(err)
	}

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: deployFile,
	}

	// We set the second container's log configuration after the deployFile has been marshaled
	// to ensure that log configuration is properly rendered.
	cntr2.SetLogConfiguration(logConfig2)

	registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{}
	registerTaskDefinitionInput.SetContainerDefinitions(containers)
	registerTaskDefinitionInput.SetCpu("256")
	registerTaskDefinitionInput.SetExecutionRoleArn("arn:aws:iam::012345678910:role/ecsTaskExecutionRole")
	registerTaskDefinitionInput.SetFamily("l0-test-dpl_name")
	registerTaskDefinitionInput.SetMemory("1GB")
	registerTaskDefinitionInput.SetRequiresCompatibilities(compatibilities)
	registerTaskDefinitionInput.SetTaskRoleArn("arn:aws:iam::012345678910:role/test-role")

	taskDefinitionOutput := &ecs.TaskDefinition{}
	taskDefinitionOutput.SetCompatibilities(compatibilities)
	taskDefinitionOutput.SetContainerDefinitions(containers)
	taskDefinitionOutput.SetCpu("256")
	taskDefinitionOutput.SetFamily("l0-test-dpl_name")
	taskDefinitionOutput.SetMemory("1GB")
	taskDefinitionOutput.SetRevision(1)
	taskDefinitionOutput.SetRequiresCompatibilities(compatibilities)
	taskDefinitionOutput.SetTaskRoleArn("arn:aws:iam::012345678910:role/test-role")
	taskDefinitionOutput.SetTaskDefinitionArn("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id:1")

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
			Value:      "1",
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
