package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/models"
)

func (e *LoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	loadBalancerNames, err := e.listLoadBalancerNames()
	if err != nil {
		return nil, err
	}

	loadBalancerIDs := make([]string, len(loadBalancerNames))
	for i, loadBalancerName := range loadBalancerNames {
		loadBalancerID := delLayer0Prefix(e.Config.Instance(), loadBalancerName)
		loadBalancerIDs[i] = loadBalancerID
	}

	return e.newSummaryModels(loadBalancerIDs)
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

func (e *LoadBalancerProvider) newSummaryModels(loadBalancerIDs []string) ([]models.LoadBalancerSummary, error) {
	environmentTags, err := e.TagStore.SelectByType("environment")
	if err != nil {
		return nil, err
	}

	loadBalancerTags, err := e.TagStore.SelectByType("load_balancer")
	if err != nil {
		return nil, err
	}

	models := make([]models.LoadBalancerSummary, len(loadBalancerIDs))
	for i, loadBalancerID := range loadBalancerIDs {
		models[i].LoadBalancerID = loadBalancerID

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("name").First(); ok {
			models[i].LoadBalancerName = tag.Value
		}

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("environment_id").First(); ok {
			models[i].EnvironmentID = tag.Value

			if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
				models[i].EnvironmentName = t.Value
			}
		}
	}

	return models, nil
}
