package test_aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsc "github.com/quintilesims/layer0/common/aws"
)

func describeSecurityGroupHelper(mockAWS *awsc.MockClient, securityGroupName, securityGroupID string) {
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

func deleteSecurityGroupHelper(mockAWS *awsc.MockClient, securityGroupID string) {
	input := &ec2.DeleteSecurityGroupInput{}
	input.SetGroupId(securityGroupID)

	mockAWS.EC2.EXPECT().
		DeleteSecurityGroup(input).
		Return(&ec2.DeleteSecurityGroupOutput{}, nil)
}
