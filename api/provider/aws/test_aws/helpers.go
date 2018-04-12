package test_aws

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func createSGHelper(t *testing.T, mockAWS *awsc.MockClient, groupName, vpcID string) {
	// use this method so we can ignore validation for the Description field
	validateInput := func(input *ec2.CreateSecurityGroupInput) {
		assert.Equal(t, groupName, aws.StringValue(input.GroupName))
		assert.Equal(t, vpcID, aws.StringValue(input.VpcId))
	}

	mockAWS.EC2.EXPECT().
		CreateSecurityGroup(gomock.Any()).
		Do(validateInput).
		Return(&ec2.CreateSecurityGroupOutput{}, nil)
}

func readSGHelper(mockAWS *awsc.MockClient, securityGroupName, securityGroupID string) {
	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String(securityGroupName)})

	input := &ec2.DescribeSecurityGroupsInput{}
	input.SetFilters([]*ec2.Filter{filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName(securityGroupName)
	securityGroup.SetGroupId(securityGroupID)

	output := &ec2.DescribeSecurityGroupsOutput{}
	output.SetSecurityGroups([]*ec2.SecurityGroup{securityGroup})

	// todo: check input
	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(input).
		Return(output, nil)
}

func deleteSGHelper(mockAWS *awsc.MockClient, securityGroupID string) {
	input := &ec2.DeleteSecurityGroupInput{}
	input.SetGroupId(securityGroupID)

	mockAWS.EC2.EXPECT().
		DeleteSecurityGroup(input).
		Return(&ec2.DeleteSecurityGroupOutput{}, nil)
}

func healthCheckHelper(healthCheck *models.HealthCheck) *elb.HealthCheck {
	elbHealthCheck := &elb.HealthCheck{}
	elbHealthCheck.SetTarget(healthCheck.Target)
	elbHealthCheck.SetInterval(int64(healthCheck.Interval))
	elbHealthCheck.SetTimeout(int64(healthCheck.Timeout))
	elbHealthCheck.SetHealthyThreshold(int64(healthCheck.HealthyThreshold))
	elbHealthCheck.SetUnhealthyThreshold(int64(healthCheck.UnhealthyThreshold))

	return elbHealthCheck
}

func healthCheckTargetGroupHelper(healthCheck *models.HealthCheck) *alb.ModifyTargetGroupInput {
	albHealthCheck := &alb.ModifyTargetGroupInput{}
	albHealthCheck.SetHealthCheckIntervalSeconds(int64(healthCheck.Interval))
	albHealthCheck.SetHealthCheckPath(healthCheck.Path)
	albHealthCheck.SetHealthCheckTimeoutSeconds(int64(healthCheck.Timeout))
	albHealthCheck.SetHealthyThresholdCount(int64(healthCheck.HealthyThreshold))
	albHealthCheck.SetUnhealthyThresholdCount(int64(healthCheck.UnhealthyThreshold))

	return albHealthCheck
}

func listenerHelper(port models.Port) *elb.Listener {
	listener := &elb.Listener{}
	listener.SetProtocol(port.Protocol)
	listener.SetLoadBalancerPort(port.HostPort)
	listener.SetInstancePort(port.ContainerPort)

	switch strings.ToLower(port.Protocol) {
	case "http", "https":
		listener.SetInstanceProtocol("http")
	case "tcp", "ssl":
		listener.SetInstanceProtocol("tcp")
	}

	return listener
}

func authorizeSGIngressHelper(port models.Port) *ec2.AuthorizeSecurityGroupIngressInput {
	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetCidrIp("0.0.0.0/0")
	ingressInput.SetIpProtocol("TCP")
	ingressInput.SetFromPort(int64(port.HostPort))
	ingressInput.SetToPort(int64(port.HostPort))

	return ingressInput
}

func revokeSGIngressHelper(port models.Port) *ec2.RevokeSecurityGroupIngressInput {
	ingressInput := &ec2.RevokeSecurityGroupIngressInput{}
	ingressInput.SetCidrIp("0.0.0.0/0")
	ingressInput.SetIpProtocol("TCP")
	ingressInput.SetFromPort(port.HostPort)
	ingressInput.SetToPort(port.HostPort)

	return ingressInput
}
