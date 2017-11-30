package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceLayer0APIRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	config := &models.APIConfig{
		Instance:       "test",
		VPCID:          "some_vpc",
		Version:        "some_version",
		PublicSubnets:  []string{"pub1", "pub2"},
		PrivateSubnets: []string{"pri1", "pri2"},
	}

	mockClient.EXPECT().
		ReadConfig().
		Return(config, nil)

	apiDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_api"]
	d := schema.TestResourceDataRaw(t, apiDataSource.Schema, map[string]interface{}{})
	if err := dataSourceLayer0APIRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test", d.Id())
	assert.Equal(t, "test", d.Get("instance"))
	assert.Equal(t, "some_vpc", d.Get("vpc_id"))
	assert.Equal(t, "some_version", d.Get("version"))
	assert.Equal(t, []interface{}{"pub1", "pub2"}, d.Get("public_subnets"))
	assert.Equal(t, []interface{}{"pri1", "pri2"}, d.Get("private_subnets"))
}
