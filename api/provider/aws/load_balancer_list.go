package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of Load Balancers from ELB and returns a list of Load
// Balancer Summaries. A Load Balancer Summary consists of the Load Balancer ID,
// Load Balancer name, Environment ID, and Environment name.
func (l *LoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	loadBalancerNames, err := l.listLoadBalancerNames()
	if err != nil {
		return nil, err
	}

	loadBalancerIDs := make([]string, len(loadBalancerNames))
	for i, loadBalancerName := range loadBalancerNames {
		loadBalancerID := delLayer0Prefix(l.Config.Instance(), loadBalancerName)
		loadBalancerIDs[i] = loadBalancerID
	}

	return l.makeLoadBalancerSummaryModels(loadBalancerIDs)
}

func (l *LoadBalancerProvider) listLoadBalancerNames() ([]string, error) {
	loadBalancerNames := []string{}

	// list classic load balancer
	fnELB := func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, description := range output.LoadBalancerDescriptions {
			loadBalancerName := aws.StringValue(description.LoadBalancerName)

			if hasLayer0Prefix(l.Config.Instance(), loadBalancerName) {
				loadBalancerNames = append(loadBalancerNames, loadBalancerName)
			}
		}

		return !lastPage
	}

	if err := l.AWS.ELB.DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{}, fnELB); err != nil {
		return nil, err
	}

	// list application load balancers
	fnALB := func(output *alb.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, description := range output.LoadBalancers {
			loadBalancerName := aws.StringValue(description.LoadBalancerName)

			if hasLayer0Prefix(l.Config.Instance(), loadBalancerName) {
				loadBalancerNames = append(loadBalancerNames, loadBalancerName)
			}
		}

		return !lastPage
	}

	if err := l.AWS.ALB.DescribeLoadBalancersPages(&alb.DescribeLoadBalancersInput{}, fnALB); err != nil {
		return nil, err
	}

	return loadBalancerNames, nil
}

func (l *LoadBalancerProvider) makeLoadBalancerSummaryModels(loadBalancerIDs []string) ([]models.LoadBalancerSummary, error) {
	environmentTags, err := l.TagStore.SelectByType("environment")
	if err != nil {
		return nil, err
	}

	loadBalancerTags, err := l.TagStore.SelectByType("load_balancer")
	if err != nil {
		return nil, err
	}

	models := make([]models.LoadBalancerSummary, len(loadBalancerIDs))
	for i, loadBalancerID := range loadBalancerIDs {
		models[i].LoadBalancerID = loadBalancerID

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("name").First(); ok {
			models[i].LoadBalancerName = tag.Value
		}

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("type").First(); ok {
			models[i].LoadBalancerType = tag.Value
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
