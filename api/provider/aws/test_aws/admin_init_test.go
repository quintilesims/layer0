package test_aws

import (
	"testing"

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

func TestAdminInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	describeServicesInput := &ecs.DescribeServicesInput{}
	describeServicesInput.SetCluster("l0-test-api")
	describeServicesInput.SetServices([]*string{aws.String("l0-test-api")})

	service := &ecs.Service{}
	taskDefinitionARN := "arn:aws:ecs:region:123:task-definition/l0-test-api:1"
	service.SetTaskDefinition(taskDefinitionARN)

	describeServicesOutput := &ecs.DescribeServicesOutput{}
	describeServicesOutput.SetServices([]*ecs.Service{service})

	mockAWS.ECS.EXPECT().
		DescribeServices(describeServicesInput).
		Return(describeServicesOutput, nil)

	target := provider.NewAdminProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Init(); err != nil {
		t.Fatal(err)
	}

	expectedTags := []models.Tag{
		{EntityID: "api", EntityType: "deploy", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "deploy", Key: "arn", Value: taskDefinitionARN},
		{EntityID: "api", EntityType: "deploy", Key: "version", Value: "1"},
		{EntityID: "api", EntityType: "environment", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "environment", Key: "os", Value: "linux"},
		{EntityID: "api", EntityType: "load_balancer", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "load_balancer", Key: "environment_id", Value: "api"},
		{EntityID: "api", EntityType: "service", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "service", Key: "environment_id", Value: "api"},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}
