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

func TestDataSourceLayer0ServiceRead_stateless(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "service")
	query.Set(models.TagQueryParamName, "svc_name")
	query.Set(models.TagQueryParamEnvironmentID, "env_id")

	mockClient.EXPECT().
		ListTags(query).
		Return([]models.Tag{{EntityID: "svc_id"}}, nil)

	service := &models.Service{
		ServiceID:       "svc_id",
		ServiceName:     "svc_name",
		ServiceType:     models.DeployCompatibilityStateless,
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DesiredCount:    2,
	}

	mockClient.EXPECT().
		ReadService("svc_id").
		Return(service, nil)

	serviceDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceDataSource.Schema, map[string]interface{}{
		"name":           "svc_name",
		"environment_id": "env_id",
	})

	if err := dataSourceLayer0ServiceRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "svc_id", d.Id())
	assert.Equal(t, "svc_name", d.Get("name"))
	assert.Equal(t, "env_id", d.Get("environment_id"))
	assert.Equal(t, "env_name", d.Get("environment_name"))
	assert.Equal(t, 2, d.Get("scale"))
	assert.Equal(t, false, d.Get("stateful"))
}

func TestDataSourceLayer0ServiceRead_stateful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "service")
	query.Set(models.TagQueryParamName, "svc_name")
	query.Set(models.TagQueryParamEnvironmentID, "env_id")

	mockClient.EXPECT().
		ListTags(query).
		Return([]models.Tag{{EntityID: "svc_id"}}, nil)

	service := &models.Service{
		ServiceID:       "svc_id",
		ServiceName:     "svc_name",
		ServiceType:     models.DeployCompatibilityStateful,
		EnvironmentID:   "env_id",
		EnvironmentName: "env_name",
		DesiredCount:    2,
	}

	mockClient.EXPECT().
		ReadService("svc_id").
		Return(service, nil)

	serviceDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceDataSource.Schema, map[string]interface{}{
		"name":           "svc_name",
		"environment_id": "env_id",
	})

	if err := dataSourceLayer0ServiceRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "svc_id", d.Id())
	assert.Equal(t, "svc_name", d.Get("name"))
	assert.Equal(t, "env_id", d.Get("environment_id"))
	assert.Equal(t, "env_name", d.Get("environment_name"))
	assert.Equal(t, 2, d.Get("scale"))
	assert.Equal(t, true, d.Get("stateful"))
}
