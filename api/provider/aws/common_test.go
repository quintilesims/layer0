package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func Test_lookupDeployIDFromTaskDefinitionARN(t *testing.T) {
	tagStore := tag.NewMemoryStore()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "task_definition_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := lookupDeployIDFromTaskDefinitionARN(tagStore, "task_definition_arn")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", result)
}

func Test_lookupTaskDefinitionARNFromDeployID(t *testing.T) {
	tagStore := tag.NewMemoryStore()

	tags := models.Tags{
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "task_definition_arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := lookupTaskDefinitionARNFromDeployID(tagStore, "dpl_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "task_definition_arn", result)
}
