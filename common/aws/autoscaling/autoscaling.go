package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/quintilesims/layer0/common/aws/provider"
)

type Provider interface {
	AttachLoadBalancer(autoScalingGroupName, loadBalancerName string) error
	CreateLaunchConfiguration(name, amiID, iamInstanceProfile, instanceType, keyName, userData *string, securityGroups []*string, volSize map[string]int) error
	CreateAutoScalingGroup(name, launchConfigName, subnets string, minCount, maxCount int) error
	SetDesiredCapacity(name string, size int) error
	UpdateAutoScalingGroupMaxSize(name string, size int) error
	UpdateAutoScalingGroupMinSize(name string, size int) error
	DescribeAutoScalingGroups(names []*string) ([]*Group, error)
	DescribeAutoScalingGroup(name string) (*Group, error)
	DescribeLaunchConfigurations(names []*string) ([]*LaunchConfiguration, error)
	DescribeLaunchConfiguration(name string) (*LaunchConfiguration, error)
	DeleteAutoScalingGroup(name *string) error
	DeleteLaunchConfiguration(name *string) error
	TerminateInstanceInAutoScalingGroup(instanceID string, decrement bool) (*Activity, error)
}

type AutoScaling struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (AutoScalingInternal, error)
}

type AutoScalingInternal interface {
	AttachLoadBalancers(input *autoscaling.AttachLoadBalancersInput) (output *autoscaling.AttachLoadBalancersOutput, err error)
	CreateLaunchConfiguration(input *autoscaling.CreateLaunchConfigurationInput) (*autoscaling.CreateLaunchConfigurationOutput, error)
	CreateAutoScalingGroup(input *autoscaling.CreateAutoScalingGroupInput) (*autoscaling.CreateAutoScalingGroupOutput, error)
	DescribeAutoScalingGroups(input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	DescribeLaunchConfigurations(input *autoscaling.DescribeLaunchConfigurationsInput) (*autoscaling.DescribeLaunchConfigurationsOutput, error)
	SetDesiredCapacity(input *autoscaling.SetDesiredCapacityInput) (*autoscaling.SetDesiredCapacityOutput, error)
	UpdateAutoScalingGroup(input *autoscaling.UpdateAutoScalingGroupInput) (*autoscaling.UpdateAutoScalingGroupOutput, error)
	DeleteAutoScalingGroup(input *autoscaling.DeleteAutoScalingGroupInput) (*autoscaling.DeleteAutoScalingGroupOutput, error)
	DeleteLaunchConfiguration(input *autoscaling.DeleteLaunchConfigurationInput) (*autoscaling.DeleteLaunchConfigurationOutput, error)
	TerminateInstanceInAutoScalingGroup(input *autoscaling.TerminateInstanceInAutoScalingGroupInput) (*autoscaling.TerminateInstanceInAutoScalingGroupOutput, error)
}

type Group struct {
	*autoscaling.Group
}

type LaunchConfiguration struct {
	*autoscaling.LaunchConfiguration
}

type Activity struct {
	*autoscaling.Activity
}

func NewGroup() *Group {
	return &Group{&autoscaling.Group{
		MinSize:             aws.Int64(0),
		MaxSize:             aws.Int64(0),
		DesiredCapacity:     aws.Int64(0),
		AutoScalingGroupARN: aws.String(""),
	}}
}

func NewLaunchConfiguration(size, ami string) *LaunchConfiguration {
	return &LaunchConfiguration{&autoscaling.LaunchConfiguration{
		InstanceType: &size,
		ImageId:      &ami,
	}}
}

func NewAutoScaling(credProvider provider.CredProvider, region string) (Provider, error) {
	autoScaling := AutoScaling{
		credProvider,
		region,
		func() (AutoScalingInternal, error) {
			return Connect(credProvider, region)
		},
	}
	_, err := autoScaling.Connect()
	if err != nil {
		return nil, err
	}
	return &autoScaling, nil
}

func Connect(credProvider provider.CredProvider, region string) (AutoScalingInternal, error) {
	connection, err := provider.GetAutoScalingConnection(credProvider, region)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (this *AutoScaling) AttachLoadBalancer(autoScalingGroupName, loadBalancerName string) (err error) {

	input := &autoscaling.AttachLoadBalancersInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
		LoadBalancerNames:    []*string{aws.String(loadBalancerName)},
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}
	_, err = connection.AttachLoadBalancers(input)
	return err
}

func (this *AutoScaling) CreateLaunchConfiguration(
	name *string,
	amiID *string,
	iamInstanceProfile *string,
	instanceType *string,
	keyName *string,
	userData *string,
	securityGroups []*string,
	volSizes map[string]int,
) error {
	if *keyName == "" {
		keyName = nil
	}

	blocks := []*autoscaling.BlockDeviceMapping{}
	for vol, size := range volSizes {
		block := &autoscaling.BlockDeviceMapping{
			DeviceName: aws.String(vol),
			//VirtualName?
			Ebs: &autoscaling.Ebs{
				DeleteOnTermination: aws.Bool(true),
				VolumeSize:          aws.Int64(int64(size)),
				VolumeType:          aws.String("gp2"),
			},
		}

		blocks = append(blocks, block)
	}

	input := &autoscaling.CreateLaunchConfigurationInput{
		ImageId:                 amiID,
		IamInstanceProfile:      iamInstanceProfile,
		InstanceType:            instanceType,
		KeyName:                 keyName,
		UserData:                userData,
		LaunchConfigurationName: name,
		SecurityGroups:          securityGroups,
		BlockDeviceMappings:     blocks,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}
	_, err = connection.CreateLaunchConfiguration(input)
	return err
}

func (this *AutoScaling) CreateAutoScalingGroup(name, launchConfigName, subnets string, minSize, maxSize int) error {
	input := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName:             aws.String(name),
		DesiredCapacity:                  aws.Int64(int64(maxSize)),
		MinSize:                          aws.Int64(int64(minSize)),
		MaxSize:                          aws.Int64(int64(maxSize)),
		LaunchConfigurationName:          aws.String(launchConfigName),
		VPCZoneIdentifier:                aws.String(subnets),
		NewInstancesProtectedFromScaleIn: aws.Bool(true),
		Tags: []*autoscaling.Tag{
			{
				Key:               aws.String("Name"),
				Value:             aws.String(name),
				PropagateAtLaunch: aws.Bool(true),
			},
		},
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}
	_, err = connection.CreateAutoScalingGroup(input)
	return err
}

func (this *AutoScaling) SetDesiredCapacity(name string, size int) error {
	size64 := int64(size)
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      &size64,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}
	_, err = connection.SetDesiredCapacity(input)
	return err
}

func (this *AutoScaling) DescribeAutoScalingGroup(name string) (*Group, error) {
	groups, err := this.DescribeAutoScalingGroups([]*string{&name})
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("Autoscaling group '%s' not found", name)
	}

	if len(groups) > 1 {
		return nil, fmt.Errorf("Unexpected number of autoscaling groups with name '%s'", name)
	}

	return groups[0], nil
}

func (this *AutoScaling) DescribeAutoScalingGroups(names []*string) ([]*Group, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: names,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}

	groups := []*Group{}
	for _, g := range out.AutoScalingGroups {
		groups = append(groups, &Group{g})
	}

	return groups, nil
}

func (this *AutoScaling) DescribeLaunchConfiguration(name string) (*LaunchConfiguration, error) {
	configs, err := this.DescribeLaunchConfigurations([]*string{&name})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("Launch configuration '%s' not found.", name)
	}

	if len(configs) > 1 {
		return nil, fmt.Errorf("Unexpected number of launch configurations with name '%s'.", name)
	}

	return configs[0], nil
}

func (this *AutoScaling) DescribeLaunchConfigurations(names []*string) ([]*LaunchConfiguration, error) {
	input := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: names,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.DescribeLaunchConfigurations(input)
	if err != nil {
		return nil, err
	}

	configs := []*LaunchConfiguration{}
	for _, s := range out.LaunchConfigurations {
		configs = append(configs, &LaunchConfiguration{s})
	}

	return configs, nil
}

func (this *AutoScaling) UpdateAutoScalingGroupMaxSize(name string, size int) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(int64(size)),
		MaxSize:              aws.Int64(int64(size)),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.UpdateAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}

func (this *AutoScaling) UpdateAutoScalingGroupMinSize(name string, size int) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		MinSize:              aws.Int64(int64(size)),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.UpdateAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}

func (this *AutoScaling) DeleteAutoScalingGroup(name *string) error {
	input := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: name,
		ForceDelete:          aws.Bool(true),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteAutoScalingGroup(input)
	if err != nil {
		return err
	}
	return nil
}

func (this *AutoScaling) DeleteLaunchConfiguration(name *string) error {
	input := &autoscaling.DeleteLaunchConfigurationInput{
		LaunchConfigurationName: name,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteLaunchConfiguration(input)
	if err != nil {
		return err
	}
	return nil
}

func (this *AutoScaling) TerminateInstanceInAutoScalingGroup(instanceID string, decrement bool) (*Activity, error) {
	input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
		InstanceId:                     &instanceID,
		ShouldDecrementDesiredCapacity: &decrement,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.TerminateInstanceInAutoScalingGroup(input)
	if err != nil {
		return nil, err
	}

	return &Activity{out.Activity}, nil
}
