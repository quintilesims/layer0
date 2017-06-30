package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestEnvironmentDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	environmentId := "environment-id"
	environmentName := "environment-name"

	params := map[string]string{
		"type": "environment",
		"fuzz": environmentName,
	}

	mockClient.EXPECT().
		SelectByQuery(params).
		Return([]*models.EntityWithTags{
			&models.EntityWithTags{
				EntityID:   environmentId,
				EntityType: "environment",
			},
		}, nil)

	mockClient.EXPECT().
		GetEnvironment(environmentId).
		Return(&models.Environment{}, nil)

	environmentResource := provider.DataSourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"name": environmentName,
	})

	if err := environmentResource.Read(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
