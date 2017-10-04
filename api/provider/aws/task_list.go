package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of all Task ARNs from all of the user's clusters (Environments)
// from ECS and returns a list of Task summaries. A Task summary consists of the Task ID,
// Task name, Environment ID, and Environment name.
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

	taskARNMatches := map[string]bool{}
	for _, taskARN := range taskARNs {
		taskARNMatches[taskARN] = true
	}

	taskModels := make([]models.TaskSummary, 0, len(taskARNs))
	for _, tag := range taskTags.WithKey("arn") {
		if taskARNMatches[tag.Value] {
			model := models.TaskSummary{
				TaskID: tag.EntityID,
			}

			if tag, ok := taskTags.WithID(model.TaskID).WithKey("name").First(); ok {
				model.TaskName = tag.Value
			}

			if tag, ok := taskTags.WithID(model.TaskID).WithKey("environment_id").First(); ok {
				model.EnvironmentID = tag.Value

				if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
					model.EnvironmentName = t.Value
				}
			}

			taskModels = append(taskModels, model)
		}
	}

	return taskModels, nil
}
