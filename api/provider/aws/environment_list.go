package aws

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of ECS Cluster names and returns a list of Environment
// Summaries. An Environment Summary consists of the Environment ID, Environment name,
// and Operating System.
func (e *EnvironmentProvider) List() ([]models.EnvironmentSummary, error) {
	instance := e.Context.String(config.FlagInstance.GetName())
	clusterNames, err := listClusterNames(e.AWS.ECS, instance)
	if err != nil {
		return nil, err
	}

	environmentIDs := make([]string, len(clusterNames))
	for i, clusterName := range clusterNames {
		environmentID := delLayer0Prefix(e.Context, clusterName)
		environmentIDs[i] = environmentID
	}

	return e.makeEnvironmentSummaryModels(environmentIDs)
}

func (e *EnvironmentProvider) makeEnvironmentSummaryModels(environmentIDs []string) ([]models.EnvironmentSummary, error) {
	tags, err := e.TagStore.SelectByType("environment")
	if err != nil {
		return nil, err
	}

	models := make([]models.EnvironmentSummary, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		models[i].EnvironmentID = environmentID

		if tag, ok := tags.WithID(environmentID).WithKey("name").First(); ok {
			models[i].EnvironmentName = tag.Value
		}

		if tag, ok := tags.WithID(environmentID).WithKey("os").First(); ok {
			models[i].OperatingSystem = tag.Value
		}
	}

	return models, nil
}
