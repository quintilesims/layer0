package test_aws

import (
	"testing"

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

func TestDeployDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
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
			Value:      "dpl_version",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	deregisterTaskInput := &ecs.DeregisterTaskDefinitionInput{}
	deregisterTaskInput.SetTaskDefinition("arn:aws:ecs:region:012345678910:task/arn")

	mockAWS.ECS.EXPECT().
		DeregisterTaskDefinition(deregisterTaskInput).
		Return(&ecs.DeregisterTaskDefinitionOutput{}, nil)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("dpl_id"); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tagStore.Tags(), 0)
}

func TestDeleteDeployIdempotence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
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
			Value:      "dpl_version",
		},
		{
			EntityID:   "dpl_id",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task/arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	mockAWS.ECS.EXPECT().
		DeregisterTaskDefinition(gomock.Any()).
		Return(nil, awserr.New("", "task definition does not exist", nil))

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("dpl_id"); err != nil {
		t.Fatal(err)
	}
}

func TestDeployDelete_idempotenceViaTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	if err := target.Delete("dpl_id"); err != nil {
		t.Fatal(err)
	}
}
