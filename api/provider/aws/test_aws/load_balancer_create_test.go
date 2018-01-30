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
	"github.com/quintilesims/layer0/common/config"
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
				Certificate:   "cert",
				ContainerPort: 88,
				HostPort:      80,
				Protocol:      "http",
			},
			{
				Certificate:   "cert",
				ContainerPort: 4444,
				HostPort:      443,
				Protocol:      "https",
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

	authorizeIngressInput := authorizeSGIngressHelper(req.Ports[0])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	authorizeIngressInput = authorizeSGIngressHelper(req.Ports[1])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
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

	certificateARN := "arn:aws:iam::123456789012:server-certificate/cert"
	serverCertificateMetadata := &iam.ServerCertificateMetadata{}
	serverCertificateMetadata.SetArn(certificateARN)
	serverCertificateMetadata.SetServerCertificateName("cert")
	serverCertificateMetadataList := []*iam.ServerCertificateMetadata{serverCertificateMetadata}

	listServerCertificatesOutput := &iam.ListServerCertificatesOutput{}
	listServerCertificatesOutput.SetServerCertificateMetadataList(serverCertificateMetadataList)

	mockAWS.IAM.EXPECT().
		ListServerCertificates(&iam.ListServerCertificatesInput{}).
		Return(listServerCertificatesOutput, nil).
		AnyTimes()

	listener1 := listenerHelper(req.Ports[0])
	listener1.SetSSLCertificateId(certificateARN)
	listener2 := listenerHelper(req.Ports[1])
	listener2.SetSSLCertificateId(certificateARN)
	listeners := []*elb.Listener{listener1, listener2}

	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internet-facing")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("pub1"), aws.String("pub2")})
	createLoadBalancerInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := healthCheckHelper(&req.HealthCheck)
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

	listeners := []*elb.Listener{listenerHelper(config.DefaultLoadBalancerPort)}
	createLoadBalancerInput := &elb.CreateLoadBalancerInput{}
	createLoadBalancerInput.SetLoadBalancerName("l0-test-lb_id")
	createLoadBalancerInput.SetScheme("internal")
	createLoadBalancerInput.SetSecurityGroups([]*string{aws.String("env_sg"), aws.String("lb_sg")})
	createLoadBalancerInput.SetSubnets([]*string{aws.String("priv1"), aws.String("priv2")})
	createLoadBalancerInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancer(createLoadBalancerInput).
		Return(&elb.CreateLoadBalancerOutput{}, nil)

	healthCheck := healthCheckHelper(&config.DefaultLoadBalancerHealthCheck)
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
