package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestServiceCreate_defaults(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateService("test-svc", "test-env", "test-dep", "").
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		GetService("sid").
		Return(&models.Service{}, nil)

	mockClient.EXPECT().
		WaitForDeployment("sid", gomock.Any()).
		Return(&models.Service{ServiceID: "sid"}, nil)

	serviceResource := provider.ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":        "test-svc",
		"environment": "test-env",
		"deploy":      "test-dep",
	})

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := serviceResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestServiceCreate_specifyOptional(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateService("test-svc", "test-env", "test-dep", "test-lb").
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		ScaleService("sid", 2).
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		GetService("sid").
		Return(&models.Service{}, nil)

	mockClient.EXPECT().
		WaitForDeployment("sid", gomock.Any()).
		Return(&models.Service{ServiceID: "sid"}, nil)

	serviceResource := provider.ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":          "test-svc",
		"environment":   "test-env",
		"deploy":        "test-dep",
		"load_balancer": "test-lb",
		"scale":         2,
	})

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := serviceResource.Create(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestServiceRead(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetService("sid").
		Return(&models.Service{}, nil)

	serviceResource := provider.ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{})
	d.SetId("sid")

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := serviceResource.Read(d, client); err != nil {
		t.Fatal(err)
	}
}

func TestServiceUpdate(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateService("test-svc", "test-env", "test-dep", "").
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		GetService("sid").
		Return(&models.Service{}, nil)

	mockClient.EXPECT().
		UpdateService("sid", "test-dep2").
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		ScaleService("sid", 2).
		Return(&models.Service{ServiceID: "sid"}, nil)

	mockClient.EXPECT().
		WaitForDeployment("sid", gomock.Any()).
		Return(&models.Service{ServiceID: "sid"}, nil).
		Times(2)

	mockClient.EXPECT().
		GetService("sid").
		Return(&models.Service{}, nil)

	serviceResource := provider.ResourcesMap["layer0_service"]
	d1 := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":        "test-svc",
		"environment": "test-env",
		"deploy":      "test-dep",
	})

	d2 := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":        "test-svc",
		"environment": "test-env",
		"deploy":      "test-dep2",
		"scale":       2,
	})

	d2.SetId("sid")

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := serviceResource.Create(d1, client); err != nil {
		t.Fatal(err)
	}

	if err := serviceResource.Update(d2, client); err != nil {
		t.Fatal(err)
	}
}

func TestServiceDelete(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		DeleteService("sid").
		Return("jid", nil)

	mockClient.EXPECT().
		WaitForJob("jid", gomock.Any()).
		Return(nil)

	serviceResource := provider.ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{})
	d.SetId("sid")

	client := &Layer0Client{API: mockClient, StopContext: context.Background()}
	if err := serviceResource.Delete(d, client); err != nil {
		t.Fatal(err)
	}
}
