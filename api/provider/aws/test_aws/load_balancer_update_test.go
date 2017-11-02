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

	requestPorts := &[]models.Port{
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
		Ports:          requestPorts,
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

	listeners := make([]*elb.Listener, len(*requestPorts))
	for i, port := range *requestPorts {
		listener := listenerHelper(&port)
		if port.CertificateName != "" {
			serverCertificateMetadata := &iam.ServerCertificateMetadata{}
			serverCertificateMetadata.SetArn(port.CertificateName)
			serverCertificateMetadata.SetServerCertificateName(port.CertificateName)
			serverCertificateMetadataList := []*iam.ServerCertificateMetadata{serverCertificateMetadata}

			listServerCertificatesOutput := &iam.ListServerCertificatesOutput{}
			listServerCertificatesOutput.SetServerCertificateMetadataList(serverCertificateMetadataList)

			mockAWS.IAM.EXPECT().
				ListServerCertificates(&iam.ListServerCertificatesInput{}).
				Return(listServerCertificatesOutput, nil)
		}

		listeners[i] = listener
	}

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

	portNumbers := make([]int64, len(lb.ListenerDescriptions))
	for i, listenerDescription := range lb.ListenerDescriptions {
		ingressInput := &ec2.RevokeSecurityGroupIngressInput{}
		ingressInput.SetGroupId("lb_sg")
		ingressInput.SetCidrIp("0.0.0.0/0")
		ingressInput.SetIpProtocol("TCP")
		ingressInput.SetFromPort(aws.Int64Value(listenerDescription.Listener.LoadBalancerPort))
		ingressInput.SetToPort(aws.Int64Value(listenerDescription.Listener.LoadBalancerPort))

		mockAWS.EC2.EXPECT().
			RevokeSecurityGroupIngress(ingressInput).
			Return(&ec2.RevokeSecurityGroupIngressOutput{}, nil)

		portNumber := aws.Int64Value(listenerDescription.Listener.LoadBalancerPort)
		portNumbers[i] = portNumber
	}

	deleteLoadBalancerListenersInput := &elb.DeleteLoadBalancerListenersInput{}
	deleteLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")

	ports := make([]*int64, len(portNumbers))
	for i, p := range portNumbers {
		ports[i] = aws.Int64(p)
	}

	deleteLoadBalancerListenersInput.SetLoadBalancerPorts(ports)

	mockAWS.ELB.EXPECT().
		DeleteLoadBalancerListeners(deleteLoadBalancerListenersInput).
		Return(&elb.DeleteLoadBalancerListenersOutput{}, nil)

	createLoadBalancerListenersInput := &elb.CreateLoadBalancerListenersInput{}
	createLoadBalancerListenersInput.SetLoadBalancerName("l0-test-lb_name")
	createLoadBalancerListenersInput.SetListeners(listeners)

	mockAWS.ELB.EXPECT().
		CreateLoadBalancerListeners(createLoadBalancerListenersInput).
		Return(&elb.CreateLoadBalancerListenersOutput{}, nil)

	for _, listener := range listeners {
		loadBalancerListenerPort := aws.Int64Value(listener.LoadBalancerPort)

		ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
		ingressInput.SetGroupId("lb_sg")
		ingressInput.SetCidrIp("0.0.0.0/0")
		ingressInput.SetIpProtocol("TCP")
		ingressInput.SetFromPort(loadBalancerListenerPort)
		ingressInput.SetToPort(loadBalancerListenerPort)

		mockAWS.EC2.EXPECT().
			AuthorizeSecurityGroupIngress(ingressInput).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil)
	}

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Update(req); err != nil {
		t.Fatal(err)
	}
}
