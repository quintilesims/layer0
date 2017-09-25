package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) Read(taskID string) (*models.Task, error) {
	// todo: use tlake's
	environmentID, err := t.lookupTaskEnvironmentID(taskID)
	if err != nil {
		return nil, err
	}
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)

	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	clusterName := fqEnvironmentID
	if _, err := t.readTask(clusterName, taskARN); err != nil {
		return nil, err
	}

	model := &models.Task{
		TaskID:        taskID,
		EnvironmentID: environmentID,
	}

	return model, nil
}

// todo remove
func (t *TaskProvider) lookupTaskEnvironmentID(taskID string) (string, error) {
	tags, err := t.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Failed to find environment_id for task '%s'", taskID)
}

func (t *TaskProvider) readTask(clusterName, taskARN string) (*ecs.Task, error) {
	input := &ecs.DescribeTasksInput{}
	input.SetCluster(clusterName)
	input.SetTasks([]*string{aws.String(taskARN)})

	if err := input.Validate(); err != nil {
		return nil, err
	}

	// todo: catch does ntoe xist error
	output, err := t.AWS.ECS.DescribeTasks(input)
	if err != nil {
		return nil, err
	}

	if len(output.Failures) > 0 {
		// todo: catch does not exist eror rhere?
		return nil, fmt.Errorf("Failed to describe task: %s", aws.StringValue(output.Failures[0].Reason))
	}

	return output.Tasks[0], nil
}

func (t *TaskProvider) lookupTaskARN(taskID string) (string, error) {
	tags, err := t.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.Newf(errors.TaskDoesNotExist, "Task '%s' does not exist", taskID)
	}

	if tag, ok := tags.WithKey("arn").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Failed to find ARN for task '%s'", taskID)
}
