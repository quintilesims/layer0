package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of Load Balancers from ELB and returns a list of Load
// Balancer Summaries. A Load Balancer Summary consists of the Load Balancer ID,
// Load Balancer name, Environment ID and Environment name.
func (e *LoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	loadBalancerNames, err := e.listLoadBalancerNames()
	if err != nil {
		return nil, err
	}

	summaries := make([]models.LoadBalancerSummary, len(loadBalancerNames))
	for i, loadBalancerName := range loadBalancerNames {
		loadBalancerID := delLayer0Prefix(e.Config.Instance(), loadBalancerName)
		summary := models.LoadBalancerSummary{
			LoadBalancerID: loadBalancerID,
		}

		summaries[i] = summary
	}

	if err := e.populateSummariesTags(summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (e *LoadBalancerProvider) listLoadBalancerNames() ([]string, error) {
	loadBalancerNames := []string{}
	fn := func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, description := range output.LoadBalancerDescriptions {
			loadBalancerName := aws.StringValue(description.LoadBalancerName)

			if hasLayer0Prefix(e.Config.Instance(), loadBalancerName) {
				loadBalancerNames = append(loadBalancerNames, loadBalancerName)
			}
		}

		return !lastPage
	}

	if err := e.AWS.ELB.DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{}, fn); err != nil {
		return nil, err
	}

	return loadBalancerNames, nil
}

func (e *LoadBalancerProvider) populateSummariesTags(summaries []models.LoadBalancerSummary) error {
	environmentTags, err := e.TagStore.SelectByType("environment")
	if err != nil {
		return err
	}

	loadBalancerTags, err := e.TagStore.SelectByType("load_balancer")
	if err != nil {
		return err
	}

	for i, summary := range summaries {
		if tag, ok := loadBalancerTags.WithID(summary.LoadBalancerID).WithKey("name").First(); ok {
			summaries[i].LoadBalancerName = tag.Value
		}

		if tag, ok := loadBalancerTags.WithID(summary.LoadBalancerID).WithKey("environment_id").First(); ok {
			summaries[i].EnvironmentID = tag.Value

			if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
				summaries[i].EnvironmentName = t.Value
			}
		}
	}

	return nil
}
