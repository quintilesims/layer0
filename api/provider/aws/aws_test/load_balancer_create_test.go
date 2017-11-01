package test_aws

import (
	"strings"
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

	defer provider.SetEntityIDGenerator("lb_id")()

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         false,
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
			Interval:           60,
			Timeout:            60,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
	}

	readSGHelper(mockAWS, "l0-test-env_id-env", "env_sg")
	createSGHelper(t, mockAWS, "l0-test-lb_id-lb", "vpc_id")
	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")

	for _, port := range req.Ports {
		ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
		ingressInput.SetGroupId("lb_sg")
		ingressInput.SetCidrIp("0.0.0.0/0")
		ingressInput.SetIpProtocol("TCP")
		ingressInput.SetFromPort(int64(port.HostPort))
		ingressInput.SetToPort(int64(port.HostPort))

		mockAWS.EC2.EXPECT().
			AuthorizeSecurityGroupIngress(ingressInput).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)
	}

	iamRoleInput := &iam.CreateRoleInput{}
	iamRoleInput.SetRoleName("l0-test-lb_id-lb")
	iamRoleInput.SetAssumeRolePolicyDocument(provider.DEFAULT_ASSUME_ROLE_POLICY)

	mockAWS.IAM.EXPECT().
		CreateRole(iamRoleInput).
		Return(&iam.CreateRoleOutput{}, nil)

	renderedLBRolePolicy, err := provider.RenderLoadBalancerRolePolicy(
		"region",
		"123456789012",
		"l0-test-lb_id",
		provider.DEFAULT_LB_ROLE_POLICY_TEMPLATE)
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

	listeners := make([]*elb.Listener, len(req.Ports))
	for i, port := range req.Ports {
		listener := &elb.Listener{}
		listener.SetProtocol(port.Protocol)
		listener.SetLoadBalancerPort(port.HostPort)
		listener.SetInstancePort(port.ContainerPort)

		if port.CertificateName != "" {
			serverCertificateMetadataList := []*iam.ServerCertificateMetadata{
				&iam.ServerCertificateMetadata{
					Arn: aws.String(port.CertificateName),
					ServerCertificateName: aws.String(port.CertificateName),
				},
			}

			listServerCertificatesOutput := &iam.ListServerCertificatesOutput{}
			listServerCertificatesOutput.SetServerCertificateMetadataList(serverCertificateMetadataList)

			mockAWS.IAM.EXPECT().
				ListServerCertificates(&iam.ListServerCertificatesInput{}).
				Return(listServerCertificatesOutput, nil)

			listener.SetSSLCertificateId(port.CertificateName)
		}

		switch strings.ToLower(port.Protocol) {
		case "http", "https":
			listener.SetInstanceProtocol("http")
		case "tcp", "ssl":
			listener.SetInstanceProtocol("tcp")
		}

		listeners[i] = listener
	}

	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internal")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("priv1"), aws.String("priv2")})
	createLoadBalancerInput.SetListeners(listeners)

	validateFN := func(input *elb.CreateLoadBalancerInput) {
		for i, listener := range input.Listeners {
			assert.Equal(t, listeners[i], listener)
		}
	}

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Do(validateFN).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := &elb.HealthCheck{}
	healthCheck.SetTarget(req.HealthCheck.Target)
	healthCheck.SetInterval(int64(req.HealthCheck.Interval))
	healthCheck.SetTimeout(int64(req.HealthCheck.Timeout))
	healthCheck.SetHealthyThreshold(int64(req.HealthCheck.HealthyThreshold))
	healthCheck.SetUnhealthyThreshold(int64(req.HealthCheck.UnhealthyThreshold))

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
	mockConfig.EXPECT().PublicSubnets().Return([]string{"pub1", "pub2"}).AnyTimes()

	defer provider.SetEntityIDGenerator("lb_id")()

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         true,
		Ports: []models.Port{
			{
				ContainerPort: 80,
				HostPort:      80,
				Protocol:      "tcp",
			},
		},
		HealthCheck: models.HealthCheck{
			Target:             "TCP:80",
			Interval:           30,
			Timeout:            5,
			HealthyThreshold:   2,
			UnhealthyThreshold: 2,
		},
	}

	readSGHelper(mockAWS, "l0-test-env_id-env", "env_sg")
	createSGHelper(t, mockAWS, "l0-test-lb_id-lb", "vpc_id")
	readSGHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg")

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId("lb_sg")
	ingressInput.SetCidrIp("0.0.0.0/0")
	ingressInput.SetIpProtocol("TCP")
	ingressInput.SetFromPort(int64(80))
	ingressInput.SetToPort(int64(80))

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	mockAWS.IAM.EXPECT().
		CreateRole(gomock.Any()).
		Return(&iam.CreateRoleOutput{}, nil)

	mockAWS.IAM.EXPECT().
		PutRolePolicy(gomock.Any()).
		Return(&iam.PutRolePolicyOutput{}, nil)

	listeners := make([]*elb.Listener, len(req.Ports))
	listener := &elb.Listener{}
	listener.SetProtocol("tcp")
	listener.SetLoadBalancerPort(80)
	listener.SetInstancePort(80)
	listener.SetInstanceProtocol("tcp")
	listeners[0] = listener

	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internet-facing")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("pub1"), aws.String("pub2")})
	createLoadBalancerInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := &elb.HealthCheck{}
	healthCheck.SetTarget("TCP:80")
	healthCheck.SetInterval(int64(30))
	healthCheck.SetTimeout(int64(5))
	healthCheck.SetHealthyThreshold(int64(2))
	healthCheck.SetUnhealthyThreshold(int64(2))

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
