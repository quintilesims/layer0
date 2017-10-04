package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestService_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	if err := service.createTags("svc_id", "svc_name", "env_id"); err != nil {
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

	for _, tag := range expected {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestService_lookupTaskDefinitionARNFromDeployID(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "dpl_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := service.lookupTaskDefinitionARNFromDeployID("dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_arn", result)
}
