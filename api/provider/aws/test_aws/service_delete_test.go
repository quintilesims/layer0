package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := mock_config.NewMockAPIConfig(ctrl)

	c.EXPECT().Instance().Return("test").AnyTimes()

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
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "load_balancer_id",
			Value:      "lb_id",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	taskARNs := []*string{
		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-svc_id1"),
		aws.String("arn:aws:ecs:region:012345678910:task/l0-test-svc_id2"),
	}

	listTasksInput := &ecs.ListTasksInput{}
	listTasksInput.SetCluster("l0-test-env_id")
	listTasksInput.SetDesiredStatus(ecs.DesiredStatusRunning)
	listTasksInput.SetStartedBy("l0-test-svc_id")

	listTasksPagesFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
		output := &ecs.ListTasksOutput{}
		output.SetTaskArns(taskARNs)
		fn(output, true)

		return nil
	}

	mockAWS.ECS.EXPECT().
		ListTasksPages(listTasksInput, gomock.Any()).
		Do(listTasksPagesFN).
		Return(nil)

	describeServicesInput := &ecs.DescribeServicesInput{}
	describeServicesInput.SetCluster("l0-test-env_id")
	describeServicesInput.SetServices([]*string{
		aws.String("l0-test-svc_id"),
	})

	deployment := &ecs.Deployment{}
	deployment.SetId("l0-test-svc_id")
	deployment.SetStatus(ecs.DesiredStatusRunning)

	deployments := []*ecs.Deployment{
		deployment,
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

	for _, taskARN := range taskARNs {
		stopTaskInput := &ecs.StopTaskInput{}
		stopTaskInput.SetCluster("l0-test-env_id")
		stopTaskInput.SetTask(*taskARN)

		mockAWS.ECS.EXPECT().
			StopTask(stopTaskInput).
			Return(&ecs.StopTaskOutput{}, nil)
	}

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.SetDesiredCount(int64(0))
	updateServiceInput.SetCluster("l0-test-env_id")
	updateServiceInput.SetService("l0-test-svc_id")

	updateServiceOutput := &ecs.UpdateServiceOutput{}

	mockAWS.ECS.EXPECT().
		UpdateService(updateServiceInput).
		Return(updateServiceOutput, nil)

	deleteServiceInput := &ecs.DeleteServiceInput{}
	deleteServiceInput.SetCluster("l0-test-env_id")
	deleteServiceInput.SetService("l0-test-svc_id")

	deleteServiceOutput := &ecs.DeleteServiceOutput{}

	mockAWS.ECS.EXPECT().
		DeleteService(deleteServiceInput).
		Return(deleteServiceOutput, nil)

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, c)
	if err := target.Delete("svc_id"); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestServiceDelete_idempotenceViaTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := mock_config.NewMockAPIConfig(ctrl)

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, c)
	if err := target.Delete("svc_id"); err != nil {
		t.Fatal(err)
	}
}

func TestServiceDelete_idempotenceViaAWS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	c := mock_config.NewMockAPIConfig(ctrl)

	c.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_id",
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

	service := &ecs.Service{}
	service.SetServiceName("l0-test-svc_id")

	services := []*ecs.Service{
		service,
	}

	describeServicesOutput := &ecs.DescribeServicesOutput{}
	describeServicesOutput.SetServices(services)

	mockAWS.ECS.EXPECT().
		DescribeServices(gomock.Any()).
		Return(describeServicesOutput, nil)

	mockAWS.ECS.EXPECT().
		UpdateService(gomock.Any()).
		Return(&ecs.UpdateServiceOutput{}, nil)

	deleteServiceInput := &ecs.DeleteServiceInput{}
	deleteServiceInput.SetCluster("l0-test-env_id")
	deleteServiceInput.SetService("l0-test-svc_id")

	mockAWS.ECS.EXPECT().
		DeleteService(deleteServiceInput).
		Return(&ecs.DeleteServiceOutput{}, awserr.New("ServiceNotFoundException", "", nil))

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, c)
	if err := target.Delete("svc_id"); err != nil {
		t.Fatal(err)
	}
}
