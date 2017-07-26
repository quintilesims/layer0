package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestLoadBalancerCreate_defaults(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	ports := []models.Port{
		{
			ContainerPort: 80,
			HostPort:      80,
			Protocol:      "http",
		},
	}

	mockClient.EXPECT().
		CreateLoadBalancer("test-lb", "test-env", models.HealthCheck{"TCP:80", 30, 5, 2, 2}, ports, true).
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	mockClient.EXPECT().
		GetLoadBalancer("lbid").
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
		"port":        flattenPorts(ports),
	})

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerCreate_specifyPorts(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	ports := []models.Port{
		{
			ContainerPort: 80,
			HostPort:      80,
			Protocol:      "http",
		},
		{
			CertificateName: "certname",
			ContainerPort:   8080,
			HostPort:        8080,
			Protocol:        "https",
		},
	}

	mockClient.EXPECT().
		CreateLoadBalancer("test-lb", "test-env", models.HealthCheck{"TCP:80", 30, 5, 2, 2}, ports, false).
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	mockClient.EXPECT().
		GetLoadBalancer("lbid").
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
		"port":        flattenPorts(ports),
		"private":     true,
	})

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerCreate_specifyHealthCheck(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateLoadBalancer("test-lb", "test-env", models.HealthCheck{"HTTP:80/admin/healthcheck", 25, 10, 4, 3}, []models.Port{}, true).
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	mockClient.EXPECT().
		GetLoadBalancer("lbid").
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
		"health_check": flattenHealthCheck(models.HealthCheck{
			Target:             "HTTP:80/admin/healthcheck",
			Interval:           25,
			Timeout:            10,
			HealthyThreshold:   4,
			UnhealthyThreshold: 3,
		}),
	})

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerRead(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetLoadBalancer("lbid").
		Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{})
	d.SetId("lbid")

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Read(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerUpdate_ports(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	gomock.InOrder(
		mockClient.EXPECT().
			CreateLoadBalancer("test-lb", "test-env", models.HealthCheck{"TCP:80", 30, 5, 2, 2}, []models.Port{}, true).
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			GetLoadBalancer("lbid").
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			UpdateLoadBalancerPorts("lbid", []models.Port{{"", 80, 80, "http"}}).
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			GetLoadBalancer("lbid").
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),
	)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d1 := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
	})

	d2 := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
		"port": flattenPorts([]models.Port{
			{
				ContainerPort: 80,
				HostPort:      80,
				Protocol:      "http",
			},
		}),
	})

	d2.SetId("lbid")

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Create(d1, client); err != nil {
		t.Fatal(err)
	}

	if err := loadBalancerResource.Update(d2, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerUpdate_healthCheck(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	gomock.InOrder(
		mockClient.EXPECT().
			CreateLoadBalancer("test-lb", "test-env", models.HealthCheck{"TCP:80", 30, 5, 2, 2}, []models.Port{}, true).
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			GetLoadBalancer("lbid").
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			UpdateLoadBalancerHealthCheck("lbid", models.HealthCheck{"HTTP:80/admin/healthcheck", 25, 10, 4, 3}).
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),

		mockClient.EXPECT().
			GetLoadBalancer("lbid").
			Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil),
	)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d1 := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
	})

	d2 := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{
		"name":        "test-lb",
		"environment": "test-env",
		"health_check": flattenHealthCheck(models.HealthCheck{
			Target:             "HTTP:80/admin/healthcheck",
			Interval:           25,
			Timeout:            10,
			HealthyThreshold:   4,
			UnhealthyThreshold: 3,
		}),
	})

	d2.SetId("lbid")

	client := &Layer0Client{API: mockClient}
	if err := loadBalancerResource.Create(d1, client); err != nil {
		t.Fatal(err)
	}

	if err := loadBalancerResource.Update(d2, client); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerDelete(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		DeleteLoadBalancer("lbid").
		Return("jid", nil)

	mockClient.EXPECT().
		WaitForJob("jid", gomock.Any()).
		Return(nil)

	loadBalancerResource := provider.ResourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerResource.Schema, map[string]interface{}{})
	d.SetId("lbid")

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := loadBalancerResource.Delete(d, client); err != nil {
		t.Fatal(err)
	}
}
