package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestServiceDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	serviceName := "service-name"
	environmentID := "l0-env-id"

	params := map[string]string{
		"type":           "service",
		"environment_id": environmentID,
		"fuzz":           serviceName,
	}

	mockClient.EXPECT().
		SelectByQuery(params)

	serviceResource := provider.DataSourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{})
	d.Set("name", serviceName)
	d.Set("environment_id", environmentID)

	if err := serviceResource.Read(d, mockClient); err.Error() != fmt.Errorf("No entities found").Error() {
		t.Fatal(err)
	}
}
