package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Create runs an ECS Task using the specified Create Task Request. The Create Task
// Request contains the Task name, the Environment ID, and the Deploy ID.
// The Deploy ID is used to look up the ECS Task Definition family and version of the
// Task to run.
func (t *TaskProvider) Create(req models.CreateTaskRequest) (string, error) {
	taskID := generateEntityID(req.TaskName)
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), req.EnvironmentID)

	deployName, deployVersion, err := lookupDeployNameAndVersion(t.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	clusterName := fqEnvironmentID
	startedBy := t.Config.Instance()
	taskDefinitionFamily := addLayer0Prefix(t.Config.Instance(), deployName)
	taskDefinitionVersion := deployVersion
	overrides := req.ContainerOverrides

	task, err := t.runTask(clusterName, startedBy, taskDefinitionFamily, taskDefinitionVersion, overrides)
	if err != nil {
		return "", err
	}

	taskARN := aws.StringValue(task.TaskArn)
	if err := t.createTags(taskID, req.TaskName, req.EnvironmentID, taskARN); err != nil {
		return "", err
	}

	return taskID, nil
}

func (t *TaskProvider) runTask(clusterName, startedBy, taskDefinitionFamily, taskDefinitionRevision string, overrides []models.ContainerOverride) (*ecs.Task, error) {
	input := &ecs.RunTaskInput{}
	input.SetCluster(clusterName)
	input.SetStartedBy(startedBy)

	taskFamilyRevision := fmt.Sprintf("%s:%s", taskDefinitionFamily, taskDefinitionRevision)
	input.SetTaskDefinition(taskFamilyRevision)

	newContainerOverride := func(envVars map[string]string) []*ecs.KeyValuePair {
		environment := []*ecs.KeyValuePair{}
		for k, v := range envVars {
			name := k
			value := v
			environment = append(environment, &ecs.KeyValuePair{
				Name:  &name,
				Value: &value,
			})
		}
		return environment
	}

	// Convert to Task Overrides
	var taskOverride *ecs.TaskOverride
	if overrides != nil {
		containerOverrides := []*ecs.ContainerOverride{}
		for _, c := range overrides {
			containerOverride := &ecs.ContainerOverride{}
			containerOverride.SetName(c.ContainerName)

			// Convert Map to Key Value
			overrideKV := newContainerOverride(c.EnvironmentOverrides)
			containerOverride.SetEnvironment(overrideKV)
			containerOverrides = append(containerOverrides, containerOverride)
		}

		taskOverride = &ecs.TaskOverride{
			ContainerOverrides: containerOverrides,
		}

		input.SetOverrides(taskOverride)
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := t.AWS.ECS.RunTask(input)
	if err != nil {
		return nil, err
	}

	if len(output.Failures) > 0 {
		return nil, fmt.Errorf("Failed to create task: %s", aws.StringValue(output.Failures[0].Reason))
	}

	return output.Tasks[0], nil
}

func (t *TaskProvider) createTags(taskID, taskName, environmentID, taskARN string) error {
	tags := []models.Tag{
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "name",
			Value:      taskName,
		},
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "environment_id",
			Value:      environmentID,
		},
		{
			EntityID:   taskID,
			EntityType: "task",
			Key:        "arn",
			Value:      taskARN,
		},
	}

	for _, tag := range tags {
		if err := t.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
