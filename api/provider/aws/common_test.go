package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func Test_getLaunchTypeFromEnvironmentID(t *testing.T) {
	tagStore := tag.NewMemoryStore()

	envIDs := []string{"env_id0", "env_id1"}

	tags := models.Tags{
		{
			EntityID:   envIDs[0],
			EntityType: "environment",
			Key:        "type",
			Value:      models.EnvironmentTypeDynamic,
		},
		{
			EntityID:   envIDs[1],
			EntityType: "environment",
			Key:        "type",
			Value:      models.EnvironmentTypeStatic,
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	cases := map[string]string{
		envIDs[0]: ecs.LaunchTypeFargate,
		envIDs[1]: ecs.LaunchTypeEc2,
	}

	for id, expected := range cases {
		t.Run(id, func(t *testing.T) {
			result, err := getLaunchTypeFromEnvironmentID(tagStore, id)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, result)
		})
	}
}

func Test_getLaunchTypeFromEnvironmentID_Errors(t *testing.T) {
	tagStore := tag.NewMemoryStore()

	envIDs := []string{"env_id0", "env_id1", "env_id2"}

	tags := models.Tags{
		{
			EntityID:   envIDs[0],
			EntityType: "environment",
			Key:        "type",
			Value:      "",
		},
		{
			EntityID:   envIDs[1],
			EntityType: "environment",
			Key:        "type",
			Value:      "neither static nor dynamic",
		},
		{
			EntityID:   envIDs[2],
			EntityType: "environment",
			Key:        "",
			Value:      "",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	for _, id := range envIDs {
		t.Run(id, func(t *testing.T) {
			if _, err := getLaunchTypeFromEnvironmentID(tagStore, id); err == nil {
				t.Fatal("Err was nil!")
			}
		})
	}
}

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
