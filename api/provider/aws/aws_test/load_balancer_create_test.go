package aws_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/aws/mock_aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "lb_id1",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name1",
		},
		{
			EntityID:   "lb_id1",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id1",
		},
		{
			EntityID:   "lb_id2",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "env_name2",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "os",
			Value:      "os2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name1",
		EnvironmentID:    "env_id1",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck: models.HealthCheck{
			Target:             "80",
			Interval:           60,
			Timeout:            60,
			HealthyThreshold:   3,
			UnhealthyThreshold: 3,
		},
	}

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"jid1", "jid2"}
	assert.Equal(t, expected, result)
}
