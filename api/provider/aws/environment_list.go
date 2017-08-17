package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) List() ([]models.EnvironmentSummary, error) {
	clusterNames, err := e.listClusterNames()
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

	if err := e.listTags(summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (e *EnvironmentProvider) listClusterNames() ([]string, error) {
	output, err := e.AWS.ECS.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	clusterNames := []string{}
	for _, arn := range output.ClusterArns {
		// cluster arn format: arn:aws:ecs:region:012345678910:cluster/name
		clusterName := strings.Split(aws.StringValue(arn), "/")[1]

		if hasLayer0Prefix(e.Config.Instance(), clusterName) {
			clusterNames = append(clusterNames, clusterName)
		}
	}

	return clusterNames, nil
}

func (e *EnvironmentProvider) listTags(summaries []models.EnvironmentSummary) error {
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
