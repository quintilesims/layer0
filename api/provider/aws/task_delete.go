package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Delete stops an ECS Task using the specified Task ID. The user's active cluster (Environment) is used as a filter
// when the request to stop the Task is made.
func (t *TaskProvider) Delete(taskID string) error {
	environmentID, err := t.lookupTaskEnvironmentID(taskID)
	if err != nil {
		return err
	}
	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)

	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		return err
	}

	clusterName := fqEnvironmentID
	if err := t.stopTask(clusterName, taskARN); err != nil {
		return err
	}

	if err := deleteEntityTags(t.TagStore, "task", taskID); err != nil {
		return err
	}

	return nil
}

func (t *TaskProvider) stopTask(clusterName, taskARN string) error {
	input := &ecs.StopTaskInput{}
	input.SetCluster(clusterName)
	input.SetTask(taskARN)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := t.AWS.ECS.StopTask(input); err != nil {
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "task was not found") {
			return nil
		}
	}

	return nil
}
