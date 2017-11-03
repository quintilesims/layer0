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
			CertificateName: "cert",
			ContainerPort:   8080,
			HostPort:        8088,
			Protocol:        "http",
		},
		models.Port{
			CertificateName: "cert",
			ContainerPort:   4444,
			HostPort:        444,
			Protocol:        "https",
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
		LoadBalancerID: "lb_name",
		Ports:          &requestPorts,
		HealthCheck:    requestHealthCheck,
	}

	healthCheck := healthCheckHelper(requestHealthCheck)

	configureHealthCheckInput := &elb.ConfigureHealthCheckInput{}
	configureHealthCheckInput.SetLoadBalancerName("l0-test-lb_name")
	configureHealthCheckInput.SetHealthCheck(healthCheck)

	configureHealthCheckOutput := &elb.ConfigureHealthCheckOutput{}
	configureHealthCheckOutput.SetHealthCheck(healthCheck)

	mockAWS.ELB.EXPECT().
		ConfigureHealthCheck(configureHealthCheckInput).
		Return(configureHealthCheckOutput, nil)

	serverCertificateMetadata := &iam.ServerCertificateMetadata{}
	serverCertificateMetadata.SetArn("cert")
	serverCertificateMetadata.SetServerCertificateName("cert")

	serverCertificateMetadataList1 := []*iam.ServerCertificateMetadata{serverCertificateMetadata}
	listServerCertificatesOutput1 := &iam.ListServerCertificatesOutput{}
	listServerCertificatesOutput1.SetServerCertificateMetadataList(serverCertificateMetadataList1)

	mockAWS.IAM.EXPECT().
		ListServerCertificates(&iam.ListServerCertificatesInput{}).
		Return(listServerCertificatesOutput1, nil)

	serverCertificateMetadataList2 := []*iam.ServerCertificateMetadata{serverCertificateMetadata}
	listServerCertificatesOutput2 := &iam.ListServerCertificatesOutput{}
	listServerCertificatesOutput2.SetServerCertificateMetadataList(serverCertificateMetadataList2)

	mockAWS.IAM.EXPECT().
		ListServerCertificates(&iam.ListServerCertificatesInput{}).
		Return(listServerCertificatesOutput2, nil)

	readSGHelper(mockAWS, "l0-test-lb_name-lb", "lb_sg")

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_name")})
	describeLoadBalancersInput.SetPageSize(1)

	healthCheck = healthCheckHelper(nil)
	listener := listenerHelper(nil)
	listenerDescription := &elb.ListenerDescription{}
	listenerDescription.SetListener(listener)

	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_name")
	lb.SetHealthCheck(healthCheck)
	lb.SetListenerDescriptions([]*elb.ListenerDescription{listenerDescription})

	describeLoadBalancersOutput := &elb.DescribeLoadBalancersOutput{}
	describeLoadBalancersOutput.SetLoadBalancerDescriptions([]*elb.LoadBalancerDescription{lb})

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(describeLoadBalancersOutput, nil)

	ingressInput := &ec2.RevokeSecurityGroupIngressInput{}
	ingressInput.SetGroupId("lb_sg")
	ingressInput.SetCidrIp("0.0.0.0/0")
	ingressInput.SetIpProtocol("TCP")
	ingressInput.SetFromPort(80)
	ingressInput.SetToPort(80)

	mockAWS.EC2.EXPECT().
		RevokeSecurityGroupIngress(ingressInput).
		Return(&ec2.RevokeSecurityGroupIngressOutput{}, nil)

	port := int64(80)
	deleteLoadBalancerListenersInput := &elb.DeleteLoadBalancerListenersInput{}
	deleteLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	deleteLoadBalancerListenersInput.SetLoadBalancerPorts([]*int64{&port})

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancerListeners(deleteLoadBalancerListenersInput).
		Return(&elb.DeleteLoadBalancerListenersOutput{}, nil)

	listener1 := listenerHelper(&requestPorts[0])
	listener2 := listenerHelper(&requestPorts[1])
	listeners := []*elb.Listener{listener1, listener2}

	createLoadBalancerListenersInput := &elb.CreateLoadBalancerListenersInput{}
	createLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	createLoadBalancerListenersInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancerListeners(createLoadBalancerListenersInput).
		Return(&elb.CreateLoadBalancerListenersOutput{}, nil)

	ingressInput1 := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput1.SetGroupId("lb_sg")
	ingressInput1.SetCidrIp("0.0.0.0/0")
	ingressInput1.SetIpProtocol("TCP")
	ingressInput1.SetFromPort(int64(8088))
	ingressInput1.SetToPort(int64(8088))

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput1).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	ingressInput2 := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput2.SetGroupId("lb_sg")
	ingressInput2.SetCidrIp("0.0.0.0/0")
	ingressInput2.SetIpProtocol("TCP")
	ingressInput2.SetFromPort(int64(444))
	ingressInput2.SetToPort(int64(444))

	mockAWS.EC2.EXPECT().
		AuthorizeSecurityGroupIngress(ingressInput2).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update(req); err != nil {
		t.Fatal(err)
	}
}
