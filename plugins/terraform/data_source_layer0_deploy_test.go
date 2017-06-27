package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestDeployDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	deployName := "deploy-name"
	version := "1"

	params := map[string]string{
		"type":    "deploy",
		"version": version,
		"fuzz":    deployName,
	}

	mockClient.EXPECT().
		SelectByQuery(params)

	deployResource := provider.DataSourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{})
	d.Set("name", deployName)
	d.Set("version", version)

	if err := deployResource.Read(d, mockClient); err.Error() != fmt.Errorf("No entities found").Error() {
		t.Fatal(err)
	}
}
