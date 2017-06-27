package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestEnvironmentDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	environmentName := "environment-name"

	params := map[string]string{
		"type": "environment",
		"fuzz": environmentName,
	}

	mockClient.EXPECT().
		SelectByQuery(params)

	environmentResource := provider.DataSourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{})
	d.Set("name", environmentName)

	if err := environmentResource.Read(d, mockClient); err.Error() != fmt.Errorf("No entities found").Error() {
		t.Fatal(err)
	}
}
