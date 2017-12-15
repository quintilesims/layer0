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
	c := mock_config.NewMockAPIConfig(ctrl)

	c.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
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
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	updateServiceScaleInput := &ecs.UpdateServiceInput{}
	updateServiceScaleInput.SetCluster("l0-test-env_id")
	updateServiceScaleInput.SetDesiredCount(2)
	updateServiceScaleInput.SetService("l0-test-svc_id")

	mockAWS.ECS.EXPECT().
		UpdateService(updateServiceScaleInput).
		Return(&ecs.UpdateServiceOutput{}, nil)

	updateServiceTaskDefinitionInput := &ecs.UpdateServiceInput{}
	updateServiceTaskDefinitionInput.SetCluster("l0-test-env_id")
	updateServiceTaskDefinitionInput.SetService("l0-test-svc_id")
	updateServiceTaskDefinitionInput.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")

	mockAWS.ECS.EXPECT().
		UpdateService(updateServiceTaskDefinitionInput).
		Return(&ecs.UpdateServiceOutput{}, nil)

	deployID := "dpl_id"
	scale := 2
	req := models.UpdateServiceRequest{
		DeployID: &deployID,
		Scale:    &scale,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, c)
	if err := target.Update("svc_id", req); err != nil {
		t.Fatal(err)
	}
}
