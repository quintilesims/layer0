package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of Task Definition ARNs from ECS and returns a list of Deploy summaries.
// A Deploy summary consists of the Deploy ID, Deploy name, and Version.
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
	taskDefinitionFamilies := []string{}
	listTaskDefinitionFamiliesPagesfn := func(output *ecs.ListTaskDefinitionFamiliesOutput, lastPage bool) bool {
		for _, taskDefinitionFamily := range output.Families {
			taskDefinitionFamilies = append(taskDefinitionFamilies, aws.StringValue(taskDefinitionFamily))
		}
		return !lastPage
	}

	familyPrefix := addLayer0Prefix(d.Config.Instance(), "")
	input := &ecs.ListTaskDefinitionFamiliesInput{}
	input.SetFamilyPrefix(familyPrefix)
	// TODO: Revisit how Inactive and Active Task Definitions might want to be returned to the client
	input.SetStatus(ecs.TaskDefinitionFamilyStatusActive)
	if err := d.AWS.ECS.ListTaskDefinitionFamiliesPages(input, listTaskDefinitionFamiliesPagesfn); err != nil {
		return nil, err
	}

	taskDefinitionARNs := []string{}
	listTaskDefinitionPagesfn := func(output *ecs.ListTaskDefinitionsOutput, lastPage bool) bool {
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
