package test_aws

import (
	"testing"

	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "linux",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeSecurityGroupHelper(mockAWS, "l0-test-env_id-env", "env_sg_id")
	describeSecurityGroupHelper(mockAWS, "l0-test-lb_id-lb", "lb_sg_id")

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: "lb_name",
		EnvironmentID:    "env_id",
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

	expected := []string{"jid"}
	assert.Equal(t, expected, result)
}
