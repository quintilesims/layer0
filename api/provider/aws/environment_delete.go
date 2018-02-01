package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
)

// Delete is used to delete an ECS Cluster using the specified environmentID. The environmentID
// is used as the Environment's Auto Scaling Group, and Launch Configuration names when
// DeleteAutoScalingGroup and DeleteLaunchConfiguration requests are made to AWS, respectively.
// The environmentID is also used to look up the Environment's Security Group and the Security Group name
// is subsequently used when the DeleteSecurityGroup request is made to AWS. The ECS Cluster is deleted
// by making a DeleteCluster request to AWS.
func (e *EnvironmentProvider) Delete(environmentID string) error {
	if err := e.checkEnvironmentDependencies(environmentID); err != nil {
		return err
	}

	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	autoScalingGroupName := fqEnvironmentID
	if err := e.deleteASG(autoScalingGroupName); err != nil {
		return err
	}

	launchConfigName := fqEnvironmentID
	if err := e.deleteLC(launchConfigName); err != nil {
		return err
	}

	securityGroupName := getEnvironmentSGName(fqEnvironmentID)
	securityGroup, err := readSG(e.AWS.EC2, securityGroupName)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return err
	}

	if securityGroup != nil {
		groupID := aws.StringValue(securityGroup.GroupId)
		if err := deleteSG(e.AWS.EC2, groupID); err != nil {
			return err
		}
	}

	clusterName := fqEnvironmentID
	if err := e.deleteCluster(clusterName); err != nil {
		return err
	}

	if err := deleteEntityTags(e.TagStore, "environment", environmentID); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) checkEnvironmentDependencies(environmentID string) error {
	checkDependentEntityTags := func(entityType string) error {
		tags, err := e.TagStore.SelectByType(entityType)
		if err != nil {
			return err
		}

		dependentTags := tags.WithKey("environment_id").WithValue(environmentID)
		if len(dependentTags) > 0 {
			entityIDs := []string{}
			for _, tag := range dependentTags {
				entityIDs = append(entityIDs, tag.EntityID)
			}

			msg := fmt.Sprintf("Cannot delete environment '%s' because it contains dependent %ss:", environmentID, entityType)
			msg += strings.Join(entityIDs, ", ")
			return errors.Newf(errors.DependencyError, msg)
		}
		return nil
	}

	if err := checkDependentEntityTags("load_balancer"); err != nil {
		return err
	}

	if err := checkDependentEntityTags("service"); err != nil {
		return err
	}

	if err := checkDependentEntityTags("task"); err != nil {
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
