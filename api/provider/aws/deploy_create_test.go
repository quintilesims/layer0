package aws

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/config"
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

	if err := deploy.createTags(id, name, version, arn); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
		{
			EntityID:   id,
			EntityType: "deploy",
			Key:        "name",
			Value:      name,
		},
		{
			EntityID:   id,
			EntityType: "deploy",
			Key:        "version",
			Value:      version,
		},
		{
			EntityID:   id,
			EntityType: "deploy",
			Key:        "arn",
			Value:      arn,
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestDeploy_renderTaskDefinition(t *testing.T) {
	c := config.NewTestContext(t, nil, map[string]interface{}{
		config.FlagAWSLogGroup.GetName(): "test_group",
		config.FlagAWSRegion.GetName():   "test_region",
	})

	deploy := NewDeployProvider(nil, nil, c)

	container := &ecs.ContainerDefinition{}
	container.SetName("test_name")
	logContext := &ecs.LogConfiguration{
		LogDriver: aws.String("awslogs"),
		Options: map[string]*string{
			"awslogs-group":         aws.String("test_group"),
			"awslogs-region":        aws.String("test_region"),
			"awslogs-stream-prefix": aws.String("test_prefix"),
		},
	}
	container.SetLogConfiguration(logContext)

	containers := []*ecs.ContainerDefinition{}
	containers = append(containers, container)

	taskDef := &ecs.TaskDefinition{}
	taskDef.SetContainerDefinitions(containers)
	taskDef.SetFamily("test_family")
	taskDef.SetNetworkMode("bridge")

	bytes, err := json.Marshal(taskDef)
	if err != nil {
		t.Fatal("Failed to extract deploy file")
	}

	renderedTaskDef, err := deploy.renderTaskDefinition(bytes, "test_family")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, taskDef, renderedTaskDef)
}

func TestDeploy_renderTaskDefinition_Errors(t *testing.T) {
	deploy := NewDeployProvider(nil, nil, nil)
	model := &ecs.TaskDefinition{}

	container := &ecs.ContainerDefinition{}
	container.SetName("container_name")

	testCases := map[string]*ecs.TaskDefinition{
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
