package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/aws/provider"
	"github.com/zpatrick/go-bytesize"
)

type Provider interface {
	CreateSecurityGroup(name, desc, vpcId string) (*string, error)
	AuthorizeSecurityGroupIngress(input []*SecurityGroupIngress) error
	RevokeSecurityGroupIngress(input []*SecurityGroupIngress) error
	RevokeSecurityGroupIngressHelper(groupID string, permission IpPermission) error
	AuthorizeSecurityGroupIngressFromGroup(groupId, sourceGroupId string) error
	DescribeSecurityGroup(name string) (*SecurityGroup, error)
	DescribeSubnet(subnetId string) (*Subnet, error)
	DeleteSecurityGroup(*SecurityGroup) error
	DescribeInstance(instanceId string) (*Instance, error)
	DescribeVPC(vpcID string) (*VPC, error)
	DescribeVPCByName(vpcName string) (*VPC, error)
	DescribeVPCSubnets(vpcId string) ([]*Subnet, error)
	DescribeVPCGateways(vpcId string) ([]*InternetGateway, error)
	DescribeVPCRoutes(vpcId string) ([]*RouteTable, error)
}

type EC2 struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (EC2Internal, error)
}

type EC2Internal interface {
	CreateSecurityGroup(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupIngress(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	RevokeSecurityGroupIngress(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error)
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DeleteSecurityGroup(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error)
	DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error)
}

// https://aws.amazon.com/ec2/instance-types/
var InstanceSizes = map[string]bytesize.Bytesize{
	"a1.medium":      2 * bytesize.GiB,
	"a1.large":       4 * bytesize.GiB,
	"a1.xlarge":      8 * bytesize.GiB,
	"a1.2xlarge":     16 * bytesize.GiB,
	"a1.4xlarge":     32 * bytesize.GiB,
	"c1.medium":      1.7 * bytesize.GiB,
	"c1.xlarge":      7 * bytesize.GiB,
	"c3.large":       3.75 * bytesize.GiB,
	"c3.xlarge":      7.5 * bytesize.GiB,
	"c3.2xlarge":     15 * bytesize.GiB,
	"c3.4xlarge":     30 * bytesize.GiB,
	"c3.8xlarge":     60 * bytesize.GiB,
	"c4.large":       3.75 * bytesize.GiB,
	"c4.xlarge":      7.5 * bytesize.GiB,
	"c4.2xlarge":     15 * bytesize.GiB,
	"c4.4xlarge":     30 * bytesize.GiB,
	"c4.8xlarge":     60 * bytesize.GiB,
	"c5.large":       4 * bytesize.GiB,
	"c5.xlarge":      8 * bytesize.GiB,
	"c5.2xlarge":     16 * bytesize.GiB,
	"c5.4xlarge":     32 * bytesize.GiB,
	"c5.9xlarge":     72 * bytesize.GiB,
	"c5.12xlarge":    96 * bytesize.GiB,
	"c5.18xlarge":    144 * bytesize.GiB,
	"c5.24xlarge":    192 * bytesize.GiB,
	"c5.metal":       192 * bytesize.GiB,
	"c5d.large":      4 * bytesize.GiB,
	"c5d.xlarge":     8 * bytesize.GiB,
	"c5d.2xlarge":    16 * bytesize.GiB,
	"c5d.4xlarge":    32 * bytesize.GiB,
	"c5d.9xlarge":    72 * bytesize.GiB,
	"c5d.18xlarge":   144 * bytesize.GiB,
	"c5n.large":      5.25 * bytesize.GiB,
	"c5n.xlarge":     10.5 * bytesize.GiB,
	"c5n.2xlarge":    21 * bytesize.GiB,
	"c5n.4xlarge":    42 * bytesize.GiB,
	"c5n.9xlarge":    96 * bytesize.GiB,
	"c5n.18xlarge":   192 * bytesize.GiB,
	"c5n.metal":      192 * bytesize.GiB,
	"cc2.8xlarge":    60.5 * bytesize.GiB,
	"cr1.8xlarge":    244 * bytesize.GiB,
	"d2.xlarge":      30.5 * bytesize.GiB,
	"d2.2xlarge":     61 * bytesize.GiB,
	"d2.4xlarge":     122 * bytesize.GiB,
	"d2.8xlarge":     244 * bytesize.GiB,
	"f1.2xlarge":     122 * bytesize.GiB,
	"f1.4xlarge":     244 * bytesize.GiB,
	"f1.16xlarge":    976 * bytesize.GiB,
	"g2.2xlarge":     15 * bytesize.GiB,
	"g2.8xlarge":     60 * bytesize.GiB,
	"g3.4xlarge":     122 * bytesize.GiB,
	"g3.8xlarge":     244 * bytesize.GiB,
	"g3.16xlarge":    488 * bytesize.GiB,
	"g3s.xlarge":     30.5 * bytesize.GiB,
	"g4dn.xlarge":    16 * bytesize.GiB,
	"g4dn.2xlarge":   32 * bytesize.GiB,
	"g4dn.4xlarge":   64 * bytesize.GiB,
	"g4dn.8xlarge":   128 * bytesize.GiB,
	"g4dn.16xlarge":  256 * bytesize.GiB,
	"g4dn.12xlarge":  192 * bytesize.GiB,
	"g4dn.metal":     384 * bytesize.GiB, // "coming soon" as of 2019/10/11
	"h1.2xlarge":     32 * bytesize.GiB,
	"h1.4xlarge":     64 * bytesize.GiB,
	"h1.8xlarge":     128 * bytesize.GiB,
	"h1.16xlarge":    256 * bytesize.GiB,
	"hs1.8xlarge":    117 * bytesize.GiB,
	"i2.xlarge":      30.5 * bytesize.GiB,
	"i2.2xlarge":     61 * bytesize.GiB,
	"i2.4xlarge":     122 * bytesize.GiB,
	"i2.8xlarge":     244 * bytesize.GiB,
	"i3.large":       15.25 * bytesize.GiB,
	"i3.xlarge":      30.5 * bytesize.GiB,
	"i3.2xlarge":     61 * bytesize.GiB,
	"i3.4xlarge":     122 * bytesize.GiB,
	"i3.8xlarge":     244 * bytesize.GiB,
	"i3.16xlarge":    488 * bytesize.GiB,
	"i3.metal":       512 * bytesize.GiB,
	"i3en.large":     16 * bytesize.GiB,
	"i3en.xlarge":    32 * bytesize.GiB,
	"i3en.2xlarge":   64 * bytesize.GiB,
	"i3en.3xlarge":   96 * bytesize.GiB,
	"i3en.6xlarge":   192 * bytesize.GiB,
	"i3en.12xlarge":  384 * bytesize.GiB,
	"i3en.24xlarge":  768 * bytesize.GiB,
	"i3en.metal":     768 * bytesize.GiB,
	"m1.small":       1.7 * bytesize.GiB,
	"m1.medium":      3.75 * bytesize.GiB,
	"m1.large":       7.5 * bytesize.GiB,
	"m1.xlarge":      15 * bytesize.GiB,
	"m2.xlarge":      17.1 * bytesize.GiB,
	"m2.2xlarge":     34.2 * bytesize.GiB,
	"m2.4xlarge":     68.4 * bytesize.GiB,
	"m3.medium":      3.75 * bytesize.GiB,
	"m3.large":       7.5 * bytesize.GiB,
	"m3.xlarge":      15 * bytesize.GiB,
	"m3.2xlarge":     30 * bytesize.GiB,
	"m4.large":       8 * bytesize.GiB,
	"m4.xlarge":      16 * bytesize.GiB,
	"m4.2xlarge":     32 * bytesize.GiB,
	"m4.4xlarge":     64 * bytesize.GiB,
	"m4.10xlarge":    160 * bytesize.GiB,
	"m4.16xlarge":    256 * bytesize.GiB,
	"m5.large":       8 * bytesize.GiB,
	"m5.xlarge":      16 * bytesize.GiB,
	"m5.2xlarge":     32 * bytesize.GiB,
	"m5.4xlarge":     64 * bytesize.GiB,
	"m5.8xlarge":     128 * bytesize.GiB,
	"m5.12xlarge":    192 * bytesize.GiB,
	"m5.16xlarge":    256 * bytesize.GiB,
	"m5.24xlarge":    384 * bytesize.GiB,
	"m5.metal":       384 * bytesize.GiB,
	"m5a.large":      8 * bytesize.GiB,
	"m5a.xlarge":     16 * bytesize.GiB,
	"m5a.2xlarge":    32 * bytesize.GiB,
	"m5a.4xlarge":    64 * bytesize.GiB,
	"m5a.8xlarge":    128 * bytesize.GiB,
	"m5a.12xlarge":   192 * bytesize.GiB,
	"m5a.16xlarge":   256 * bytesize.GiB,
	"m5a.24xlarge":   384 * bytesize.GiB,
	"m5ad.large":     8 * bytesize.GiB,
	"m5ad.xlarge":    16 * bytesize.GiB,
	"m5ad.2xlarge":   32 * bytesize.GiB,
	"m5ad.4xlarge":   64 * bytesize.GiB,
	"m5ad.12xlarge":  192 * bytesize.GiB,
	"m5ad.24xlarge":  384 * bytesize.GiB,
	"m5d.large":      8 * bytesize.GiB,
	"m5d.xlarge":     16 * bytesize.GiB,
	"m5d.2xlarge":    32 * bytesize.GiB,
	"m5d.4xlarge":    64 * bytesize.GiB,
	"m5d.8xlarge":    128 * bytesize.GiB,
	"m5d.12xlarge":   192 * bytesize.GiB,
	"m5d.16xlarge":   256 * bytesize.GiB,
	"m5d.24xlarge":   384 * bytesize.GiB,
	"m5d.metal":      384 * bytesize.GiB,
	"p2.xlarge":      61 * bytesize.GiB,
	"p2.8xlarge":     488 * bytesize.GiB,
	"p2.16xlarge":    732 * bytesize.GiB,
	"p3.2xlarge":     61 * bytesize.GiB,
	"p3.8xlarge":     244 * bytesize.GiB,
	"p3.16xlarge":    488 * bytesize.GiB,
	"p3dn.24xlarge":  768 * bytesize.GiB,
	"r3.large":       15.25 * bytesize.GiB,
	"r3.xlarge":      30.5 * bytesize.GiB,
	"r3.2xlarge":     61 * bytesize.GiB,
	"r3.4xlarge":     122 * bytesize.GiB,
	"r3.8xlarge":     244 * bytesize.GiB,
	"r4.large":       15.25 * bytesize.GiB,
	"r4.xlarge":      30.5 * bytesize.GiB,
	"r4.2xlarge":     61 * bytesize.GiB,
	"r4.4xlarge":     122 * bytesize.GiB,
	"r4.8xlarge":     244 * bytesize.GiB,
	"r4.16xlarge":    488 * bytesize.GiB,
	"r5.large":       16 * bytesize.GiB,
	"r5.xlarge":      32 * bytesize.GiB,
	"r5.2xlarge":     64 * bytesize.GiB,
	"r5.4xlarge":     128 * bytesize.GiB,
	"r5.8xlarge":     256 * bytesize.GiB,
	"r5.12xlarge":    384 * bytesize.GiB,
	"r5.16xlarge":    512 * bytesize.GiB,
	"r5.24xlarge":    768 * bytesize.GiB,
	"r5.metal":       768 * bytesize.GiB,
	"r5a.large":      16 * bytesize.GiB,
	"r5a.xlarge":     32 * bytesize.GiB,
	"r5a.2xlarge":    64 * bytesize.GiB,
	"r5a.4xlarge":    128 * bytesize.GiB,
	"r5a.8xlarge":    256 * bytesize.GiB,
	"r5a.12xlarge":   384 * bytesize.GiB,
	"r5a.16xlarge":   512 * bytesize.GiB,
	"r5a.24.xlarge":  768 * bytesize.GiB,
	"r5ad.large":     16 * bytesize.GiB,
	"r5ad.xlarge":    32 * bytesize.GiB,
	"r5ad.2xlarge":   64 * bytesize.GiB,
	"r5ad.4xlarge":   128 * bytesize.GiB,
	"r5ad.12xlarge":  384 * bytesize.GiB,
	"r5ad.24.xlarge": 768 * bytesize.GiB,
	"r5d.large":      16 * bytesize.GiB,
	"r5d.xlarge":     32 * bytesize.GiB,
	"r5d.2xlarge":    64 * bytesize.GiB,
	"r5d.4xlarge":    128 * bytesize.GiB,
	"r5d.8xlarge":    256 * bytesize.GiB,
	"r5d.12xlarge":   384 * bytesize.GiB,
	"r5d.16xlarge":   512 * bytesize.GiB,
	"r5d.24xlarge":   768 * bytesize.GiB,
	"r5d.metal":      768 * bytesize.GiB,
	"t1.micro":       0.613 * bytesize.GiB,
	"t2.nano":        0.5 * bytesize.GiB,
	"t2.micro":       1 * bytesize.GiB,
	"t2.small":       2 * bytesize.GiB,
	"t2.medium":      4 * bytesize.GiB,
	"t2.large":       8 * bytesize.GiB,
	"t2.xlarge":      16 * bytesize.GiB,
	"t2.2xlarge":     32 * bytesize.GiB,
	"t3.nano":        0.5 * bytesize.GiB,
	"t3.micro":       1 * bytesize.GiB,
	"t3.small":       2 * bytesize.GiB,
	"t3.medium":      4 * bytesize.GiB,
	"t3.large":       8 * bytesize.GiB,
	"t3.xlarge":      16 * bytesize.GiB,
	"t3.2xlarge":     32 * bytesize.GiB,
	"t3a.nano":       0.5 * bytesize.GiB,
	"t3a.micro":      1 * bytesize.GiB,
	"t3a.small":      2 * bytesize.GiB,
	"t3a.medium":     4 * bytesize.GiB,
	"t3a.large":      8 * bytesize.GiB,
	"t3a.xlarge":     16 * bytesize.GiB,
	"t3a.2xlarge":    32 * bytesize.GiB,
	"u-6tb1.metal":   6144 * bytesize.GiB,
	"u-9tb1.metal":   9216 * bytesize.GiB,
	"u-12tb1.metal":  12288 * bytesize.GiB,
	"u-18tb1.metal":  18432 * bytesize.GiB,
	"u-24tb1.metal":  24576 * bytesize.GiB,
	"x1.16xlarge":    976 * bytesize.GiB,
	"x1.32xlarge":    1952 * bytesize.GiB,
	"x1e.xlarge":     122 * bytesize.GiB,
	"x1e.2xlarge":    244 * bytesize.GiB,
	"x1e.4xlarge":    488 * bytesize.GiB,
	"x1e.8xlarge":    976 * bytesize.GiB,
	"x1e.16xlarge":   1952 * bytesize.GiB,
	"x1e.32xlarge":   3904 * bytesize.GiB,
	"z1d.large":      16 * bytesize.GiB,
	"z1d.xlarge":     32 * bytesize.GiB,
	"z1d.2xlarge":    64 * bytesize.GiB,
	"z1d.3xlarge":    96 * bytesize.GiB,
	"z1d.6xlarge":    192 * bytesize.GiB,
	"z1d.12xlarge":   384 * bytesize.GiB,
	"z1d.metal":      384 * bytesize.GiB,
}

type SecurityGroup struct {
	*ec2.SecurityGroup
}

func NewSecurityGroup(id string) *SecurityGroup {
	return &SecurityGroup{
		&ec2.SecurityGroup{
			GroupId: aws.String(id),
		},
	}
}

type SecurityGroupIngress struct {
	*ec2.AuthorizeSecurityGroupIngressInput
	*ec2.RevokeSecurityGroupIngressInput
}

func NewSecurityGroupIngress(groupID, cidrIP, protocol string, fromPort, toPort int) *SecurityGroupIngress {
	return &SecurityGroupIngress{
		AuthorizeSecurityGroupIngressInput: &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:    aws.String(groupID),
			CidrIp:     aws.String(cidrIP),
			FromPort:   aws.Int64(int64(fromPort)),
			ToPort:     aws.Int64(int64(toPort)),
			IpProtocol: aws.String(protocol),
		},
		RevokeSecurityGroupIngressInput: &ec2.RevokeSecurityGroupIngressInput{
			GroupId:    aws.String(groupID),
			CidrIp:     aws.String(cidrIP),
			FromPort:   aws.Int64(int64(fromPort)),
			ToPort:     aws.Int64(int64(toPort)),
			IpProtocol: aws.String(protocol),
		},
	}
}

type IpPermission struct {
	*ec2.IpPermission
}

type UserIdGroupPair struct {
	*ec2.UserIdGroupPair
}

type Subnet struct {
	*ec2.Subnet
}

func NewSubnet() *Subnet {
	return &Subnet{&ec2.Subnet{}}
}

type Instance struct {
	*ec2.Instance
}

func NewInstance() *Instance {
	return &Instance{&ec2.Instance{}}
}

type VPC struct {
	*ec2.Vpc
}

func NewVpc() *VPC {
	return &VPC{&ec2.Vpc{}}
}

type InternetGateway struct {
	*ec2.InternetGateway
}

func NewInternetGateway() *InternetGateway {
	return &InternetGateway{&ec2.InternetGateway{}}
}

type RouteTable struct {
	*ec2.RouteTable
}

func NewRouteTable() *RouteTable {
	return &RouteTable{&ec2.RouteTable{}}
}

func NewEC2(credProvider provider.CredProvider, region string) (Provider, error) {
	ec2 := EC2{
		credProvider,
		region,
		func() (EC2Internal, error) {
			return Connect(credProvider, region)
		},
	}
	_, err := ec2.Connect()
	if err != nil {
		return nil, err
	}
	return &ec2, nil
}

func Connect(credProvider provider.CredProvider, region string) (EC2Internal, error) {
	connection, err := provider.GetEC2Connection(credProvider, region)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (this *EC2) CreateSecurityGroup(name, desc, vpcId string) (*string, error) {
	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String(desc),
		GroupName:   aws.String(name),
		VpcId:       aws.String(vpcId),
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}
	output, err := connection.CreateSecurityGroup(input)
	return output.GroupId, err
}

func (this *EC2) AuthorizeSecurityGroupIngress(ingresses []*SecurityGroupIngress) error {
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	for _, ingress := range ingresses {
		input := ingress.AuthorizeSecurityGroupIngressInput
		if _, err := connection.AuthorizeSecurityGroupIngress(input); err != nil {
			return err
		}
	}

	return nil
}

func (this *EC2) RevokeSecurityGroupIngress(ingresses []*SecurityGroupIngress) error {
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	for _, ingress := range ingresses {
		input := ingress.RevokeSecurityGroupIngressInput

		if _, err := connection.RevokeSecurityGroupIngress(input); err != nil {
			return err
		}
	}

	return nil
}

func (this *EC2) RevokeSecurityGroupIngressHelper(groupID string, permission IpPermission) error {
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	input := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       aws.String(groupID),
		IpPermissions: []*ec2.IpPermission{permission.IpPermission},
	}

	if _, err := connection.RevokeSecurityGroupIngress(input); err != nil {
		return err
	}

	return nil
}

func (this *EC2) AuthorizeSecurityGroupIngressFromGroup(groupId, sourceGroupId string) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(groupId),
		IpPermissions: []*ec2.IpPermission{
			{
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{GroupId: aws.String(sourceGroupId)},
				},
				IpProtocol: aws.String("-1"),
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return err
	}
	_, err = connection.AuthorizeSecurityGroupIngress(input)
	return err
}

func (this *EC2) DescribeSecurityGroup(name string) (*SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(name)},
			},
		},
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	groups, err := connection.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	if len(groups.SecurityGroups) == 0 {
		return nil, nil
	}

	return &SecurityGroup{groups.SecurityGroups[0]}, nil
}

func (this *EC2) DescribeSubnet(subnetId string) (*Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: []*string{aws.String(subnetId)},
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	subnets, err := connection.DescribeSubnets(input)
	var subnet *Subnet
	if err == nil {
		if len(subnets.Subnets) > 0 {
			subnet = &Subnet{subnets.Subnets[0]}
		}
	}
	return subnet, err
}

func (this *EC2) DeleteSecurityGroup(group *SecurityGroup) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: group.GroupId,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteSecurityGroup(input)
	return err
}

func (this *EC2) DescribeInstance(instanceId string) (*Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceId},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	resp, err := connection.DescribeInstances(input)
	var instance *Instance
	if err == nil {
		if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
			instance = &Instance{resp.Reservations[0].Instances[0]}
		}
	}
	return instance, err
}

func (this *EC2) DescribeVPC(vpcId string) (*VPC, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{&vpcId},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	resp, err := connection.DescribeVpcs(input)
	var vpc *VPC
	if err == nil {
		if len(resp.Vpcs) > 0 {
			vpc = &VPC{resp.Vpcs[0]}
		}
	}
	return vpc, err
}

func (this *EC2) DescribeVPCByName(vpcName string) (*VPC, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{&vpcName},
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	var vpc *VPC
	resp, err := connection.DescribeVpcs(input)
	if err == nil {
		if len(resp.Vpcs) > 0 {
			vpc = &VPC{resp.Vpcs[0]}
		}
	}
	return vpc, err
}

func (this *EC2) DescribeVPCSubnets(vpcId string) ([]*Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	subnets, err := connection.DescribeSubnets(input)
	var subnetList []*Subnet
	if err == nil {
		fmt.Printf("Found %d subnets\n", len(subnets.Subnets))

		length := len(subnets.Subnets)
		subnetList = make([]*Subnet, length, length)
		for i, net := range subnets.Subnets {
			subnetList[i] = &Subnet{net}
		}
	}
	return subnetList, err
}

func (this *EC2) DescribeVPCGateways(vpcId string) ([]*InternetGateway, error) {
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	gateways, err := connection.DescribeInternetGateways(input)
	var result []*InternetGateway
	if err == nil {
		fmt.Printf("Found %d gateways\n", len(gateways.InternetGateways))

		length := len(gateways.InternetGateways)
		result = make([]*InternetGateway, length, length)
		for i, net := range gateways.InternetGateways {
			result[i] = &InternetGateway{net}
		}
	}
	return result, err
}

func (this *EC2) DescribeVPCRoutes(vpcId string) ([]*RouteTable, error) {
	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeRouteTables(input)
	var result []*RouteTable
	if err == nil {
		list := output.RouteTables
		fmt.Printf("Found %d routes\n", len(list))

		length := len(list)
		result = make([]*RouteTable, length, length)
		for i, o := range list {
			result[i] = &RouteTable{o}
		}
	}
	return result, err
}
