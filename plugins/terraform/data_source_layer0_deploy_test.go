package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestDeployDataResourceSelectByQueryParams(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	deployId := "deploy-id"
	deployName := "deploy-name"
	version := "1"

	params := map[string]string{
		"type":    "deploy",
		"version": version,
		"fuzz":    deployName,
	}

	mockClient.EXPECT().
		SelectByQuery(params).
		Return([]*models.EntityWithTags{
			&models.EntityWithTags{
				EntityID:   deployId,
				EntityType: "deploy",
			},
		}, nil)

	mockClient.EXPECT().
		GetDeploy(deployId).
		Return(&models.Deploy{}, nil)

	deployResource := provider.DataSourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{
		"name":    deployName,
		"version": version,
	})

	client := &Layer0Client{API: mockClient}
	if err := deployResource.Read(d, client); err != nil {
		t.Fatal(err)
	}
}
