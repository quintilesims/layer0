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

func TestServiceScale(t *testing.T) {
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

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.SetCluster("l0-test-env_id")
	updateServiceInput.SetDesiredCount(2)
	updateServiceInput.SetService("l0-test-svc_id")

	updateServiceOutput := &ecs.UpdateServiceOutput{}

	mockAWS.ECS.EXPECT().
		UpdateService(updateServiceInput).
		Return(updateServiceOutput, nil)

	scale := 2
	req := models.UpdateServiceRequest{
		ServiceID: "svc_id",
		Scale:     &scale,
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	err := target.Update(req)
	if err != nil {
		t.Fatal(err)
	}
}
