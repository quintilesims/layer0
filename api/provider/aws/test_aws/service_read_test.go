package test_aws

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name",
		},
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task-definition/dpl_id:1",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl_name",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "1",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "load_balancer_id",
			Value:      "lb_id",
		},
		{
			EntityID:   "lb_id",
			EntityType: "load_balancer",
			Key:        "name",
			Value:      "lb_name",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeServicesInput := &ecs.DescribeServicesInput{}
	describeServicesInput.SetCluster("l0-test-env_id")
	describeServicesInput.SetServices([]*string{aws.String("l0-test-svc_id")})

	deployment := &ecs.Deployment{}
	deployment.SetCreatedAt(time.Time{})
	deployment.SetDesiredCount(int64(1))
	deployment.SetPendingCount(int64(0))
	deployment.SetRunningCount(int64(1))
	deployment.SetStatus(ecs.DesiredStatusRunning)
	deployment.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")
	deployment.SetUpdatedAt(time.Time{})
	deployments := []*ecs.Deployment{
		deployment,
	}

	loadBalancer := &ecs.LoadBalancer{}
	loadBalancer.SetLoadBalancerName("l0-test-lb_id")
	loadBalancers := []*ecs.LoadBalancer{
		loadBalancer,
	}

	service := &ecs.Service{}
	service.SetDeployments(deployments)
	service.SetDesiredCount(int64(1))
	service.SetLoadBalancers(loadBalancers)
	service.SetPendingCount(int64(0))
	service.SetRunningCount(int64(1))
	services := []*ecs.Service{
		service,
	}

	describeServicesOutput := &ecs.DescribeServicesOutput{}
	describeServicesOutput.SetServices(services)

	mockAWS.ECS.EXPECT().
		DescribeServices(describeServicesInput).
		Return(describeServicesOutput, nil)

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.Read("svc_id")
	if err != nil {
		t.Fatal(err)
	}

	expectedDeployments := []models.Deployment{
		{
			Created:       time.Time{},
			DeployID:      "dpl_id",
			DeployName:    "dpl_name",
			DeployVersion: "1",
			DesiredCount:  1,
			PendingCount:  0,
			RunningCount:  1,
			Status:        ecs.DesiredStatusRunning,
			Updated:       time.Time{},
		},
	}

	expected := &models.Service{
		Deployments:      expectedDeployments,
		DesiredCount:     1,
		EnvironmentID:    "env_id",
		EnvironmentName:  "env_name",
		LoadBalancerID:   "lb_id",
		LoadBalancerName: "lb_name",
		PendingCount:     0,
		RunningCount:     1,
		ServiceID:        "svc_id",
		ServiceName:      "svc_name",
	}

	assert.Equal(t, expected, result)
}
