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

func TestDataSourceLayer0LoadBalancerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "load_balancer")
	query.Set(models.TagQueryParamName, "lb_name")
	query.Set(models.TagQueryParamEnvironmentID, "env_id")

	mockClient.EXPECT().
		ListTags(query).
		Return([]models.Tag{{EntityID: "lb_id"}}, nil)

	loadBalancer := &models.LoadBalancer{
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		IsPublic:         false,
		URL:              "some_url",
	}

	mockClient.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(loadBalancer, nil)

	loadBalancerDataSource := Provider().(*schema.Provider).DataSourcesMap["layer0_load_balancer"]
	d := schema.TestResourceDataRaw(t, loadBalancerDataSource.Schema, map[string]interface{}{
		"name":           "lb_name",
		"environment_id": "env_id",
	})

	if err := dataSourceLayer0LoadBalancerRead(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "lb_id", d.Id())
	assert.Equal(t, "lb_name", d.Get("name"))
	assert.Equal(t, "env_id", d.Get("environment_id"))
	assert.Equal(t, "env_name", d.Get("environment_name"))
	assert.Equal(t, true, d.Get("private"))
	assert.Equal(t, "some_url", d.Get("url"))
}
