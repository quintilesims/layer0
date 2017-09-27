package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of Task ARNs from the user's active cluster (Environment) from ECS
// and returns a list of Task summaries. A Task summary consists of the Task ID,
// Task name, Environment ID, and Environment name.
func (t *TaskProvider) List() ([]models.TaskSummary, error) {
	clusterNames, err := listClusterNames(t.AWS.ECS, t.Config.Instance())
	if err != nil {
		return nil, err
	}

	taskARNs := []string{}
	for _, clusterName := range clusterNames {
		startedBy := t.Config.Instance()
		clusterTaskARNs, err := t.listClusterTaskARNs(clusterName, startedBy)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, clusterTaskARNs...)
	}

	summaries, err := t.populateSummariesFromTaskARNs(taskARNs)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (t *TaskProvider) listClusterTaskARNs(clusterName, startedBy string) ([]string, error) {
	taskARNs := []string{}
	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {
		for _, taskARN := range output.TaskArns {
			taskARNs = append(taskARNs, aws.StringValue(taskARN))
		}

		return !lastPage
	}

	for _, status := range []string{ecs.DesiredStatusRunning, ecs.DesiredStatusStopped} {
		input := &ecs.ListTasksInput{}
		input.SetCluster(clusterName)
		input.SetDesiredStatus(status)
		input.SetStartedBy(startedBy)

		if err := t.AWS.ECS.ListTasksPages(input, fn); err != nil {
			return nil, err
		}
	}

	return taskARNs, nil
}

func (t *TaskProvider) populateSummariesFromTaskARNs(taskARNs []string) ([]models.TaskSummary, error) {
	environmentTags, err := t.TagStore.SelectByType("environment")
	if err != nil {
		return nil, err
	}

	taskTags, err := t.TagStore.SelectByType("task")
	if err != nil {
		return nil, err
	}

	summaries := make([]models.TaskSummary, 0, len(taskARNs))
	for _, tag := range taskTags.WithKey("arn") {

		summary := models.TaskSummary{
			TaskID: tag.EntityID,
		}

		if tag, ok := taskTags.WithID(summary.TaskID).WithKey("name").First(); ok {
			summary.TaskName = tag.Value
		}

		if tag, ok := taskTags.WithID(summary.TaskID).WithKey("arn").First(); ok {
			summary.TaskName = tag.Value
		}

		if tag, ok := taskTags.WithID(summary.TaskID).WithKey("environment_id").First(); ok {
			summary.EnvironmentID = tag.Value

			if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
				summary.EnvironmentName = t.Value
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
