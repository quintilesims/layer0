package aws

import (
	"testing"

	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer_populateSummariesTags(t *testing.T) {
	tagStore := tag_store.NewMemoryTagStore()
	loadBalancer := NewLoadBalancerProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "lid1",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lname1",
		},
		{
			EntityID:   "lid1",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "eid1",
		},
		{
			EntityID:   "eid1",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename1",
		},
		{
			EntityID:   "lid2",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lname2",
		},
		{
			EntityID:   "lid2",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "eid2",
		},
		{
			EntityID:   "eid2",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename2",
		},
		{
			EntityID:   "someid",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "badname",
		},
		{
			EntityID:   "lid1",
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

	results := []models.LoadBalancerSummary{
		{LoadBalancerID: "lid1"},
		{LoadBalancerID: "lid2"},
	}

	if err := loadBalancer.populateSummariesTags(results); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, results, 2)
	assert.Equal(t, "lname1", results[0].LoadBalancerName)
	assert.Equal(t, "eid1", results[0].EnvironmentID)
	assert.Equal(t, "ename1", results[0].EnvironmentName)
	assert.Equal(t, "lname2", results[1].LoadBalancerName)
	assert.Equal(t, "eid2", results[1].EnvironmentID)
	assert.Equal(t, "ename2", results[1].EnvironmentName)
}
