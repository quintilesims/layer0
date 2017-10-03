package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/quintilesims/layer0/common/models"
)

// Update updates an ECS Cluster using the specified Update Environment Request.
// The Update Environment Request contains the Environment ID and the minimum size
// of the Cluster's Auto Scaling Group. The Cluster's Auto Scaling Group size is
// updated by making an UpdateAutoScalingGroup request to AWS.
func (e *EnvironmentProvider) Update(req models.UpdateEnvironmentRequest) error {
	environmentID := req.EnvironmentID
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	if req.MinClusterCount != nil {
		minSize := int64(*req.MinClusterCount)

		autoScalingGroupName := fqEnvironmentID
		asg, err := e.readASG(autoScalingGroupName)
		if err != nil {
			return err
		}

		maxSize := aws.Int64Value(asg.MaxSize)
		if maxSize < minSize {
			maxSize = minSize
		}

		if err := e.updateASGSize(autoScalingGroupName, minSize, maxSize); err != nil {
			return err
		}
	}

	return nil
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
