package aws

import (
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) List() ([]models.EnvironmentSummary, error) {
	clusterNames, err := listClusterNames(e.AWS.ECS, e.Config.Instance())
	if err != nil {
		return nil, err
	}

	summaries := make([]models.EnvironmentSummary, len(clusterNames))
	for i, clusterName := range clusterNames {
		environmentID := delLayer0Prefix(e.Config.Instance(), clusterName)
		summary := models.EnvironmentSummary{
			EnvironmentID: environmentID,
		}

		summaries[i] = summary
	}

	if err := e.populateSummariesTags(summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (e *EnvironmentProvider) populateSummariesTags(summaries []models.EnvironmentSummary) error {
	tags, err := e.TagStore.SelectByType("environment")
	if err != nil {
		return err
	}

	for i, summary := range summaries {
		if tag, ok := tags.WithID(summary.EnvironmentID).WithKey("name").First(); ok {
			summaries[i].EnvironmentName = tag.Value
		}

		if tag, ok := tags.WithID(summary.EnvironmentID).WithKey("os").First(); ok {
			summaries[i].OperatingSystem = tag.Value
		}
	}

	return nil
}
