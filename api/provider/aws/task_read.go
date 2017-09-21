package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) Read(taskID string) (*models.Task, error) {
	environmentID, err := lookupEntityEnvironmentID(t.TagStore, "task", taskID)
	if err != nil {
		return nil, err
	}
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)

	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	clusterName := fqEnvironmentID
	task, err := t.readTask(clusterName, taskARN)
	if err != nil {
		return nil, err
	}

	taskFamily, _ := taskFamilyRevisionFromARN(aws.StringValue(task.TaskDefinitionArn))
	deployID := delLayer0Prefix(t.Config.Instance(), taskFamily)

	containers := make([]models.Container, len(task.Containers))
	for i, c := range task.Containers {
		containers[i] = models.Container{
			ContainerName: aws.StringValue(c.Name),
			Status:        aws.StringValue(c.LastStatus),
			ExitCode:      int(aws.Int64Value(c.ExitCode)),
			Meta:          aws.StringValue(c.Reason),
		}
	}

	model := &models.Task{
		TaskID:        taskID,
		EnvironmentID: environmentID,
		DeployID:      deployID,
		Status:        aws.StringValue(task.LastStatus),
		Containers:    containers,
	}

	if err := t.populateModelTags(taskID, environmentID, deployID, model); err != nil {
		return nil, err
	}

	return model, nil
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

// todo: we should standardize this pattern for tag lookups
func (t *TaskProvider) populateModelTags(taskID, environmentID, deployID string, model *models.Task) error {
	taskTags, err := t.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return err
	}

	if tag, ok := taskTags.WithKey("name").First(); ok {
		model.TaskName = tag.Value
	}

	environmentTags, err := t.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return err
	}

	if tag, ok := environmentTags.WithKey("name").First(); ok {
		model.EnvironmentName = tag.Value
	}

	deployTags, err := t.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return err
	}

	if tag, ok := deployTags.WithKey("name").First(); ok {
		model.DeployName = tag.Value
	}

	if tag, ok := deployTags.WithKey("version").First(); ok {
		model.DeployVersion = tag.Value
	}

	return nil
}
