package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestService_makeDeploymentModel(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl_name",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "0",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := service.makeDeploymentModel("dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", result.DeployID)
	assert.Equal(t, "dpl_name", result.DeployName)
	assert.Equal(t, "0", result.DeployVersion)
}

func TestService_makeServiceModel(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := service.makeServiceModel("env_id", "lb_id", "svc_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id", result.EnvironmentID)
	assert.Equal(t, "env_name", result.EnvironmentName)
	assert.Equal(t, "lb_id", result.LoadBalancerID)
	assert.Equal(t, "lb_name", result.LoadBalancerName)
	assert.Equal(t, "svc_id", result.ServiceID)
	assert.Equal(t, "svc_name", result.ServiceName)
}
