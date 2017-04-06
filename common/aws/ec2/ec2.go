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
	RevokeSecurityGroupIngressHelper(groupID string, permission *IpPermission) error
	AuthorizeSecurityGroupIngressFromGroup(groupId, sourceGroupId *string) error
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
	"t2.nano":     0.5 * bytesize.GiB,
	"t2.micro":    1 * bytesize.GiB,
	"t2.small":    2 * bytesize.GiB,
	"t2.medium":   4 * bytesize.GiB,
	"t2.large":    8 * bytesize.GiB,
	"m4.large":    8 * bytesize.GiB,
	"m4.xlarge":   16 * bytesize.GiB,
	"m4.2xlarge":  32 * bytesize.GiB,
	"m4.4xlarge":  64 * bytesize.GiB,
	"m4.10xlarge": 160 * bytesize.GiB,
	"m3.medium":   3.75 * bytesize.GiB,
	"m3.large":    7.5 * bytesize.GiB,
	"m3.xlarge":   15 * bytesize.GiB,
	"m3.2xlarge":  30 * bytesize.GiB,
	"c4.large":    3.75 * bytesize.GiB,
	"c4.xlarge":   7.5 * bytesize.GiB,
	"c4.2xlarge":  15 * bytesize.GiB,
	"c4.4xlarge":  30 * bytesize.GiB,
	"c4.8xlarge":  60 * bytesize.GiB,
	"c3.large":    3.75 * bytesize.GiB,
	"c3.xlarge":   7.5 * bytesize.GiB,
	"c3.2xlarge":  15 * bytesize.GiB,
	"c3.4xlarge":  30 * bytesize.GiB,
	"c3.8xlarge":  60 * bytesize.GiB,
	"g2.2xlarge":  15 * bytesize.GiB,
	"g2.8xlarge":  60 * bytesize.GiB,
	"x1.32xlarge": 1952 * bytesize.GiB,
	"r3.large":    15.25 * bytesize.GiB,
	"r3.xlarge":   30.5 * bytesize.GiB,
	"r3.2xlarge":  61 * bytesize.GiB,
	"r3.4xlarge":  122 * bytesize.GiB,
	"r3.8xlarge":  244 * bytesize.GiB,
	"i3.large":    15.25 * bytesize.GiB,
	"i3.xlarge":   30.5 * bytesize.GiB,
	"i3.2xlarge":  61 * bytesize.GiB,
	"i3.4xlarge":  122 * bytesize.GiB,
	"i3.8xlarge":  244 * bytesize.GiB,
	"d2.xlarge":   30.5 * bytesize.GiB,
	"d2.2xlarge":  61 * bytesize.GiB,
	"d2.4xlarge":  122 * bytesize.GiB,
	"d2.8xlarge":  244 * bytesize.GiB,
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

func (this *EC2) RevokeSecurityGroupIngressHelper(groupID string, permission *IpPermission) error {
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

func (this *EC2) AuthorizeSecurityGroupIngressFromGroup(groupId, sourceGroupId *string) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: groupId,
		IpPermissions: []*ec2.IpPermission{
			&ec2.IpPermission{
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					&ec2.UserIdGroupPair{GroupId: sourceGroupId},
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
			&ec2.Filter{
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
			&ec2.Filter{
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
			&ec2.Filter{
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
			&ec2.Filter{
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
			&ec2.Filter{
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
			&ec2.Filter{
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
