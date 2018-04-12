package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreate_stateless(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"subnet-test"})

	defer provider.SetEntityIDGenerator("svc_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// define expected ec2.DescribeSecurityGroups
	// (part of "awsvpc" NetworkMode workflow)
	ec2Filter := &ec2.Filter{}
	ec2Filter.SetName("group-name")
	ec2Filter.SetValues([]*string{aws.String("l0-test-env_id-env")})

	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{}
	describeSecurityGroupsInput.SetFilters([]*ec2.Filter{ec2Filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-env_id-env")
	securityGroup.SetGroupId("sg-test")

	securityGroups := []*ec2.SecurityGroup{securityGroup}

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups(securityGroups)

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(describeSecurityGroupsInput).
		Return(describeSecurityGroupsOutput, nil)

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	networkMode := ecs.NetworkModeAwsvpc

	taskDefinition := &ecs.TaskDefinition{
		Compatibilities:      []*string{aws.String(ecs.LaunchTypeEc2), aws.String(ecs.LaunchTypeFargate)},
		ContainerDefinitions: containerDefinitions,
		NetworkMode:          &networkMode,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected ecs.CreateService
	awsvpcConfig := &ecs.AwsVpcConfiguration{}
	awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
	awsvpcConfig.SetSecurityGroups([]*string{aws.String("sg-test")})
	awsvpcConfig.SetSubnets([]*string{aws.String("subnet-test")})

	networkConfig := &ecs.NetworkConfiguration{}
	networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

	createServiceInput := &ecs.CreateServiceInput{}
	createServiceInput.SetCluster("l0-test-env_id")
	createServiceInput.SetDesiredCount(0)
	createServiceInput.SetLaunchType(ecs.LaunchTypeFargate)
	createServiceInput.SetNetworkConfiguration(networkConfig)
	createServiceInput.SetPlatformVersion(config.DefaultFargatePlatformVersion)
	createServiceInput.SetServiceName("l0-test-svc_id")
	createServiceInput.SetTaskDefinition("dpl_arn")

	mockAWS.ECS.EXPECT().
		CreateService(createServiceInput).
		Return(&ecs.CreateServiceOutput{}, nil)

	// define request
	req := models.CreateServiceRequest{
		DeployID:      "dpl_id",
		EnvironmentID: "env_id",
		ServiceName:   "svc_name",
		Stateful:      false,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result, "svc_id")

	expectedTags := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestServiceCreate_stateful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	defer provider.SetEntityIDGenerator("svc_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	taskDefinition := &ecs.TaskDefinition{
		Compatibilities:      []*string{aws.String(ecs.LaunchTypeEc2)},
		ContainerDefinitions: containerDefinitions,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected CreateService
	createServiceInput := &ecs.CreateServiceInput{}
	createServiceInput.SetCluster("l0-test-env_id")
	createServiceInput.SetDesiredCount(2)
	createServiceInput.SetLaunchType(ecs.LaunchTypeEc2)
	createServiceInput.SetServiceName("l0-test-svc_id")
	createServiceInput.SetTaskDefinition("dpl_arn")

	mockAWS.ECS.EXPECT().
		CreateService(createServiceInput).
		Return(&ecs.CreateServiceOutput{}, nil)

	// define request
	req := models.CreateServiceRequest{
		DeployID:      "dpl_id",
		EnvironmentID: "env_id",
		Scale:         2,
		ServiceName:   "svc_name",
		Stateful:      true,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result, "svc_id")

	expectedTags := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestServiceCreate_statelessLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"subnet-test"})

	defer provider.SetEntityIDGenerator("svc_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// define expected ec2.DescribeSecurityGroups
	// (part of "awsvpc" NetworkMode workflow)
	ec2Filter := &ec2.Filter{}
	ec2Filter.SetName("group-name")
	ec2Filter.SetValues([]*string{aws.String("l0-test-env_id-env")})

	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{}
	describeSecurityGroupsInput.SetFilters([]*ec2.Filter{ec2Filter})

	securityGroup := &ec2.SecurityGroup{}
	securityGroup.SetGroupName("l0-test-env_id-env")
	securityGroup.SetGroupId("sg-test")
	securityGroups := []*ec2.SecurityGroup{securityGroup}

	describeSecurityGroupsOutput := &ec2.DescribeSecurityGroupsOutput{}
	describeSecurityGroupsOutput.SetSecurityGroups(securityGroups)

	mockAWS.EC2.EXPECT().
		DescribeSecurityGroups(describeSecurityGroupsInput).
		Return(describeSecurityGroupsOutput, nil)

	// Define expected load balancers
	//
	// First we check for an ELB load balancer.
	// This needs to return a LoadBalancerNotFound error
	// so that we can then check for an ALB load balancer.
	elbInput := &elb.DescribeLoadBalancersInput{}
	elbInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_id")})
	elbInput.SetPageSize(1)

	elbErr := awserr.New(alb.ErrCodeLoadBalancerNotFoundException, "", nil)

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(elbInput).
		Return(nil, elbErr)

	// Now we can check for an ALB load balancer.
	albInput := &alb.DescribeLoadBalancersInput{}
	albInput.SetNames([]*string{aws.String("l0-test-lb_id")})

	albLoadBalancer := &alb.LoadBalancer{}
	albLoadBalancers := []*alb.LoadBalancer{
		albLoadBalancer,
	}

	albOutput := &alb.DescribeLoadBalancersOutput{}
	albOutput.SetLoadBalancers(albLoadBalancers)

	mockAWS.ALB.EXPECT().
		DescribeLoadBalancers(albInput).
		Return(albOutput, nil)

	// Once there is an ALB, the TargetGroup is checked.
	describeTargetGroupsInput := &alb.DescribeTargetGroupsInput{}
	describeTargetGroupsInput.SetNames([]*string{aws.String("l0-test-lb_id")})

	targetGroup := &alb.TargetGroup{}
	targetGroup.SetTargetGroupArn("lb_arn")
	targetGroup.SetHealthCheckPath("")
	targetGroup.SetHealthCheckIntervalSeconds(0)
	targetGroup.SetHealthCheckTimeoutSeconds(0)
	targetGroup.SetHealthyThresholdCount(0)
	targetGroup.SetUnhealthyThresholdCount(0)

	targetGroups := []*alb.TargetGroup{
		targetGroup,
	}

	describeTargetGroupsOutput := &alb.DescribeTargetGroupsOutput{}
	describeTargetGroupsOutput.SetTargetGroups(targetGroups)

	mockAWS.ALB.EXPECT().
		DescribeTargetGroups(describeTargetGroupsInput).
		Return(describeTargetGroupsOutput, nil)

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	networkMode := ecs.NetworkModeAwsvpc

	taskDefinition := &ecs.TaskDefinition{
		Compatibilities:      []*string{aws.String(ecs.LaunchTypeEc2), aws.String(ecs.LaunchTypeFargate)},
		ContainerDefinitions: containerDefinitions,
		NetworkMode:          &networkMode,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected CreateService
	awsvpcConfig := &ecs.AwsVpcConfiguration{}
	awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
	awsvpcConfig.SetSecurityGroups([]*string{aws.String("sg-test")})
	awsvpcConfig.SetSubnets([]*string{aws.String("subnet-test")})

	networkConfig := &ecs.NetworkConfiguration{}
	networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

	createServiceInput := &ecs.CreateServiceInput{}
	createServiceInput.SetCluster("l0-test-env_id")
	createServiceInput.SetDesiredCount(1)
	createServiceInput.SetLaunchType(ecs.LaunchTypeFargate)
	createServiceInput.SetNetworkConfiguration(networkConfig)
	createServiceInput.SetPlatformVersion(config.DefaultFargatePlatformVersion)
	createServiceInput.SetServiceName("l0-test-svc_id")
	createServiceInput.SetTaskDefinition("dpl_arn")

	loadBalancer := &ecs.LoadBalancer{}
	loadBalancer.SetContainerName("ctn_name")
	loadBalancer.SetContainerPort(80)
	loadBalancer.SetTargetGroupArn("lb_arn")

	loadBalancers := []*ecs.LoadBalancer{loadBalancer}

	createServiceInput.SetLoadBalancers(loadBalancers)

	mockAWS.ECS.EXPECT().
		CreateService(createServiceInput).
		Return(&ecs.CreateServiceOutput{}, nil)

	// define request
	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		Scale:          1,
		ServiceName:    "svc_name",
		Stateful:       false,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result, "svc_id")

	expectedTags := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "load_balancer_id",
			Value:      "lb_id",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestServiceCreate_statefulLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	defer provider.SetEntityIDGenerator("svc_id")()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// Define expected load balancers
	//
	// First we check for an ELB load balancer.
	// This needs to return an ELB load balancer so that we
	// don't then check for an ALB load balancer.
	loadBalancerInput := &elb.DescribeLoadBalancersInput{}
	loadBalancerInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_id")})
	loadBalancerInput.SetPageSize(1)

	listener := &elb.Listener{}
	listener.SetInstancePort(config.DefaultLoadBalancerPort().HostPort)

	listenerDescription := &elb.ListenerDescription{}
	listenerDescription.SetListener(listener)

	listenerDescriptions := []*elb.ListenerDescription{
		listenerDescription,
	}

	loadBalancerDescription := &elb.LoadBalancerDescription{}
	loadBalancerDescription.SetListenerDescriptions(listenerDescriptions)
	loadBalancerDescription.SetLoadBalancerName("l0-test-lb_id")

	loadBalancerDescriptions := []*elb.LoadBalancerDescription{
		loadBalancerDescription,
	}

	loadBalancerOutput := &elb.DescribeLoadBalancersOutput{}
	loadBalancerOutput.SetLoadBalancerDescriptions(loadBalancerDescriptions)

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(loadBalancerInput).
		Return(loadBalancerOutput, nil)

	// define expected ecs.DescribeTaskDefinition
	taskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	taskDefinitionInput.SetTaskDefinition("dpl_arn")

	portMapping := &ecs.PortMapping{}
	portMapping.SetContainerPort(int64(80))

	portMappings := []*ecs.PortMapping{
		portMapping,
	}

	containerDefinition := &ecs.ContainerDefinition{}
	containerDefinition.SetName("ctn_name")
	containerDefinition.SetPortMappings(portMappings)

	containerDefinitions := []*ecs.ContainerDefinition{
		containerDefinition,
	}

	taskDefinition := &ecs.TaskDefinition{
		Compatibilities:      []*string{aws.String(ecs.LaunchTypeEc2)},
		ContainerDefinitions: containerDefinitions,
	}

	taskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	taskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(taskDefinitionInput).
		Return(taskDefinitionOutput, nil)

	// define expected CreateService
	createServiceInput := &ecs.CreateServiceInput{}
	createServiceInput.SetCluster("l0-test-env_id")
	createServiceInput.SetDesiredCount(1)
	createServiceInput.SetLaunchType(ecs.LaunchTypeEc2)
	createServiceInput.SetRole("l0-test-lb_id-lb")
	createServiceInput.SetServiceName("l0-test-svc_id")
	createServiceInput.SetTaskDefinition("dpl_arn")

	loadBalancer := &ecs.LoadBalancer{}
	loadBalancer.SetContainerName("ctn_name")
	loadBalancer.SetContainerPort(80)
	loadBalancer.SetLoadBalancerName("l0-test-lb_id")

	loadBalancers := []*ecs.LoadBalancer{loadBalancer}

	createServiceInput.SetLoadBalancers(loadBalancers)

	mockAWS.ECS.EXPECT().
		CreateService(createServiceInput).
		Return(&ecs.CreateServiceOutput{}, nil)

	// define request
	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		Scale:          1,
		ServiceName:    "svc_name",
		Stateful:       true,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result, "svc_id")

	expectedTags := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "load_balancer_id",
			Value:      "lb_id",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
