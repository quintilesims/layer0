package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) List() ([]models.TaskSummary, error) {
	clusterNames, err := listClusterNames(t.AWS.ECS, t.Config.Instance())
	if err != nil {
		return nil, err
	}

	summaries := []models.TaskSummary{}
	for _, clusterName := range clusterNames {
		runningTaskIDs, err := t.listTaskIDs(clusterName, ecs.DesiredStatusRunning)
		if err != nil {
			return nil, err
		}

		stoppedTaskIDs, err := t.listTaskIDs(clusterName, ecs.DesiredStatusStopped)
		if err != nil {
			return nil, err
		}

		for _, taskID := range append(runningTaskIDs, stoppedTaskIDs...) {
			summary := models.TaskSummary{
				TaskID: taskID,
			}
			summaries = append(summaries, summary)
		}

	}

	if err := t.populateTaskSummaries(summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (t *TaskProvider) listTaskIDs(clusterName, status string) ([]string, error) {
	taskIDs := []string{}
	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {
		for _, arn := range output.TaskArns {
			// task arn format: arn:aws:ecs:region:012345678910:task/taskID
			taskID := strings.Split(aws.StringValue(arn), "/")[1]
			taskIDs = append(taskIDs, taskID)
		}
		return !lastPage
	}

	input := &ecs.ListTasksInput{}
	input.SetCluster(clusterName)
	input.SetDesiredStatus(status)

	if err := t.AWS.ECS.ListTasksPages(input, fn); err != nil {
		return nil, err
	}

	return taskIDs, nil
}

func (t *TaskProvider) populateTaskSummaries(summaries []models.TaskSummary) error {
	environmentTags, err := t.TagStore.SelectByType("environment")
	if err != nil {
		return err
	}

	taskTags, err := t.TagStore.SelectByType("task")
	if err != nil {
		return err
	}

	for i, summary := range summaries {
		if tag, ok := taskTags.WithID(summary.TaskID).WithKey("name").First(); ok {
			summaries[i].TaskName = tag.Value
		}

		if tag, ok := taskTags.WithID(summary.TaskID).WithKey("environment_id").First(); ok {
			summaries[i].EnvironmentID = tag.Value

			if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
				summaries[i].EnvironmentName = t.Value
			}
		}
	}

	return nil
}
