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
)

func TestServiceLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
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
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	describeServicesInput := &ecs.DescribeServicesInput{}
	describeServicesInput.SetCluster("l0-test-env_id")
	describeServicesInput.SetServices([]*string{aws.String("l0-test-svc_id")})

	deployment1 := &ecs.Deployment{}
	deployment1.SetId("dpl_id:1")
	deployment1.SetStatus(ecs.DesiredStatusRunning)
	deployment1.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")

	deployment2 := &ecs.Deployment{}
	deployment2.SetId("dpl_id:2")
	deployment2.SetStatus(ecs.DesiredStatusStopped)
	deployment2.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:2")

	deployments := []*ecs.Deployment{
		deployment1,
		deployment2,
	}

	service := &ecs.Service{}
	service.SetDeployments(deployments)
	services := []*ecs.Service{
		service,
	}

	describeServicesOutput := &ecs.DescribeServicesOutput{}
	describeServicesOutput.SetServices(services)

	mockAWS.ECS.EXPECT().
		DescribeServices(describeServicesInput).
		Return(describeServicesOutput, nil)

	generateListTasksPagesFN := func(input *ecs.ListTasksInput) func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
		var taskARN *string

		if *input.DesiredStatus == ecs.DesiredStatusRunning && *input.StartedBy == "dpl_id:1" {
			taskARN = aws.String("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")
		}

		if *input.DesiredStatus == ecs.DesiredStatusStopped && *input.StartedBy == "dpl_id:2" {
			aws.String("arn:aws:ecs:region:012345678910:task-definition/dpl_id:2")
		}

		listTasksPagesFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
			taskARNs := []*string{
				taskARN,
			}

			output := &ecs.ListTasksOutput{}
			output.SetTaskArns(taskARNs)
			fn(output, true)

			return nil
		}

		return listTasksPagesFN
	}

	deploymentIDs := []string{
		"dpl_id:1",
		"dpl_id:2",
	}

	statuses := []string{
		ecs.DesiredStatusRunning,
		ecs.DesiredStatusStopped,
	}

	for _, deploymentID := range deploymentIDs {
		for _, status := range statuses {
			input := &ecs.ListTasksInput{}
			input.SetCluster("l0-test-env_id")
			input.SetDesiredStatus(status)
			input.SetStartedBy(deploymentID)

			mockAWS.ECS.EXPECT().
				ListTasksPages(input, gomock.Any()).
				Do(generateListTasksPagesFN(input)).
				Return(nil)
		}
	}

	mockConfig.EXPECT().
		LogGroupName().
		Return("l0-test")

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	_, err := target.Logs("svc_id", 0, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
}
