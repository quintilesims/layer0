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

	name := "deploy_name"
	id := "deploy_id"
	version := "deploy_version"
	arn := "deploy_arn"

	if err := deploy.createTags(name, id, version, arn); err != nil {
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
		{
			EntityID:   "deploy_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "deploy_arn",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestDeploy_renderTaskDefinition(t *testing.T) {
	deploy := NewDeployProvider(nil, nil, nil)

	container := &ecs.ContainerDefinition{}
	container.SetName("test_name")
	logConfig := &ecs.LogConfiguration{
		LogDriver: aws.String("awslogs"),
		Options: map[string]*string{
			"awslogs-group":         aws.String("test_group"),
			"awslogs-region":        aws.String("test_region"),
			"awslogs-stream-prefix": aws.String("test_prefix"),
		},
	}
	container.SetLogConfiguration(logConfig)

	containers := []*ecs.ContainerDefinition{}
	containers = append(containers, container)

	taskDef := &ecs.TaskDefinition{}
	taskDef.SetContainerDefinitions(containers)

	bytes, err := json.Marshal(taskDef)
	if err != nil {
		t.Fatal("Failed to extract deploy file")
	}

	if _, err := deploy.renderTaskDefinition(bytes, "familyName"); err != nil {
		t.Fatal(err)
	}
}

func TestDeploy_renderTaskDefinition_Errors(t *testing.T) {
	deploy := NewDeployProvider(nil, nil, nil)
	model := &ecs.TaskDefinition{}

	testCases := map[string]*ecs.TaskDefinition{
		"Custom Family Name": &ecs.TaskDefinition{
			Family: aws.String("customName"),
		},
		"No Container Definitions": &ecs.TaskDefinition{},
	}

	bytes, err := json.Marshal(model)
	if err != nil {
		t.Fatal("Failed to extract deploy file")
	}

	for _, test := range testCases {
		if _, err := deploy.renderTaskDefinition(bytes, aws.StringValue(test.Family)); err == nil {
			t.Fatal("Expected error was nil")
		}
	}
}
