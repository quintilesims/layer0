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
		InstanceSize:     "m3.large",
		UserDataTemplate: []byte("template"),
		MinClusterCount:  1,
		OperatingSystem:  "linux",
		AMIID:            "ami123",
	}

	mockClient.EXPECT().
		CreateEnvironment(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: models.CompletedJobStatus,
		Result: "env_id",
	}

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	environment := &models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		ClusterCount:    1,
		InstanceSize:    "m3.large",
		SecurityGroupID: "sgid",
		OperatingSystem: "linux",
		AMIID:           "ami123",
	}

	mockClient.EXPECT().
		ReadEnvironment("env_id").
		Return(environment, nil)

	environmentResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{
		"name":      "env_name",
		"size":      "m3.large",
		"user_data": "template",
		"min_count": 1,
		"os":        "linux",
		"ami":       "ami123",
	})

	if err := resourceLayer0EnvironmentCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", d.Id())
	assert.Equal(t, "env_name", d.Get("name").(string))
	assert.Equal(t, "m3.large", d.Get("size").(string))
	assert.Equal(t, 1, d.Get("cluster_count").(int))
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
		Return("job_id", nil)

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{Status: models.CompletedJobStatus}, nil)

	environmentResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentResource.Schema, map[string]interface{}{})
	d.SetId("env_id")

	if err := resourceLayer0EnvironmentDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
