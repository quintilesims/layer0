package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) Read(taskID string) (*models.Task, error) {
	environmentID, err := lookupEntityEnvironmentID(t.TagStore, "task", taskID)
	if err != nil {
		return nil, err
	}

	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	task, err := t.readTask(clusterName, taskARN)
	if err != nil {
		return nil, err
	}

	taskFamily, _ := taskFamilyRevisionFromARN(aws.StringValue(task.TaskDefinitionArn))
	deployID := delLayer0Prefix(t.Config.Instance(), taskFamily)

	model, err := t.newModel(taskID, environmentID, deployID)
	if err != nil {
		return nil, err
	}

	model.Containers = make([]models.Container, len(task.Containers))
	for i, c := range task.Containers {
		model.Containers[i] = models.Container{
			ContainerName: aws.StringValue(c.Name),
			Status:        aws.StringValue(c.LastStatus),
			ExitCode:      int(aws.Int64Value(c.ExitCode)),
			Meta:          aws.StringValue(c.Reason),
		}
	}

	model.Status = aws.StringValue(task.LastStatus)
	return model, nil
}

func (t *TaskProvider) readTask(clusterName, taskARN string) (*ecs.Task, error) {
	input := &ecs.DescribeTasksInput{}
	input.SetCluster(clusterName)
	input.SetTasks([]*string{aws.String(taskARN)})

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := t.AWS.ECS.DescribeTasks(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "task was not found") {
			return nil, errors.Newf(errors.TaskDoesNotExist, "The specified task does not exist")
		}

		return nil, err
	}

	if len(output.Failures) > 0 {
		reason := aws.StringValue(output.Failures[0].Reason)
		if strings.Contains(reason, "MISSING") {
			return nil, errors.Newf(errors.TaskDoesNotExist, "The specified task does not exist")
		}

		return nil, fmt.Errorf("Failed to describe task: %s", reason)
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

func (t *TaskProvider) newModel(taskID, environmentID, deployID string) (*models.Task, error) {
	model := &models.Task{
		TaskID:        taskID,
		EnvironmentID: environmentID,
		DeployID:      deployID,
	}

	taskTags, err := t.TagStore.SelectByTypeAndID("task", taskID)
	if err != nil {
		return nil, err
	}

	if tag, ok := taskTags.WithKey("name").First(); ok {
		model.TaskName = tag.Value
	}

	environmentTags, err := t.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return nil, err
	}

	if tag, ok := environmentTags.WithKey("name").First(); ok {
		model.EnvironmentName = tag.Value
	}

	deployTags, err := t.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return nil, err
	}

	if tag, ok := deployTags.WithKey("name").First(); ok {
		model.DeployName = tag.Value
	}

	if tag, ok := deployTags.WithKey("version").First(); ok {
		model.DeployVersion = tag.Value
	}

	return model, nil
}

func taskFamilyRevisionFromARN(taskARN string) (string, string) {
	// task definition arn format: 'arn:aws:ecs:region:account:task-definition/family:version
	familyRevision := strings.Split(taskARN, "/")[1]
	split := strings.SplitN(familyRevision, ":", 2)
	return split[0], split[1]
}
