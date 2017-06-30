package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestServiceDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	serviceID := "service-id"
	serviceName := "service-name"
	environmentID := "l0-env-id"

	params := map[string]string{
		"type":           "service",
		"environment_id": environmentID,
		"fuzz":           serviceName,
	}

	mockClient.EXPECT().
		SelectByQuery(params).
		Return([]*models.EntityWithTags{
			&models.EntityWithTags{
				EntityID:   serviceID,
				EntityType: "service",
			},
		}, nil)

	mockClient.EXPECT().
		GetService(serviceID).
		Return(&models.Service{}, nil)

	serviceResource := provider.DataSourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":           serviceName,
		"environment_id": environmentID,
	})

	if err := serviceResource.Read(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
