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
)

func TestEnvironmentUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	tags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "l0-test-env_id",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "type",
			Value:      "static",
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

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	req := models.UpdateEnvironmentRequest{
		Scale: aws.Int(2),
	}

	// ensure we update the asg's max size as well since it is greater than the current max
	updateASGInput := &autoscaling.UpdateAutoScalingGroupInput{}
	updateASGInput.SetAutoScalingGroupName("l0-test-env_id")
	updateASGInput.SetMinSize(2)
	updateASGInput.SetMaxSize(2)

	mockAWS.AutoScaling.EXPECT().
		UpdateAutoScalingGroup(updateASGInput).
		Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update("env_id", req); err != nil {
		t.Fatal(err)
	}
}
