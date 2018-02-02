package test_aws

import (
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().LinuxAMI().Return("lx_ami").AnyTimes()
	mockConfig.EXPECT().WindowsAMI().Return("win_ami").AnyTimes()
	mockConfig.EXPECT().S3Bucket().Return("bucket").AnyTimes()
	mockConfig.EXPECT().VPC().Return("vpc_id").AnyTimes()
	mockConfig.EXPECT().InstanceProfile().Return("profile").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"priv1", "priv2"}).AnyTimes()
	mockConfig.EXPECT().SSHKeyPair().Return("keypair").AnyTimes()

	defer provider.SetEntityIDGenerator("env_id")()

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "env_name",
		EnvironmentType:  "static",
		InstanceType:     "t2.small",
		UserDataTemplate: []byte("some user data"),
		AMIID:            "some ami",
		Scale:            2,
		OperatingSystem:  "windows",
	}

	// an environment's security group name is <fq environment id>-env
	createSGHelper(t, mockAWS, "l0-test-env_id-env", "vpc_id")
	readSGHelper(mockAWS, "l0-test-env_id-env", "sg_id")

	// ensure we add a self-ingress rule to the environment's security group
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId("sg_id")
	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId("sg_id")
	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	// an environent's launch configuration name is the fq environment id
	createLCInput := &autoscaling.CreateLaunchConfigurationInput{}
	createLCInput.SetLaunchConfigurationName("l0-test-env_id")
	createLCInput.SetSecurityGroups([]*string{aws.String("sg_id")})
	createLCInput.SetInstanceType("t2.small")
	createLCInput.SetIamInstanceProfile("profile")
	createLCInput.SetImageId("some ami")
	createLCInput.SetKeyName("keypair")
	base64UserData := base64.StdEncoding.EncodeToString([]byte("some user data"))
	createLCInput.SetUserData(base64UserData)

	mockAWS.AutoScaling.EXPECT().
		CreateLaunchConfiguration(createLCInput).
		Return(&autoscaling.CreateLaunchConfigurationOutput{}, nil)

	// an environment's autoscaling group name is the fq environment id
	tag := &autoscaling.Tag{}
	tag.SetKey("Name")
	tag.SetValue("l0-test-env_id")
	tag.SetPropagateAtLaunch(true)

	createASGInput := &autoscaling.CreateAutoScalingGroupInput{}
	createASGInput.SetAutoScalingGroupName("l0-test-env_id")
	createASGInput.SetLaunchConfigurationName("l0-test-env_id")
	createASGInput.SetVPCZoneIdentifier("priv1,priv2")
	createASGInput.SetMinSize(2)
	createASGInput.SetMaxSize(2)
	createASGInput.SetTags([]*autoscaling.Tag{tag})

	mockAWS.AutoScaling.EXPECT().
		CreateAutoScalingGroup(createASGInput).
		Return(&autoscaling.CreateAutoScalingGroupOutput{}, nil)

	// an environment's cluster name is the fq environment id
	createClusterInput := &ecs.CreateClusterInput{}
	createClusterInput.SetClusterName("l0-test-env_id")

	mockAWS.ECS.EXPECT().
		CreateCluster(createClusterInput).
		Return(&ecs.CreateClusterOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", result)

	expectedTags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "type",
			Value:      "static",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "windows",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestEnvironmentCreateDynamic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().LinuxAMI().Return("lx_ami").AnyTimes()
	mockConfig.EXPECT().WindowsAMI().Return("win_ami").AnyTimes()
	mockConfig.EXPECT().S3Bucket().Return("bucket").AnyTimes()
	mockConfig.EXPECT().VPC().Return("vpc_id").AnyTimes()
	mockConfig.EXPECT().InstanceProfile().Return("profile").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"priv1", "priv2"}).AnyTimes()
	mockConfig.EXPECT().SSHKeyPair().Return("keypair").AnyTimes()

	defer provider.SetEntityIDGenerator("env_id")()

	req := models.CreateEnvironmentRequest{
		EnvironmentName: "env_name",
		EnvironmentType: "dynamic",
		OperatingSystem: "linux",
	}

	// an environment's security group name is <fq environment id>-env
	createSGHelper(t, mockAWS, "l0-test-env_id-env", "vpc_id")
	readSGHelper(mockAWS, "l0-test-env_id-env", "sg_id")

	// ensure we add a self-ingress rule to the environment's security group
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId("sg_id")
	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId("sg_id")
	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	// an environment's cluster name is the fq environment id
	createClusterInput := &ecs.CreateClusterInput{}
	createClusterInput.SetClusterName("l0-test-env_id")

	mockAWS.ECS.EXPECT().
		CreateCluster(createClusterInput).
		Return(&ecs.CreateClusterOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", result)

	expectedTags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "type",
			Value:      "dynamic",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "linux",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestEnvironmentCreateDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	// todo: setup helper for config
	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().LinuxAMI().Return("lx_ami").AnyTimes()
	mockConfig.EXPECT().WindowsAMI().Return("win_ami").AnyTimes()
	mockConfig.EXPECT().S3Bucket().Return("bucket").AnyTimes()
	mockConfig.EXPECT().VPC().Return("vpc_id").AnyTimes()
	mockConfig.EXPECT().InstanceProfile().Return("profile").AnyTimes()
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"priv1", "priv2"}).AnyTimes()
	mockConfig.EXPECT().SSHKeyPair().Return("keypair").AnyTimes()

	defer provider.SetEntityIDGenerator("env_id")()

	req := models.CreateEnvironmentRequest{
		EnvironmentName: "env_name",
	}

	// using create/read helpers instead of gomock.Any() for readability
	createSGHelper(t, mockAWS, "l0-test-env_id-env", "vpc_id")
	readSGHelper(mockAWS, "l0-test-env_id-env", "sg_id")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(gomock.Any()).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	// ensure we pass the default instance type, ami id, and user data to the launch configuration
	renderedUserData, err := provider.RenderUserData(
		"l0-test-env_id",
		"bucket",
		[]byte(provider.DefaultLinuxUserdataTemplate))
	if err != nil {
		t.Fatal(err)
	}

	validateCreateLCInput := func(input *autoscaling.CreateLaunchConfigurationInput) {
		assert.Equal(t, config.DefaultEnvironmentInstanceType, aws.StringValue(input.InstanceType))
		assert.Equal(t, "lx_ami", aws.StringValue(input.ImageId))
		assert.Equal(t, renderedUserData, aws.StringValue(input.UserData))
	}

	mockAWS.AutoScaling.EXPECT().
		CreateLaunchConfiguration(gomock.Any()).
		Do(validateCreateLCInput).
		Return(&autoscaling.CreateLaunchConfigurationOutput{}, nil)

	mockAWS.AutoScaling.EXPECT().
		CreateAutoScalingGroup(gomock.Any()).
		Return(&autoscaling.CreateAutoScalingGroupOutput{}, nil)

	mockAWS.ECS.EXPECT().
		CreateCluster(gomock.Any()).
		Return(&ecs.CreateClusterOutput{}, nil)

	target := provider.NewEnvironmentProvider(mockAWS.Client(), tagStore, mockConfig)
	if _, err := target.Create(req); err != nil {
		t.Fatal(err)
	}
}
