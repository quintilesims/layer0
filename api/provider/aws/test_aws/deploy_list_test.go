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

func TestDeployList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "dpl_id1",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl1",
		},
		{
			EntityID:   "dpl_id1",
			EntityType: "deploy",
			Key:        "version",
			Value:      "1",
		},
		{
			EntityID:   "dpl_id1",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl1:1",
		},
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl2",
		},
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "version",
			Value:      "1",
		},
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl2:1",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	taskDefinitionFamilies := []*string{
		aws.String("l0-test-dpl_id1"),
		aws.String("l0-test-dpl_id2"),
		aws.String("l0-bad-dpl_id1"),
	}

	// ListTaskDefinitionsFamiliesPages Mocks
	listTaskDefinitionFamiliesPagesfn := func(input *ecs.ListTaskDefinitionFamiliesInput,
		fn func(output *ecs.ListTaskDefinitionFamiliesOutput, lastPage bool) bool) error {

		output := &ecs.ListTaskDefinitionFamiliesOutput{}
		output.SetFamilies(taskDefinitionFamilies)
		fn(output, true)

		return nil
	}

	tdFamilies := &ecs.ListTaskDefinitionFamiliesInput{}
	tdFamilies.SetFamilyPrefix("l0-test-")
	tdFamilies.SetStatus(ecs.TaskDefinitionFamilyStatusActive)

	mockAWS.ECS.EXPECT().
		ListTaskDefinitionFamiliesPages(tdFamilies, gomock.Any()).
		Do(listTaskDefinitionFamiliesPagesfn).
		Return(nil)

	// ListTaskDefinitionsPages Mocks
	taskDefinitionARNs := []*string{
		aws.String("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id1:1"),
		aws.String("arn:aws:ecs:region:012345678910:task-definition/l0-test-dpl_id2:1"),
		aws.String("arn:aws:ecs:region:012345678910:task-definition/l0-bad-dpl_id1:1"),
	}

	generateTaskDefinitionPagesFN := func(taskDefinitionARN *string) func(input *ecs.ListTaskDefinitionsInput, fn func(output *ecs.ListTaskDefinitionsOutput, lastPage bool) bool) error {
		listTaskDefinitionPagesFN := func(input *ecs.ListTaskDefinitionsInput, fn func(output *ecs.ListTaskDefinitionsOutput, lastPage bool) bool) error {
			output := &ecs.ListTaskDefinitionsOutput{}
			output.SetTaskDefinitionArns([]*string{taskDefinitionARN})

			fn(output, true)

			return nil
		}

		return listTaskDefinitionPagesFN
	}

	td := &ecs.ListTaskDefinitionsInput{}
	td.SetFamilyPrefix("l0-test-dpl_id")
	td.SetStatus(ecs.TaskDefinitionStatusActive)

	for i := range taskDefinitionFamilies {
		mockAWS.ECS.EXPECT().
			ListTaskDefinitionsPages(td, gomock.Any()).
			Do(generateTaskDefinitionPagesFN(taskDefinitionARNs[i])).
			Return(nil)
	}

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.DeploySummary{
		{
			DeployID:   "dpl_id1",
			DeployName: "dpl1",
			Version:    "1",
		},
		{
			DeployID:   "dpl_id2",
			DeployName: "dpl2",
			Version:    "1",
		},
	}

	assert.Equal(t, expected, result)
}
