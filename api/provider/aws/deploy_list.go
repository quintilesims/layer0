package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (d *DeployProvider) List() ([]models.DeploySummary, error) {
	taskDefinitionARNs, err := d.listTaskDefinitionARNs()
	if err != nil {
		return nil, err
	}

	summaries, err := d.populateSummariesFromTaskDefinitionARNs(taskDefinitionARNs)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (d *DeployProvider) listTaskDefinitionARNs() ([]string, error) {
	taskDefinitionARNs := []string{}
	fn := func(output *ecs.ListTaskDefinitionsOutput, lastPage bool) bool {
		for _, taskDefinitionARN := range output.TaskDefinitionArns {
			taskDefinitionARNs = append(taskDefinitionARNs, aws.StringValue(taskDefinitionARN))
		}
		return !lastPage
	}

	input := &ecs.ListTaskDefinitionsInput{}
	if err := d.AWS.ECS.ListTaskDefinitionsPages(input, fn); err != nil {
		return nil, err
	}

	return taskDefinitionARNs, nil
}

func (d *DeployProvider) populateSummariesFromTaskDefinitionARNs(taskDefinitionARNs []string) ([]models.DeploySummary, error) {
	deployTags, err := d.TagStore.SelectByType("deploy")
	if err != nil {
		return nil, err
	}

	taskDefinitionARNMatches := map[string]bool{}
	for _, taskDefinitionARN := range taskDefinitionARNs {
		taskDefinitionARNMatches[taskDefinitionARN] = true
	}

	deploySummaries := make([]models.DeploySummary, 0, len(taskDefinitionARNs))
	for _, tag := range deployTags.WithKey("arn") {
		// Populate only if there is a matching ARN in tags db
		if taskDefinitionARNMatches[tag.Value] {
			deploySummary := models.DeploySummary{
				DeployID: tag.EntityID,
			}

			if tag, ok := deployTags.WithID(deploySummary.DeployID).WithKey("name").First(); ok {
				deploySummary.DeployName = tag.Value
			}

			if tag, ok := deployTags.WithID(deploySummary.DeployID).WithKey("version").First(); ok {
				deploySummary.Version = tag.Value
			}

			deploySummaries = append(deploySummaries, deploySummary)
		}
	}

	return deploySummaries, nil
}
