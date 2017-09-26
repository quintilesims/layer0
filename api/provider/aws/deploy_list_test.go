package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_populateSummariesFromTaskDefinitionARNs(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	deploy := NewDeployProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "d1",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy1",
		},
		{
			EntityID:   "d1",
			EntityType: "deploy",
			Key:        "version",
			Value:      "3",
		},
		{
			EntityID:   "d1",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn1",
		},
		{
			EntityID:   "d2",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy2",
		},
		{
			EntityID:   "d2",
			EntityType: "deploy",
			Key:        "version",
			Value:      "4",
		},
		{
			EntityID:   "d2",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn2",
		},
		{
			EntityID:   "d3",
			EntityType: "environment",
			Key:        "name",
			Value:      "bad_deploy2",
		},
		{
			EntityID:   "d3",
			EntityType: "environment",
			Key:        "version",
			Value:      "5",
		},
		{
			EntityID:   "d3",
			EntityType: "environment",
			Key:        "arn",
			Value:      "arn3",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	taskDefinitionARNs := []string{
		"arn1",
		"arn2",
		"arn3",
	}

	summaries, err := deploy.populateSummariesFromTaskDefinitionARNs(taskDefinitionARNs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, summaries, 2)
	assert.Equal(t, "deploy1", summaries[0].DeployName)
	assert.Equal(t, "3", summaries[0].Version)
	assert.Equal(t, "deploy2", summaries[1].DeployName)
	assert.Equal(t, "4", summaries[1].Version)
}
