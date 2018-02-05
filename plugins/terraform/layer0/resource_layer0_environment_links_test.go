package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceEnvironmentLinksCreateRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	links := []string{"env_id2", "env_id3"}
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	mockClient.EXPECT().
		UpdateEnvironment("env_id1", req).
		Return(nil)

	environment := &models.Environment{
		EnvironmentID: "env_id1",
		Links:         []string{"env_id2", "env_id3"},
	}

	mockClient.EXPECT().
		ReadEnvironment("env_id1").
		Return(environment, nil)

	environmentLinksResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment_links"]
	d := schema.TestResourceDataRaw(t, environmentLinksResource.Schema, map[string]interface{}{
		"environment_id": "env_id1",
		"links":          []interface{}{"env_id2", "env_id3"},
	})

	if err := resourceLayer0EnvironmentLinksCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id1", d.Id())
	assert.Equal(t, []interface{}{"env_id2", "env_id3"}, d.Get("links"))
}

func TestResourceEnvironmentLinksDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	links := []string{}
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	mockClient.EXPECT().
		UpdateEnvironment("env_id1", req).
		Return(nil)

	environmentLinksResource := Provider().(*schema.Provider).ResourcesMap["layer0_environment_links"]
	d := schema.TestResourceDataRaw(t, environmentLinksResource.Schema, map[string]interface{}{})
	d.SetId("env_id1")

	if err := resourceLayer0EnvironmentLinksDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
