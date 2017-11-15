package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	loadBalancer := NewLoadBalancerProvider(nil, tagStore, nil)

	if err := loadBalancer.createTags("lb_id", "lb_name", "env_id"); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
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
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestLoadBalancer_renderLoadBalancerRolePolicy(t *testing.T) {
	template := "{{ .Region }} {{ .AccountID }} {{ .LoadBalancerID }}"

	policy, err := RenderLoadBalancerRolePolicy("region", "account_id", "lbid", template)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "region account_id lbid", policy)
}
