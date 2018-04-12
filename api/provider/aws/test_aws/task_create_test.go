package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreate_stateless(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"subnet-test"})

	defer provider.SetEntityIDGenerator("tsk_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// define expected ec2.DescribeSecurityGroups
	// (part of "awsvpc" NetworkMode workflow)
	ec2Filter := &ec2.Filter{}
	ec2Filter.SetName("group-name")
	ec2Filter.SetValues([]*string{aws.String("l0-test-env_id-env")})

	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{}
	describeSecurityGroupsInput.SetFilters([]*ec2.Filter{ec2Filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-env_id-env")
	securityGroup.SetGroupId("sg-test")

	securityGroups := []*ec2.SecurityGroup{securityGroup}

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups(securityGroups)

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(describeSecurityGroupsInput).
		Return(describeSecurityGroupsOutput, nil)

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	networkMode := ecs.NetworkModeAwsvpc

	taskDefinition := &ecs.TaskDefinition{
		Compatibilities:      []*string{aws.String(ecs.LaunchTypeEc2), aws.String(ecs.LaunchTypeFargate)},
		ContainerDefinitions: containerDefinitions,
		NetworkMode:          &networkMode,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected RunTask
	kvp := &ecs.KeyValuePair{}
	kvp.SetName("key")
	kvp.SetValue("val")

	override := &ecs.ContainerOverride{}
	override.SetEnvironment([]*ecs.KeyValuePair{kvp})
	override.SetName("container")

	taskOverride := &ecs.TaskOverride{}
	taskOverride.SetContainerOverrides([]*ecs.ContainerOverride{override})

	awsvpcConfig := &ecs.AwsVpcConfiguration{}
	awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
	awsvpcConfig.SetSecurityGroups([]*string{aws.String("sg-test")})
	awsvpcConfig.SetSubnets([]*string{aws.String("subnet-test")})

	networkConfig := &ecs.NetworkConfiguration{}
	networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

	runTaskInput := &ecs.RunTaskInput{}
	runTaskInput.SetCluster("l0-test-env_id")
	runTaskInput.SetLaunchType(ecs.LaunchTypeFargate)
	runTaskInput.SetNetworkConfiguration(networkConfig)
	runTaskInput.SetOverrides(taskOverride)
	runTaskInput.SetPlatformVersion(config.DefaultFargatePlatformVersion)
	runTaskInput.SetStartedBy("test")
	runTaskInput.SetTaskDefinition("dpl_arn")

	task := &ecs.Task{}
	task.SetTaskArn("arn:aws:ecs:region:012345678910:task/arn")

	runTaskOutput := &ecs.RunTaskOutput{}
	runTaskOutput.SetTasks([]*ecs.Task{task})

	mockAWS.ECS.EXPECT().
		RunTask(runTaskInput).
		Return(runTaskOutput, nil)

	// define request
	containerOverride := []models.ContainerOverride{
		{
			ContainerName:        "container",
			EnvironmentOverrides: map[string]string{"key": "val"},
		},
	}

	req := models.CreateTaskRequest{
		ContainerOverrides: containerOverride,
		DeployID:           "dpl_id",
		EnvironmentID:      "env_id",
		TaskName:           "tsk_name",
		Stateful:           false,
	}

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
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

func TestTaskCreate_stateful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	defer provider.SetEntityIDGenerator("tsk_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	taskDefinition := &ecs.TaskDefinition{
		ContainerDefinitions: containerDefinitions,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected RunTask
	kvp := &ecs.KeyValuePair{}
	kvp.SetName("key")
	kvp.SetValue("val")

	override := &ecs.ContainerOverride{}
	override.SetEnvironment([]*ecs.KeyValuePair{kvp})
	override.SetName("container")

	taskOverride := &ecs.TaskOverride{}
	taskOverride.SetContainerOverrides([]*ecs.ContainerOverride{override})

	awsvpcConfig := &ecs.AwsVpcConfiguration{}
	awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
	awsvpcConfig.SetSecurityGroups([]*string{aws.String("sg-test")})
	awsvpcConfig.SetSubnets([]*string{aws.String("subnet-test")})

	networkConfig := &ecs.NetworkConfiguration{}
	networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

	runTaskInput := &ecs.RunTaskInput{}
	runTaskInput.SetCluster("l0-test-env_id")
	runTaskInput.SetLaunchType(ecs.LaunchTypeEc2)
	runTaskInput.SetOverrides(taskOverride)
	runTaskInput.SetStartedBy("test")
	runTaskInput.SetTaskDefinition("dpl_arn")

	task := &ecs.Task{}
	task.SetTaskArn("arn:aws:ecs:region:012345678910:task/arn")

	runTaskOutput := &ecs.RunTaskOutput{}
	runTaskOutput.SetTasks([]*ecs.Task{task})

	mockAWS.ECS.EXPECT().
		RunTask(runTaskInput).
		Return(runTaskOutput, nil)

	// define request
	containerOverride := []models.ContainerOverride{
		{
			ContainerName:        "container",
			EnvironmentOverrides: map[string]string{"key": "val"},
		},
	}

	req := models.CreateTaskRequest{
		ContainerOverrides: containerOverride,
		DeployID:           "dpl_id",
		EnvironmentID:      "env_id",
		TaskName:           "tsk_name",
		Stateful:           true,
	}

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
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
