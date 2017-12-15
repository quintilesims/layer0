package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := config.NewTestContext(t, nil, map[string]interface{}{
		config.FlagInstance.GetName(): "test",
	})

	tags := models.Tags{
		{
			EntityID:   "lb_id1",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name1",
		},
		{
			EntityID:   "lb_id1",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id1",
		},
		{
			EntityID:   "env_id1",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name1",
		},
		{
			EntityID:   "lb_id2",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name2",
		},
		{
			EntityID:   "lb_id2",
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      "env_id2",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeLoadBalancerPagesFN := func(input *elb.DescribeLoadBalancersInput, fn func(output *elb.DescribeLoadBalancersOutput, lastPage bool) bool) error {
		loadBalancerDescriptions := []*elb.LoadBalancerDescription{
			{
				LoadBalancerName: aws.String("l0-test-lb_id1"),
			},
			{
				LoadBalancerName: aws.String("l0-test-lb_id2"),
			},
			{
				LoadBalancerName: aws.String("l0-anotherinstance-lb_id1"),
			},
			{
				LoadBalancerName: aws.String("l0-yetanotherinstance-lb_id1"),
			},
		}

		output := &elb.DescribeLoadBalancersOutput{}
		output.SetLoadBalancerDescriptions(loadBalancerDescriptions)
		fn(output, true)

		return nil
	}

	mockAWS.ELB.EXPECT().
		DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{}, gomock.Any()).
		Do(describeLoadBalancerPagesFN).
		Return(nil)

	target := provider.NewLoadBalancerProvider(mockAWS.Client(), tagStore, c)
	result, err := target.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.LoadBalancerSummary{
		{
			LoadBalancerID:   "lb_id1",
			LoadBalancerName: "lb_name1",
			EnvironmentID:    "env_id1",
			EnvironmentName:  "env_name1",
		},
		{
			LoadBalancerID:   "lb_id2",
			LoadBalancerName: "lb_name2",
			EnvironmentID:    "env_id2",
			EnvironmentName:  "env_name2",
		},
	}

	assert.Equal(t, expected, result)
}
