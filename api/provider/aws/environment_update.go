package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Update(environmentID string, minSize int64) (*models.Environment, error) {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	autoScalingGroupName := fqEnvironmentID
	asg, err := e.readASG(autoScalingGroupName)
	if err != nil {
		return nil, err
	}

	maxSize := aws.Int64Value(asg.MaxSize)
	if maxSize < minSize {
		maxSize = minSize
	}

	if err := e.updateASGSize(autoScalingGroupName, minSize, maxSize); err != nil {
		return nil, err
	}

	return e.Read(environmentID)
}

func (e *EnvironmentProvider) updateASGSize(autoScalingGroupName string, minSize, maxSize int64) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{}
	input.SetAutoScalingGroupName(autoScalingGroupName)
	input.SetMinSize(minSize)
	input.SetMaxSize(maxSize)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.UpdateAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}
