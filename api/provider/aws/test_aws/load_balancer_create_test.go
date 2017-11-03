package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().VPC().Return("vpc_id").AnyTimes()
	mockConfig.EXPECT().Region().Return("region").AnyTimes()
	mockConfig.EXPECT().AccountID().Return("123456789012")
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"priv1", "priv2"}).AnyTimes()
	mockConfig.EXPECT().PublicSubnets().Return([]string{"pub1", "pub2"}).AnyTimes()

	defer provider.SetEntityIDGenerator("lb_id")()

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports: []models.Port{
			{
				CertificateName: "cert",
				ContainerPort:   88,
				HostPort:        80,
				Protocol:        "http",
			},
			{
				CertificateName: "cert",
				ContainerPort:   4444,
				HostPort:        443,
				Protocol:        "https",
			},
		},
		HealthCheck: models.HealthCheck{
			Target:             "HTTPS:443/path/to/site",
			Interval:           20,
			Timeout:            15,
			HealthyThreshold:   4,
			UnhealthyThreshold: 3,
		},
	}

	readSGHelper(mockAWS, "l0-test-env_id-env", "env_sg")
	createSGHelper(t, mockAWS, "l0-test-lb_id-lb", "vpc_id")
	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")

	ingressInput1 := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput1.SetGroupId("lb_sg")
	ingressInput1.SetCidrIp("0.0.0.0/0")
	ingressInput1.SetIpProtocol("TCP")
	ingressInput1.SetFromPort(int64(80))
	ingressInput1.SetToPort(int64(80))

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput1).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	ingressInput2 := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput2.SetGroupId("lb_sg")
	ingressInput2.SetCidrIp("0.0.0.0/0")
	ingressInput2.SetIpProtocol("TCP")
	ingressInput2.SetFromPort(int64(443))
	ingressInput2.SetToPort(int64(443))

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput2).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	iamRoleInput := &iam.CreateRoleInput{}
	iamRoleInput.SetRoleName("l0-test-lb_id-lb")
	iamRoleInput.SetAssumeRolePolicyDocument(provider.DefaultAssumeRolePolicy)

	mockAWS.IAM.EXPECT().
		CreateRole(iamRoleInput).
		Return(&iam.CreateRoleOutput{}, nil)

	renderedLBRolePolicy, err := provider.RenderLoadBalancerRolePolicy(
		"region",
		"123456789012",
		"l0-test-lb_id",
		provider.DefaultLBRolePolicyTemplate)
	if err != nil {
		t.Fatal(err)
	}

	putIAMPolicyInput := &iam.PutRolePolicyInput{}
	putIAMPolicyInput.SetPolicyName("l0-test-lb_id-lb")
	putIAMPolicyInput.SetRoleName("l0-test-lb_id-lb")
	putIAMPolicyInput.SetPolicyDocument(renderedLBRolePolicy)

	mockAWS.IAM.EXPECT().
		PutRolePolicy(putIAMPolicyInput).
		Return(&iam.PutRolePolicyOutput{}, nil)

	serverCertificateMetadata := &iam.ServerCertificateMetadata{}
	serverCertificateMetadata.SetArn("cert")
	serverCertificateMetadata.SetServerCertificateName("cert")
	serverCertificateMetadataList := []*iam.ServerCertificateMetadata{serverCertificateMetadata}

	listServerCertificatesOutput := &iam.ListServerCertificatesOutput{}
	listServerCertificatesOutput.SetServerCertificateMetadataList(serverCertificateMetadataList)

	mockAWS.IAM.EXPECT().
		ListServerCertificates(&iam.ListServerCertificatesInput{}).
		Return(listServerCertificatesOutput, nil).
		AnyTimes()

	listener1 := &elb.Listener{}
	listener1.SetProtocol("http")
	listener1.SetLoadBalancerPort(80)
	listener1.SetInstancePort(88)
	listener1.SetSSLCertificateId("cert")
	listener1.SetInstanceProtocol("http")

	listener2 := &elb.Listener{}
	listener2.SetProtocol("https")
	listener2.SetLoadBalancerPort(443)
	listener2.SetInstancePort(4444)
	listener2.SetSSLCertificateId("cert")
	listener2.SetInstanceProtocol("http")

	listeners := []*elb.Listener{listener1, listener2}
	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internet-facing")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("pub1"), aws.String("pub2")})
	createLoadBalancerInput.SetListeners(listeners)

	validateFN := func(input *elb.CreateLoadBalancerInput) {
		assert.NotNil(t, input.Listeners)
		assert.NotNil(t, input.LoadBalancerName)

		for _, listener := range input.Listeners {
			err := listener.Validate()
			assert.NoError(t, err)
		}
	}

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Do(validateFN).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := &elb.HealthCheck{}
	healthCheck.SetTarget("HTTPS:443/path/to/site")
	healthCheck.SetInterval(int64(20))
	healthCheck.SetTimeout(int64(15))
	healthCheck.SetHealthyThreshold(int64(4))
	healthCheck.SetUnhealthyThreshold(int64(3))

	configureHealthCheckInput := &elb.ConfigureHealthCheckInput{}
	configureHealthCheckInput.SetLoadBalancerName("l0-test-lb_id")
	configureHealthCheckInput.SetHealthCheck(healthCheck)

	mockAWS.ELB.EXPECT().
		ConfigureHealthCheck(configureHealthCheckInput).
		Return(&elb.ConfigureHealthCheckOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "lb_id", result)

	expectedTags := models.Tags{
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestLoadBalancerCreateDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	mockConfig.EXPECT().VPC().Return("vpc_id").AnyTimes()
	mockConfig.EXPECT().Region().Return("region").AnyTimes()
	mockConfig.EXPECT().AccountID().Return("123456789012")
	mockConfig.EXPECT().PrivateSubnets().Return([]string{"priv1", "priv2"}).AnyTimes()

	defer provider.SetEntityIDGenerator("lb_id")()

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		Ports:            []models.Port{},
		HealthCheck:      models.HealthCheck{},
	}

	readSGHelper(mockAWS, "l0-test-env_id-env", "env_sg")
	createSGHelper(t, mockAWS, "l0-test-lb_id-lb", "vpc_id")
	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(gomock.Any()).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	mockAWS.IAM.EXPECT().
		CreateRole(gomock.Any()).
		Return(&iam.CreateRoleOutput{}, nil)

	mockAWS.IAM.EXPECT().
		PutRolePolicy(gomock.Any()).
		Return(&iam.PutRolePolicyOutput{}, nil)

	listeners := []*elb.Listener{listenerHelper(nil)}
	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internal")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("priv1"), aws.String("priv2")})
	createLoadBalancerInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := healthCheckHelper(nil)
	configureHealthCheckInput := &elb.ConfigureHealthCheckInput{}
	configureHealthCheckInput.SetLoadBalancerName("l0-test-lb_id")
	configureHealthCheckInput.SetHealthCheck(healthCheck)

	mockAWS.ELB.EXPECT().
		ConfigureHealthCheck(configureHealthCheckInput).
		Return(&elb.ConfigureHealthCheckOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "lb_id", result)
}
