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

	// list classic load balancers
	fnELB := func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, lbd := range output.LoadBalancerDescriptions {
			loadBalancerName := aws.StringValue(lbd.LoadBalancerName)

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
		for _, lb := range output.LoadBalancers {
			loadBalancerName := aws.StringValue(lb.LoadBalancerName)

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

	summaries := make([]models.LoadBalancerSummary, len(loadBalancerIDs))
	for i, loadBalancerID := range loadBalancerIDs {
		summaries[i].LoadBalancerID = loadBalancerID

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("name").First(); ok {
			summaries[i].LoadBalancerName = tag.Value
		}

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("type").First(); ok {
			summaries[i].LoadBalancerType = models.LoadBalancerType(tag.Value)
		}

		if tag, ok := loadBalancerTags.WithID(loadBalancerID).WithKey("environment_id").First(); ok {
			summaries[i].EnvironmentID = tag.Value

			if t, ok := environmentTags.WithID(tag.Value).WithKey("name").First(); ok {
				summaries[i].EnvironmentName = t.Value
			}
		}
	}

	return summaries, nil
}
