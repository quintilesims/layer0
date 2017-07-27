package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestLoadBalancerDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	loadBalancerID := "loadbalancer-id"
	loadbalancerName := "loadbalancer-name"
	environmentID := "l0-env-id"

	params := map[string]string{
		"type":           "load_balancer",
		"environment_id": environmentID,
		"fuzz":           loadbalancerName,
	}

	mockClient.EXPECT().
		SelectByQuery(params).
		Return([]*models.EntityWithTags{
			&models.EntityWithTags{
				EntityID:   loadBalancerID,
				EntityType: "load_balancer",
			},
		}, nil)

	mockClient.EXPECT().
		GetLoadBalancer(loadBalancerID).
		Return(&models.LoadBalancer{}, nil)

	loadbalancerResource := provider.DataSourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadbalancerResource.Schema, map[string]interface{}{
		"name":           loadbalancerName,
		"environment_id": environmentID,
	})

	client := &Layer0Client{API: mockClient}
	if err := loadbalancerResource.Read(d, client); err != nil {
		t.Fatal(err)
	}
}
