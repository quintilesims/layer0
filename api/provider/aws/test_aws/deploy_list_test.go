package test_aws

import (
	"testing"

	// "github.com/aws/aws-sdk-go/aws"
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
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "name",
			Value:      "dpl_name2",
		},
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "version",
			Value:      "dpl_version2",
		},
		{
			EntityID:   "dpl_id2",
			EntityType: "deploy",
			Key:        "arn",
			Value:      "arn:aws:ecs:region:012345678911:task/arn",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// taskDefinitionFamilies := []string{}
	// listTaskDefinitionFamiliesPagesfn := func(output *ecs.ListTaskDefinitionFamiliesOutput, lastPage bool) bool {
	// 	for _, taskDefinitionFamily := range output.Families {
	// 		taskDefinitionFamilies = append(taskDefinitionFamilies, aws.StringValue(taskDefinitionFamily))
	// 	}
	// 	return !lastPage
	// }

	tdFamilies := &ecs.ListTaskDefinitionFamiliesInput{}
	tdFamilies.SetFamilyPrefix("l0-test-")
	// TODO: for some reason, ListTask...() call panics when SetStatus == "ACTIVE"
	tdFamilies.SetStatus(ecs.TaskDefinitionFamilyStatusActive)

	mockAWS.ECS.EXPECT().
		ListTaskDefinitionFamiliesPages(tdFamilies, gomock.Any()).
		// Do(listTaskDefinitionFamiliesPagesfn).
		Return(nil)

	// taskDefinitionARNs := []string{}
	// listTaskDefinitionPagesfn := func(output *ecs.ListTaskDefinitionsOutput, lastPage bool) bool {
	// 	for _, taskDefinitionARN := range output.TaskDefinitionArns {
	// 		taskDefinitionARNs = append(taskDefinitionARNs, aws.StringValue(taskDefinitionARN))
	// 	}
	// 	return !lastPage
	// }

	td := &ecs.ListTaskDefinitionsInput{}
	td.SetFamilyPrefix("l0-test-")
	td.SetStatus(ecs.TaskDefinitionFamilyStatusActive)

	mockAWS.ECS.EXPECT().
		ListTaskDefinitionsPages(td, gomock.Any()).
		// Do(listTaskDefinitionPagesfn).
		Return(nil)

	target := provider.NewDeployProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.DeploySummary{
		{
			DeployID:   "dpl_id",
			DeployName: "dpl_name",
			Version:    "dpl_version",
		},
		{
			DeployID:   "dpl_id2",
			DeployName: "dpl_name2",
			Version:    "dpl_version2",
		},
	}

	assert.Equal(t, expected, result)
}
