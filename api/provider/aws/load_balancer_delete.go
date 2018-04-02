package aws

import (
	"log"
	"strings"
	"time"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/retry"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Delete is used to delete an Elastic Load Balancer using the specified loadBalancerID.
// The associated IAM Role, IAM Role inline policy, and Security Group are also
// removed as part of the process by making DeleteRole, DeleteRolePolicy and
// DeleteSecurityGroup requests to AWS, respectively. The Load Balancer is deleted
// by making a DeleteLoadBalancer request to AWS.
func (l *LoadBalancerProvider) Delete(loadBalancerID string) error {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)

	if err := l.deleteLoadBalancer(fqLoadBalancerID); err != nil {
		return err
	}

	// Check for eventually consistency
	waitUntilLBDeletedFN := func() (shouldRetry bool, err error) {
		input := &elb.DescribeLoadBalancersInput{}
		input.SetLoadBalancerNames([]*string{aws.String(fqLoadBalancerID)})
		input.SetPageSize(1)

		if _, err = l.AWS.ELB.DescribeLoadBalancers(input); err != nil {
			if err, ok := err.(awserr.Error); ok && err.Code() == "LoadBalancerNotFound" {
				return false, nil
			}

			return false, err
		}

		log.Printf("[DEBUG] Load Balancer not deleted, will retry lookup")
		return true, nil
	}

	if err := retry.Retry(waitUntilLBDeletedFN, retry.WithTimeout(time.Second*30), retry.WithDelay(time.Second)); err != nil {
		return errors.New(errors.EventualConsistencyError, err)
	}

	roleName := getLoadBalancerRoleName(fqLoadBalancerID)
	policyName := roleName
	if err := l.deleteRolePolicy(roleName, policyName); err != nil {
		return err
	}

	if err := l.deleteRole(roleName); err != nil {
		return err
	}

	securityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
	securityGroup, err := readSG(l.AWS.EC2, securityGroupName)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return err
	}

	if securityGroup != nil {
		groupID := aws.StringValue(securityGroup.GroupId)
		if err := deleteSG(l.AWS.EC2, groupID); err != nil {
			return err
		}
	}

	// Check for eventually consistency
	fn := waitUntilSGDeletedFN(l.AWS.EC2, securityGroupName)
	if err := retry.Retry(fn, retry.WithTimeout(time.Second*30), retry.WithDelay(time.Second)); err != nil {
		return errors.New(errors.EventualConsistencyError, err)
	}

	if err := l.deleteTags(loadBalancerID); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) deleteLoadBalancer(loadBalancerName string) error {
	input := &elb.DeleteLoadBalancerInput{}
	input.SetLoadBalancerName(loadBalancerName)

	if _, err := l.AWS.ELB.DeleteLoadBalancer(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "NoSuchEntity" {
			return nil
		}

		return err
	}

	return nil
}

func (l *LoadBalancerProvider) deleteRolePolicy(roleName, policyName string) error {
	input := &iam.DeleteRolePolicyInput{}
	input.SetRoleName(roleName)
	input.SetPolicyName(policyName)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.IAM.DeleteRolePolicy(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "NoSuchEntity" {
			return nil
		}

		return err
	}

	return nil
}

func (l *LoadBalancerProvider) deleteRole(roleName string) error {
	input := &iam.DeleteRoleInput{}
	input.SetRoleName(roleName)

	// todo: validate NoSuchEntity is correct error code
	if _, err := l.AWS.IAM.DeleteRole(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "NoSuchEntity" {
			return nil
		}

		return err
	}

	return nil
}

func (l *LoadBalancerProvider) deleteTags(loadBalancerID string) error {
	tags, err := l.TagStore.SelectByTypeAndID("load_balancer", loadBalancerID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := l.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			return err
		}
	}

	return nil
}
