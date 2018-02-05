package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceEnvironmentCreateRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  "env_name",
		InstanceType:     "t2.small",
		EnvironmentType:  "static",
		UserDataTemplate: []byte("template"),
		Scale:            3,
		OperatingSystem:  "linux",
		AMIID:            "ami123",
	}

	mockClient.EXPECT().
		CreateEnvironment(req).
		Return("env_id", nil)

	environment := &models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		EnvironmentType: "static",
		DesiredScale:    3,
		InstanceType:    "t2.small",
		SecurityGroupID: "sgid",
		OperatingSystem: "linux",
		AMIID:           "ami123",
	}

	mockClient.EXPECT().
		ReadEnvironment("env_id").
		Return(environment, nil)

	environmentResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"name":             "env_name",
		"instance_type":    "t2.small",
		"environment_type": "static",
		"user_data":        "template",
		"scale":            3,
		"os":               "linux",
		"ami":              "ami123",
	})

	if err := resourceLayer0EnvironmentCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", d.Id())
	assert.Equal(t, "env_name", d.Get("name").(string))
	assert.Equal(t, "static", d.Get("environment_type").(string))
	assert.Equal(t, "t2.small", d.Get("instance_type").(string))
	assert.Equal(t, 3, d.Get("scale").(int))
	assert.Equal(t, "sgid", d.Get("security_group_id").(string))
	assert.Equal(t, "linux", d.Get("os").(string))
	assert.Equal(t, "ami123", d.Get("ami").(string))
}

func TestResourceEnvironmentDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	mockClient.EXPECT().
		DeleteEnvironment("env_id").
		Return(nil)

	environmentResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{})
	d.SetId("env_id")

	if err := resourceLayer0EnvironmentDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
