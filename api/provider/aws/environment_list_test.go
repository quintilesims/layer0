package aws

import (
	"testing"

	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_listTags(t *testing.T) {
	tagStore := tag_store.NewMemoryTagStore()
	environment := NewEnvironmentProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "eid1",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename1",
		},
		{
			EntityID:   "eid1",
			EntityType: "environment",
			Key:        "os",
			Value:      "eos1",
		},
		{
			EntityID:   "eid2",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename2",
		},
		{
			EntityID:   "eid2",
			EntityType: "environment",
			Key:        "os",
			Value:      "eos2",
		},
		{
			EntityID:   "someid",
			EntityType: "environment",
			Key:        "name",
			Value:      "badname",
		},
		{
			EntityID:   "eid1",
			EntityType: "service",
			Key:        "name",
			Value:      "servicename",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	results := []models.EnvironmentSummary{
		{EnvironmentID: "eid1"},
		{EnvironmentID: "eid2"},
	}

	if err := environment.listTags(results); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, 2, results)
	assert.Equal(t, "ename1", results[0].EnvironmentName)
	assert.Equal(t, "eos1", results[0].OperatingSystem)
	assert.Equal(t, "ename2", results[1].EnvironmentName)
	assert.Equal(t, "eos2", results[1].OperatingSystem)
}
