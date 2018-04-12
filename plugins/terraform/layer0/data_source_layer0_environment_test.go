package layer0

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceLayer0EnvironmentRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "environment")
	query.Set(models.TagQueryParamName, "env_name")

	mockClient.EXPECT().
		ListTags(query).
		Return([]models.Tag{{EntityID: "env_id"}}, nil)

	environment := &models.Environment{
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DesiredScale:    2,
		InstanceType:    "t2.small",
		SecurityGroupID: "some_sg",
		OperatingSystem: "some_os",
		AMIID:           "some_ami",
	}

	mockClient.EXPECT().
		ReadEnvironment("env_id").
		Return(environment, nil)

	environmentDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_environment"]
	d := schema.TestResourceDataRaw(t, environmentDataSource.Schema, map[string]interface{}{
		"name": "env_name",
	})

	if err := dataSourceLayer0EnvironmentRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", d.Id())
	assert.Equal(t, "env_name", d.Get("name"))
	assert.Equal(t, 2, d.Get("scale"))
	assert.Equal(t, "t2.small", d.Get("instance_type"))
	assert.Equal(t, "some_sg", d.Get("security_group_id"))
	assert.Equal(t, "some_os", d.Get("os"))
	assert.Equal(t, "some_ami", d.Get("ami"))
}
