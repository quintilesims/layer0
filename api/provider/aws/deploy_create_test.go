package aws

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestDeploy_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	deploy := NewDeployProvider(nil, tagStore, nil)

	model := &models.Deploy{
		DeployID:   "deploy_id",
		Version:    "deploy_version",
		DeployName: "deploy_name",
	}

	if err := deploy.createTags(model); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "deploy_name",
		},
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "deploy_version",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestDeploy_renderTaskDefinition_errorOnInvalidRequest(t *testing.T) {
	deploy := NewDeployProvider(nil, nil, nil)
	model := &ecs.TaskDefinition{
		Family: aws.String("familyName"),
	}

	bytes, err := json.Marshal(model)
	if err != nil {
		t.Fatal("Failed to extract deploy file")
	}

	if _, err := deploy.renderTaskDefinition(bytes, "familyName"); err == nil {
		t.Fatal("Expected error was nil")
	}
}
