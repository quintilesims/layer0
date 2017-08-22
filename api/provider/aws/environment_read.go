package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/quintilesims/layer0/common/models"
)

// todo: catch 'EntityDoesNotExist' errors
func (e *EnvironmentProvider) Read(environmentID string) (*models.Environment, error) {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

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

	launchConfigName := aws.StringValue(autoScalingGroup.LaunchConfigurationName)
	launchConfig, err := e.readLC(launchConfigName)
	if err != nil {
		return nil, err
	}

	model := &models.Environment{
		EnvironmentID:   environmentID,
		ClusterCount:    len(autoScalingGroup.Instances),
		InstanceSize:    aws.StringValue(launchConfig.InstanceType),
		SecurityGroupID: aws.StringValue(securityGroup.GroupId),
		AMIID:           aws.StringValue(launchConfig.ImageId),
	}

	if err := e.populateModelTagss(environmentID, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (e *EnvironmentProvider) readLC(launchConfigName string) (*autoscaling.LaunchConfiguration, error) {
	input := &autoscaling.DescribeLaunchConfigurationsInput{}
	input.SetLaunchConfigurationNames([]*string{aws.String(launchConfigName)})

	output, err := e.AWS.AutoScaling.DescribeLaunchConfigurations(input)
	if err != nil {
		return nil, err
	}

	for _, lc := range output.LaunchConfigurations {
		if aws.StringValue(lc.LaunchConfigurationName) == launchConfigName {
			return lc, nil
		}
	}

	return nil, fmt.Errorf("Launch Configuration '%s' does not exist", launchConfigName)
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

	return nil, fmt.Errorf("AutoScaling Group '%s' does not exist", autoScalingGroupName)
}

func (e *EnvironmentProvider) populateModelTagss(environmentID string, model *models.Environment) error {
	tags, err := e.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.EnvironmentName = tag.Value
	}

	if tag, ok := tags.WithKey("os").First(); ok {
		model.OperatingSystem = tag.Value
	}

	model.Links = []string{}
	for _, tag := range tags.WithKey("link") {
		model.Links = append(model.Links, tag.Value)
	}

	return nil
}
