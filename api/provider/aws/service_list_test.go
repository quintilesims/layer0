package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceList_makeServiceSummaryModels(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	service := NewServiceProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "svc_id_0",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name_0",
		},
		{
			EntityID:   "svc_id_0",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id_0",
		},
		{
			EntityID:   "env_id_0",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name_0",
		},
		{
			EntityID:   "svc_id_1",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name_1",
		},
		{
			EntityID:   "svc_id_1",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id_1",
		},
		{
			EntityID:   "env_id_1",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name_1",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	serviceIDs := []string{"svc_id_0", "svc_id_1"}
	results, err := service.makeServiceSummaryModels(serviceIDs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, results, 2)
	assert.Equal(t, "svc_id_0", results[0].ServiceID)
	assert.Equal(t, "svc_name_0", results[0].ServiceName)
	assert.Equal(t, "env_id_0", results[0].EnvironmentID)
	assert.Equal(t, "env_name_0", results[0].EnvironmentName)
	assert.Equal(t, "svc_id_1", results[1].ServiceID)
	assert.Equal(t, "svc_name_1", results[1].ServiceName)
	assert.Equal(t, "env_id_1", results[1].EnvironmentID)
	assert.Equal(t, "env_name_1", results[1].EnvironmentName)
}
