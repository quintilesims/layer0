package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestService_createTagsWithLoadBalancer(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	if err := service.createTags("svc_id", "svc_name", "env_id", "lb_id"); err != nil {
		t.Fatal(err)
	}

	expected := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "load_balancer_id",
			Value:      "lb_id",
		},
	}

	assert.Equal(t, expected, tagStore.Tags())
}

func TestService_createTagsWithoutLoadBalancer(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	if err := service.createTags("svc_id", "svc_name", "env_id", ""); err != nil {
		t.Fatal(err)
	}

	expected := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
	}

	assert.Equal(t, expected, tagStore.Tags())
}
