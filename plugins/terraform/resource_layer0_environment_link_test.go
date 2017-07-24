package main

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestEnvironmentLinkCreate(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		CreateLink("test-env", "test-env2").
		Return(nil)

	mockClient.EXPECT().
		GetEnvironment("test-env").
		Return(&models.Environment{}, nil)

	environmentResource := provider.ResourcesMap["layer0_environment_link"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"source": "test-env",
		"dest":   "test-env2",
	})

	if err := environmentResource.Create(d, mockClient); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentLinkRead(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetEnvironment("test-env").
		Return(&models.Environment{}, nil)

	environmentResource := provider.ResourcesMap["layer0_environment_link"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"source": "test-env",
		"dest":   "test-env2",
	})

	if err := environmentResource.Read(d, mockClient); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentLinkRead_noenv(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetEnvironment("test-env").
		Return(nil, errors.New("No environment found"))

	environmentResource := provider.ResourcesMap["layer0_environment_link"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"source": "test-env",
		"dest":   "test-env2",
	})

	if err := environmentResource.Read(d, mockClient); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "" {
		t.Fatal("Id should be set to empty string")
	}
}

func TestEnvironmentLinkDelete(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		DeleteLink("test-env", "test-env2").
		Return(nil)

	environmentResource := provider.ResourcesMap["layer0_environment_link"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"source": "test-env",
		"dest":   "test-env2",
	})

	if err := environmentResource.Delete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
