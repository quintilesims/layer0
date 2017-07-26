package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestDeployCreate(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateDeploy("test-dep", []byte("sample task definition")).
		Return(&models.Deploy{DeployID: "did"}, nil)

	mockClient.EXPECT().
		GetDeploy("did").
		Return(&models.Deploy{}, nil)

	deployResource := provider.ResourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{
		"name":    "test-dep",
		"content": "sample task definition",
	})

	client := &Layer0Client{API: mockClient}
	if err := deployResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestDeployRead(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetDeploy("did").
		Return(&models.Deploy{}, nil)

	deployResource := provider.ResourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{})
	d.SetId("did")

	client := &Layer0Client{API: mockClient}
	if err := deployResource.Read(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestDeployDelete(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		DeleteDeploy("did").
		Return(nil)

	deployResource := provider.ResourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{})
	d.SetId("did")

	client := &Layer0Client{API: mockClient}
	if err := deployResource.Delete(d, client); err != nil {
		t.Fatal(err)
	}
}
