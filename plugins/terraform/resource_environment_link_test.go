package main

import (
	"errors"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/models"
)

func TestEnvironmentLink_defaults(t *testing.T) {
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

func TestEnvironmentLinkRead_noenv(t *testing.T) {
	ctrl, mockClient, provider := setupUnitTest(t)
	defer ctrl.Finish()

	mockClient.EXPECT().
		GetEnvironment("test-env").
		Return(&models.Environment{}, errors.New("No environment found"))

	environmentResource := provider.ResourcesMap["layer0_environment_link"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"source": "test-env",
		"dest":   "",
	})

	if err := environmentResource.Read(d, mockClient); err != nil {
		t.Fatal(err)
	}
	// not sure if checking the ID is needed here as the other client calls already
	// catch our desired behavior and print to stdout
	// ex:
	// 2017/04/20 14:34:45 [WARN] Error Reading Environment Link (test-env), link does not exist
	// 2017/04/20 14:34:45 [WARN] Error Reading Environment (test-env), environment does not exist

	if d.Id() == "" {
		log.Printf("[WARN] Resource does not exist (%s)", "test-env")
	}
}

func TestEnvironmentLinkDelete_default(t *testing.T) {
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
