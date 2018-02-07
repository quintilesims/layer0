package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceLoadBalancerCreateRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	ports := []models.Port{
		{
			HostPort:      80,
			ContainerPort: 8080,
			Protocol:      "tcp",
		},
		{
			CertificateARN: "arn:aws:iam::12345:server-certificate/crt_name",
			HostPort:       9000,
			ContainerPort:  9090,
			Protocol:       "https",
		},
		{
			CertificateARN: "arn:aws:acm::12345:certificate/crt_name",
			HostPort:       10000,
			ContainerPort:  10101,
			Protocol:       "https",
		},
	}

	healthCheck := models.HealthCheck{
		Target:             "tcp:80",
		Interval:           10,
		Timeout:            15,
		HealthyThreshold:   5,
		UnhealthyThreshold: 2,
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		IsPublic:         false,
		Ports:            ports,
		HealthCheck:      healthCheck,
	}

	mockClient.EXPECT().
		CreateLoadBalancer(req).
		Return("lb_id", nil)

	loadBalancer := &models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		Ports:            ports,
		HealthCheck:      healthCheck,
		URL:              "some_url",
	}

	mockClient.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(loadBalancer, nil)

	loadBalancerResource := Provider().(*schema.Provider).ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":         "lb_name",
		"environment":  "env_id",
		"private":      true,
		"port":         flattenPorts(ports),
		"health_check": flattenHealthCheck(healthCheck),
	})

	if err := resourceLayer0LoadBalancerCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "lb_id", d.Id())
	assert.Equal(t, "lb_name", d.Get("name"))
	assert.Equal(t, "env_id", d.Get("environment"))
	assert.Equal(t, true, d.Get("private"))
	assert.Equal(t, ports, expandPorts(d.Get("port").(*schema.Set).List()))
	assert.Equal(t, healthCheck, expandHealthCheck(d.Get("health_check")))
	assert.Equal(t, "some_url", d.Get("url"))
}

func TestResourceLoadBalancerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	mockClient.EXPECT().
		DeleteLoadBalancer("lb_id").
		Return(nil)

	loadBalancerResource := Provider().(*schema.Provider).ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{})
	d.SetId("lb_id")

	if err := resourceLayer0LoadBalancerDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
