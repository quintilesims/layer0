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

func TestDataSourceLayer0DeployRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "deploy")
	query.Set(models.TagQueryParamName, "dpl_name")
	query.Set(models.TagQueryParamVersion, "latest")

	mockClient.EXPECT().
		ListTags(query).
		Return([]models.Tag{{EntityID: "dpl_id"}}, nil)

	deploy := &models.Deploy{
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		Version:    "2",
	}

	mockClient.EXPECT().
		ReadDeploy("dpl_id").
		Return(deploy, nil)

	deployDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployDataSource.Schema, map[string]interface{}{
		"name":    "dpl_name",
		"version": "latest",
	})

	if err := dataSourceLayer0DeployRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", d.Id())
	assert.Equal(t, "dpl_name", d.Get("name"))
	assert.Equal(t, "2", d.Get("version"))
}
