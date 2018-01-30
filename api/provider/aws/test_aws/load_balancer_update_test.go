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
)

func TestLoadBalancerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	requestPorts := []models.Port{
		models.Port{
			Certificate:   "cert",
			ContainerPort: 8080,
			HostPort:      8088,
			Protocol:      "http",
		},
		models.Port{
			Certificate:   "cert",
			ContainerPort: 4444,
			HostPort:      444,
			Protocol:      "https",
		},
	}

	requestHealthCheck := &models.HealthCheck{
		Target:             "HTTPS:444/path/to/site",
		Interval:           15,
		Timeout:            10,
		HealthyThreshold:   5,
		UnhealthyThreshold: 4,
	}

	req := models.UpdateLoadBalancerRequest{
		Ports:       &requestPorts,
		HealthCheck: requestHealthCheck,
	}

	configureHealthCheckInput := &elb.ConfigureHealthCheckInput{}
	configureHealthCheckInput.SetLoadBalancerName("l0-test-lb_name")
	configureHealthCheckInput.SetHealthCheck(healthCheckHelper(requestHealthCheck))

	configureHealthCheckOutput := &elb.ConfigureHealthCheckOutput{}
	configureHealthCheckOutput.SetHealthCheck(healthCheckHelper(requestHealthCheck))

	mockAWS.ELB.EXPECT().
		ConfigureHealthCheck(configureHealthCheckInput).
		Return(configureHealthCheckOutput, nil)

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

	readSGHelper(mockAWS, "l0-test-lb_name-lb", "lb_sg")
	listenerDescription := &elb.ListenerDescription{}
	listenerDescription.SetListener(listenerHelper(config.DefaultLoadBalancerPort))

	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_name")
	lb.SetHealthCheck(healthCheckHelper(&config.DefaultLoadBalancerHealthCheck))
	lb.SetListenerDescriptions([]*elb.ListenerDescription{listenerDescription})

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_name")})
	describeLoadBalancersInput.SetPageSize(1)

	describeLoadBalancersOutput := &elb.DescribeLoadBalancersOutput{}
	describeLoadBalancersOutput.SetLoadBalancerDescriptions([]*elb.LoadBalancerDescription{lb})

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(describeLoadBalancersOutput, nil)

	revokeIngressInput := revokeSGIngressHelper(config.DefaultLoadBalancerPort)
	revokeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		RevokeSecurityGroupIngress(revokeIngressInput).
		Return(&ec2.RevokeSecurityGroupIngressOutput{}, nil)

	port := int64(config.DefaultLoadBalancerPort.HostPort)
	deleteLoadBalancerListenersInput := &elb.DeleteLoadBalancerListenersInput{}
	deleteLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	deleteLoadBalancerListenersInput.SetLoadBalancerPorts([]*int64{&port})

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancerListeners(deleteLoadBalancerListenersInput).
		Return(&elb.DeleteLoadBalancerListenersOutput{}, nil)

	listener1 := listenerHelper(requestPorts[0])
	listener1.SetSSLCertificateId(certificateARN)
	listener2 := listenerHelper(requestPorts[1])
	listener2.SetSSLCertificateId(certificateARN)

	createLoadBalancerListenersInput := &elb.CreateLoadBalancerListenersInput{}
	createLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	createLoadBalancerListenersInput.SetListeners([]*elb.Listener{listener1, listener2})

	mockAWS.ELB.EXPECT().
		CreateLoadBalancerListeners(createLoadBalancerListenersInput).
		Return(&elb.CreateLoadBalancerListenersOutput{}, nil)

	authorizeIngressInput := authorizeSGIngressHelper(requestPorts[0])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	authorizeIngressInput = authorizeSGIngressHelper(requestPorts[1])
	authorizeIngressInput.SetGroupId("lb_sg")

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(authorizeIngressInput).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update("lb_name", req); err != nil {
		t.Fatal(err)
	}
}
