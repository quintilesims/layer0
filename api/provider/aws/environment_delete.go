package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func (e *EnvironmentProvider) Delete(environmentID string) error {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	autoScalingGroupName := fqEnvironmentID
	if err := e.deleteASG(autoScalingGroupName); err != nil {
		return err
	}

	launchConfigName := fqEnvironmentID
	if err := e.deleteLC(launchConfigName); err != nil {
		return err
	}

	securityGroup, err := e.readSG(fqEnvironmentID)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return err
	}

	if securityGroup != nil {
		groupID := aws.StringValue(securityGroup.GroupId)
		if err := e.deleteSG(groupID); err != nil {
			return err
		}
	}

	clusterName := fqEnvironmentID
	if err := e.deleteCluster(clusterName); err != nil {
		return err
	}

	if err := e.deleteTags(environmentID); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) deleteASG(autoScalingGroupName string) error {
	input := &autoscaling.DeleteAutoScalingGroupInput{}
	input.SetAutoScalingGroupName(autoScalingGroupName)
	input.SetForceDelete(true)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.DeleteAutoScalingGroup(input); err != nil {
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "AutoScalingGroup name not found") {
			return nil
		}

		return err
	}

	return nil
}

func (e *EnvironmentProvider) deleteLC(launchConfigName string) error {
	input := &autoscaling.DeleteLaunchConfigurationInput{}
	input.SetLaunchConfigurationName(launchConfigName)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.DeleteLaunchConfiguration(input); err != nil {
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "Launch configuration name not found") {
			return nil
		}

		return err
	}

	return nil
}

func (e *EnvironmentProvider) deleteSG(securityGroupID string) error {
	input := &ec2.DeleteSecurityGroupInput{}
	input.SetGroupId(securityGroupID)

	if _, err := e.AWS.EC2.DeleteSecurityGroup(input); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) deleteCluster(clusterName string) error {
	input := &ecs.DeleteClusterInput{}
	input.SetCluster(clusterName)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.ECS.DeleteCluster(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ClusterNotFoundException" {
			return nil
		}

		return err
	}

	return nil
}

func (e *EnvironmentProvider) deleteTags(environmentID string) error {
	tags, err := e.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := e.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			return err
		}
	}

	return nil
}
