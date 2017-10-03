package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) List() ([]models.TaskSummary, error) {
	clusterNames, err := listClusterNames(t.AWS.ECS, t.Config.Instance())
	if err != nil {
		return nil, err
	}

	taskARNs := []string{}
	for _, clusterName := range clusterNames {
		startedBy := t.Config.Instance()
		clusterTaskARNsStopped, err := listClusterTaskARNs(t.AWS.ECS, clusterName, startedBy, ecs.DesiredStatusStopped)
		if err != nil {
			return nil, err
		}

		clusterTaskARNsRunning, err := listClusterTaskARNs(t.AWS.ECS, clusterName, startedBy, ecs.DesiredStatusRunning)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, clusterTaskARNsStopped...)
		taskARNs = append(taskARNs, clusterTaskARNsRunning...)
	}

	summaries, err := t.populateSummariesFromTaskARNs(taskARNs)
	if err != nil {
		return nil, err
	}

	return summaries, nil
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
