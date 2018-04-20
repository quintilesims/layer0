package aws

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/retry"
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
		if err := deleteSGWithRetry(e.AWS.EC2, groupID); err != nil {
			return err
		}

		// Check for eventually consistency
		var err error
		fn := func() (shouldRetry bool) {
			filter := &ec2.Filter{}
			filter.SetName("group-name")
			filter.SetValues([]*string{aws.String(securityGroupName)})

			input := &ec2.DescribeSecurityGroupsInput{}
			input.SetFilters([]*ec2.Filter{filter})

			var output *ec2.DescribeSecurityGroupsOutput
			output, err = e.AWS.EC2.DescribeSecurityGroups(input)
			if err != nil {
				return false
			}

			for _, group := range output.SecurityGroups {
				if aws.StringValue(group.GroupName) == securityGroupName {
					log.Printf("[DEBUG] Service group not deleted, will retry lookup")
					err = errors.New(errors.EventualConsistencyError, err)
					return true
				}
			}

			return false
		}

		retry.Retry(fn, retry.WithTimeout(time.Second*30), retry.WithDelay(time.Second))

		return err
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

			msg := fmt.Sprintf("Cannot delete environment '%s' because it contains dependent %ss: ", environmentID, entityType)
			msg += strings.Join(entityIDs, ", ")
			return errors.Newf(errors.DependencyError, msg)
		}
		return nil
	}

	for _, entity := range []string{"load_balancer", "service", "task"} {
		if err := checkDependentEntityTags(entity); err != nil {
			return err
		}
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
