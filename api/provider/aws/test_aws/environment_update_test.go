package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

func TestEnvironmentUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := config.NewTestContext(t, nil, map[string]interface{}{
		config.FlagInstance.GetName(): "test",
	})

	req := models.UpdateEnvironmentRequest{
		MinScale: aws.Int(2),
		MaxScale: aws.Int(5),
	}

	// an environment's asg name is the same as the fq environment id
	describeASGInput := &autoscaling.DescribeAutoScalingGroupsInput{}
	describeASGInput.SetAutoScalingGroupNames([]*string{aws.String("l0-test-env_id")})

	asg := &autoscaling.Group{}
	asg.SetAutoScalingGroupName("l0-test-env_id")

	describeASGOutput := &autoscaling.DescribeAutoScalingGroupsOutput{}
	describeASGOutput.SetAutoScalingGroups([]*autoscaling.Group{asg})

	mockAWS.AutoScaling.EXPECT().
		DescribeAutoScalingGroups(describeASGInput).
		Return(describeASGOutput, nil)

	// ensure we update the asg's max size as well since it is greater than the current max
	updateASGInput := &autoscaling.UpdateAutoScalingGroupInput{}
	updateASGInput.SetAutoScalingGroupName("l0-test-env_id")
	updateASGInput.SetMinSize(2)
	updateASGInput.SetMaxSize(5)

	mockAWS.AutoScaling.EXPECT().
		UpdateAutoScalingGroup(updateASGInput).
		Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, c)
	if err := target.Update("env_id", req); err != nil {
		t.Fatal(err)
	}
}
