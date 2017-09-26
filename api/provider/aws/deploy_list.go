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

	summaries := make([]models.DeploySummary, 0, len(taskDefinitionARNs))
	for _, taskDefinitionARN := range taskDefinitionARNs {
		for _, tag := range deployTags.WithKey("arn") {
			// Populate only if there is a matching ARN in tags db
			if tag.Value == taskDefinitionARN {
				summary := models.DeploySummary{
					DeployID: tag.EntityID,
				}

				if tag, ok := deployTags.WithID(summary.DeployID).WithKey("name").First(); ok {
					summary.DeployName = tag.Value
				}

				if tag, ok := deployTags.WithID(summary.DeployID).WithKey("version").First(); ok {
					summary.Version = tag.Value
				}

				summaries = append(summaries, summary)
			}
		}
	}

	return summaries, nil
}
