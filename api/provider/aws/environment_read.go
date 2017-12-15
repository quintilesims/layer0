package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a *models.Environment based on the provided environmentID. The environmentID
// is used to look up the Environment's Security Group, Auto Scaling Group, and
// Launch Configuration when DescribeSecurityGroups, DescribeAutoScalingGroups, and
// DescribeLaunchConfigurations requests respectively are made to AWS.
func (e *EnvironmentProvider) Read(environmentID string) (*models.Environment, error) {
	// todo: catch 'EntityDoesNotExist' errors
	fqEnvironmentID := addLayer0Prefix(e.Context, environmentID)

	securityGroupName := getEnvironmentSGName(fqEnvironmentID)
	securityGroup, err := readSG(e.AWS.EC2, securityGroupName)
	if err != nil {
		return nil, err
	}

	autoScalingGroupName := fqEnvironmentID
	autoScalingGroup, err := e.readASG(autoScalingGroupName)
	if err != nil {
		return nil, err
	}

	launchContextName := aws.StringValue(autoScalingGroup.LaunchConfigurationName)
	launchContext, err := e.readLC(launchContextName)
	if err != nil {
		return nil, err
	}

	model, err := e.makeEnvironmentModel(environmentID)
	if err != nil {
		return nil, err
	}

	model.MinScale = int(aws.Int64Value(autoScalingGroup.MinSize))
	model.CurrentScale = int(aws.Int64Value(autoScalingGroup.DesiredCapacity))
	model.MaxScale = int(aws.Int64Value(autoScalingGroup.MaxSize))
	model.InstanceType = aws.StringValue(launchContext.InstanceType)
	model.SecurityGroupID = aws.StringValue(securityGroup.GroupId)
	model.AMIID = aws.StringValue(launchContext.ImageId)

	return model, nil
}

func (e *EnvironmentProvider) readLC(launchContextName string) (*autoscaling.LaunchConfiguration, error) {
	input := &autoscaling.DescribeLaunchConfigurationsInput{}
	input.SetLaunchConfigurationNames([]*string{aws.String(launchContextName)})

	output, err := e.AWS.AutoScaling.DescribeLaunchConfigurations(input)
	if err != nil {
		return nil, err
	}

	for _, lc := range output.LaunchConfigurations {
		if aws.StringValue(lc.LaunchConfigurationName) == launchContextName {
			return lc, nil
		}
	}

	message := fmt.Sprintf("Launch Configuration '%s' does not exist", launchContextName)
	return nil, awserr.New("DoesNotExist", message, nil)
}

func (e *EnvironmentProvider) readASG(autoScalingGroupName string) (*autoscaling.Group, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	input.SetAutoScalingGroupNames([]*string{aws.String(autoScalingGroupName)})

	output, err := e.AWS.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}

	for _, asg := range output.AutoScalingGroups {
		if aws.StringValue(asg.AutoScalingGroupName) == autoScalingGroupName {
			return asg, nil
		}
	}

	message := fmt.Sprintf("AutoScalingGroup '%s' does not exist", autoScalingGroupName)
	return nil, awserr.New("DoesNotExist", message, nil)
}

func (e *EnvironmentProvider) makeEnvironmentModel(environmentID string) (*models.Environment, error) {
	model := &models.Environment{
		EnvironmentID: environmentID,
	}

	tags, err := e.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.EnvironmentName = tag.Value
	}

	if tag, ok := tags.WithKey("os").First(); ok {
		model.OperatingSystem = tag.Value
	}

	model.Links = []string{}
	for _, tag := range tags.WithKey("link") {
		model.Links = strings.Split(tag.Value, ",")
	}

	return model, nil
}
