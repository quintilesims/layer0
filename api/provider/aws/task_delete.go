package aws

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
)

// Delete stops an ECS Task using the specified taskID. The taskID is used to look up the name of
// the ECS Cluster (Environment) the Task resides in. The Cluster name is used when
// the StopTask request is made to AWS.
func (t *TaskProvider) Delete(taskID string) error {
	environmentID, err := lookupEntityEnvironmentID(t.TagStore, "task", taskID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.TaskDoesNotExist {
			log.Printf("[WARN] Task not found\n")
			return nil
		}
		return err
	}

	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.TaskDoesNotExist {
			log.Printf("[WARN] Task not found\n")
			return nil
		}
		return err
	}

	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)
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
