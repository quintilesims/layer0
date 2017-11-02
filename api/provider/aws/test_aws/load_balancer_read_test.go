package test_aws

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
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

	elbHealthCheck := &elb.HealthCheck{}
	elbHealthCheck.SetTarget(healthCheck.Target)
	elbHealthCheck.SetInterval(int64(healthCheck.Interval))
	elbHealthCheck.SetTimeout(int64(healthCheck.Timeout))
	elbHealthCheck.SetHealthyThreshold(int64(healthCheck.HealthyThreshold))
	elbHealthCheck.SetUnhealthyThreshold(int64(healthCheck.UnhealthyThreshold))

	listenerDescriptions := make([]*elb.ListenerDescription, len(ports))
	for i, port := range ports {
		elbListener := &elb.Listener{}
		elbListener.SetProtocol(port.Protocol)
		elbListener.SetLoadBalancerPort(port.HostPort)
		elbListener.SetInstancePort(port.ContainerPort)

		switch strings.ToLower(port.Protocol) {
		case "http", "https":
			elbListener.SetInstanceProtocol("http")
		case "tcp", "ssl":
			elbListener.SetInstanceProtocol("tcp")
		}

		if port.CertificateName != "" {
			elbListener.SetSSLCertificateId("cert")
		}

		listenerDescription := &elb.ListenerDescription{}
		listenerDescription.SetListener(elbListener)

		listenerDescriptions[i] = listenerDescription
	}

	lb := &elb.LoadBalancerDescription{}
	lb.SetLoadBalancerName("l0-test-lb_id")
	lb.SetHealthCheck(elbHealthCheck)
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
		HealthCheck:      healthCheck,
		Ports:            ports,
	}

	assert.Equal(t, expected, result)
}
