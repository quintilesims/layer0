package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestLoadBalancerDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	loadbalancerName := "loadbalancer-name"
	environmentID := "l0-env-id"

	params := map[string]string{
		"type":           "load_balancer",
		"environment_id": environmentID,
		"fuzz":           loadbalancerName,
	}

	mockClient.EXPECT().
		SelectByQuery(params)

	loadbalancerResource := provider.DataSourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadbalancerResource.Schema, map[string]interface{}{})
	d.Set("name", loadbalancerName)
	d.Set("environment_id", environmentID)

	if err := loadbalancerResource.Read(d, mockClient); err.Error() != fmt.Errorf("No entities found").Error() {
		t.Fatal(err)
	}
}
