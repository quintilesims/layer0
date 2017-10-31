package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	awsc "github.com/quintilesims/layer0/common/aws"
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
