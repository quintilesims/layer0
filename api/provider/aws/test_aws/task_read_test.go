package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskRead(t *testing.T) {
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

	// environments' asg names are the same as the fq environment id
	describeASGInput := &autoscaling.DescribeAutoScalingGroupsInput{}
	describeASGInput.SetAutoScalingGroupNames([]*string{aws.String("l0-test-env_id")})

	// environments use the asg.LaunchConfigurationName as the source of truth for
	// their launch configuration name
	asg := &autoscaling.Group{}
	asg.SetAutoScalingGroupName("l0-test-env_id")
	asg.SetLaunchConfigurationName("lc_name")
	asg.SetInstances(make([]*autoscaling.Instance, 2))

	describeASGOutput := &autoscaling.DescribeAutoScalingGroupsOutput{}
	describeASGOutput.SetAutoScalingGroups([]*autoscaling.Group{asg})

	mockAWS.AutoScaling.EXPECT().
		DescribeAutoScalingGroups(describeASGInput).
		Return(describeASGOutput, nil)

	describeLCInput := &autoscaling.DescribeLaunchConfigurationsInput{}
	describeLCInput.SetLaunchConfigurationNames([]*string{aws.String("lc_name")})

	lc := &autoscaling.LaunchConfiguration{}
	lc.SetLaunchConfigurationName("lc_name")
	lc.SetInstanceType("m3.small")
	lc.SetImageId("ami_id")

	describeLCOutput := &autoscaling.DescribeLaunchConfigurationsOutput{}
	describeLCOutput.SetLaunchConfigurations([]*autoscaling.LaunchConfiguration{lc})

	taskARNs := []*string{
		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id1"),
		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-tsk_id2"),
		aws.String("arn:aws:ecs:region:012345678910:task/l0-bad-tsk_id1"),
		aws.String("arn:aws:ecs:region:012345678910:task/bad2"),
	}

	input := &ecs.DescribeTasksInput{}
	// input.SetCluster()
	input.SetTasks(taskARNs)

	output := &ecs.DescribeTasksOutput{}
	// output.SetTasks()

	mockAWS.ECS.EXPECT().
		DescribeTasks(input).
		Return(output, nil)

	target := provider.NewTaskProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("tsk_id")
	if err != nil {
		t.Fatal(err)
	}

	// todo: Links
	expected := &models.Task{}

	assert.Equal(t, expected, result)
}
