package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

func TestLoadBalancerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
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
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
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

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeLoadBalancersInput := &elb.DescribeLoadBalancersInput{}
	describeLoadBalancersInput.SetLoadBalancerNames([]*string{aws.String("l0-test-lb_id")})
	describeLoadBalancersInput.SetPageSize(1)

	healthCheck := models.HealthCheck{
		Target:             "HTTPS:443/path/to/site",
		Interval:           10,
		Timeout:            6,
		HealthyThreshold:   3,
		UnhealthyThreshold: 2,
	}

	ports := []models.Port{
		models.Port{
			CertificateName: "cert",
			ContainerPort:   4444,
			HostPort:        443,
			Protocol:        "https",
		},
		models.Port{
			ContainerPort: 88,
			HostPort:      80,
			Protocol:      "tcp",
		},
	}

	certificateARN := "arn:aws:iam::123456789012:server-certificate/cert"
	serverCertificateMetadata := &iam.ServerCertificateMetadata{}
	serverCertificateMetadata.SetArn(certificateARN)
	serverCertificateMetadata.SetServerCertificateName("cert")

	listener1 := listenerHelper(ports[0])
	listener1.SetSSLCertificateId(certificateARN)
	listenerDescription1 := &elb.ListenerDescription{}
	listenerDescription1.SetListener(listener1)
	listenerDescription2 := &elb.ListenerDescription{}
	listenerDescription2.SetListener(listenerHelper(ports[1]))
	listenerDescriptions := []*elb.ListenerDescription{listenerDescription1, listenerDescription2}

	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_id")
	lb.SetHealthCheck(healthCheckHelper(&healthCheck))
	lb.SetListenerDescriptions(listenerDescriptions)

	describeLoadBalancersOutput := &elb.DescribeLoadBalancersOutput{}
	describeLoadBalancersOutput.SetLoadBalancerDescriptions([]*elb.LoadBalancerDescription{lb})

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancers(describeLoadBalancersInput).
		Return(describeLoadBalancersOutput, nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("lb_id")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		ServiceID:        "svc_id",
		ServiceName:      "svc_name",
		HealthCheck:      healthCheck,
		Ports:            ports,
	}

	assert.Equal(t, expected, result)
}
