package test_aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"
	"github.com/quintilesims/layer0/common/models"
)

func TestServiceUpdate(t *testing.T) {
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
		{
			EntityID:   "svc_id",
			EntityType: "service",
			Key:        "deploy_id",
			Value:      "dpl_id",
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
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task-definition/dpl_id:1",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "version",
			Value:      "2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.SetCluster("l0-test-env_id")
	updateServiceInput.SetService("l0-test-svc_id")
	updateServiceInput.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")

	updateServiceOutput := &ecs.UpdateServiceOutput{}

	mockAWS.ECS.EXPECT().
		UpdateService(updateServiceInput).
		Return(updateServiceOutput, nil)

	deployID := "dpl_id"
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		DeployID:  &deployID,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	err := target.Update(req)
	if err != nil {
		t.Fatal(err)
	}
}
