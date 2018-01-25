package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentDelete(t *testing.T) {
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

	// an environment's asg name is the same as the fq environment id
	deleteASGInput := &autoscaling.DeleteAutoScalingGroupInput{}
	deleteASGInput.SetAutoScalingGroupName("l0-test-env_id")
	deleteASGInput.SetForceDelete(true)

	mockAWS.AutoScaling.EXPECT().
		DeleteAutoScalingGroup(deleteASGInput).
		Return(&autoscaling.DeleteAutoScalingGroupOutput{}, nil)

	// an environment's lc name is the same as the fq environment id
	deleteLCInput := &autoscaling.DeleteLaunchConfigurationInput{}
	deleteLCInput.SetLaunchConfigurationName("l0-test-env_id")

	mockAWS.AutoScaling.EXPECT().
		DeleteLaunchConfiguration(deleteLCInput).
		Return(&autoscaling.DeleteLaunchConfigurationOutput{}, nil)

	// an environment's security group name is <fq environment id>-env
	readSGHelper(mockAWS, "l0-test-env_id-env", "sg_id")
	deleteSGHelper(mockAWS, "sg_id")

	// an environment's cluster name is the fq environment id
	deleteClusterInput := &ecs.DeleteClusterInput{}
	deleteClusterInput.SetCluster("l0-test-env_id")

	mockAWS.ECS.EXPECT().
		DeleteCluster(deleteClusterInput).
		Return(&ecs.DeleteClusterOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("env_id"); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestDeleteEnvironmentIdempotence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	mockAWS.AutoScaling.EXPECT().
		DeleteAutoScalingGroup(gomock.Any()).
		Return(nil, awserr.New("", "AutoScalingGroup name not found", nil))

	mockAWS.AutoScaling.EXPECT().
		DeleteLaunchConfiguration(gomock.Any()).
		Return(nil, awserr.New("", "Launch configuration name not found", nil))

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(gomock.Any()).
		Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

	mockAWS.ECS.EXPECT().
		DeleteCluster(gomock.Any()).
		Return(nil, awserr.New("ClusterNotFoundException", "", nil))

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("env_id"); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentDeleteDependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	envTag := models.Tag{
		EntityID:   "env_id",
		EntityType: "environment",
		Key:        "name",
		Value:      "env_name",
	}

	if err := tagStore.Insert(envTag); err != nil {
		t.Fatal(err)
	}

	var dependencies = []struct {
		entityID   string
		entityType string
		err        string // expected result
	}{
		{"lb_id", "load_balancer", error.Error(errors.Newf(errors.DependencyError, "Cannot delete environment 'env_id' because it contains dependent load balancers: lb_id"))},
		{"tsk_id", "task", error.Error(errors.Newf(errors.DependencyError, "Cannot delete environment 'env_id' because it contains dependent tasks: tsk_id"))},
		{"svc_id", "service", error.Error(errors.Newf(errors.DependencyError, "Cannot delete environment 'env_id' because it contains dependent services: svc_id"))},
	}

	for _, dependency := range dependencies {
		tag := models.Tag{
			EntityID:   dependency.entityID,
			EntityType: dependency.entityType,
			Key:        "environment_id",
			Value:      "env_id",
		}

		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}

		target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
		actual := target.Delete("env_id")

		if actual.Error() != dependency.err {
			t.Errorf("\nexpected %s\nactual %s", dependency.err, actual)
		}

		if err := tagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			t.Fatal(err)
		}
	}
}
