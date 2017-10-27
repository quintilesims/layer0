package test_aws

import (
	"testing"

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

	describeSecurityGroupHelper(mockAWS, "l0-test-env_id-env", "sg_id")

	asg := &autoscaling.Group{}
	asg.SetAutoScalingGroupName("l0-test-env_id")
	asg.SetLaunchConfigurationName("l0-test-env_id")

	describeASGOutput := &autoscaling.DescribeAutoScalingGroupsOutput{}
	describeASGOutput.SetAutoScalingGroups([]*autoscaling.Group{asg})

	// todo: check input
	mockAWS.AutoScaling.EXPECT().
		DescribeAutoScalingGroups(gomock.Any()).
		Return(describeASGOutput, nil)

	lc := &autoscaling.LaunchConfiguration{}
	lc.SetLaunchConfigurationName("l0-test-env_id")

	describeLCOutput := &autoscaling.DescribeLaunchConfigurationsOutput{}
	describeLCOutput.SetLaunchConfigurations([]*autoscaling.LaunchConfiguration{lc})

	// todo: check input
	mockAWS.AutoScaling.EXPECT().
		DescribeLaunchConfigurations(gomock.Any()).
		Return(describeLCOutput, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("env_id")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		OperatingSystem: "linux",
		InstanceSize:    "m3.small",
		SecurityGroupID: "sg_id",
		ClusterCount:    2,
		AMIID:           "ami_id",
		// TODO: Links
	}

	assert.Equal(t, expected, result)
}
