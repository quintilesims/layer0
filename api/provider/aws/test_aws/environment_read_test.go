package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "linux",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// an environment's security group name is <fq environment id>-env
	readSGHelper(mockAWS, "l0-test-env_id-env", "sg_id")

	// an environment's asg name is the same as the fq environment id
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

	mockAWS.AutoScaling.EXPECT().
		DescribeLaunchConfigurations(describeLCInput).
		Return(describeLCOutput, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("env_id")
	if err != nil {
		t.Fatal(err)
	}

	// todo: Links
	expected := &models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		OperatingSystem: "linux",
		InstanceSize:    "m3.small",
		SecurityGroupID: "sg_id",
		ClusterCount:    2,
		AMIID:           "ami_id",
		Links:           []string{},
	}

	assert.Equal(t, expected, result)
}
