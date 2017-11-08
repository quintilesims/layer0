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

func TestServiceList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	tagStore := tag.NewMemoryStore()
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()

	tags := models.Tags{
		{
			EntityID:   "env_id1",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name1",
		},
		{
			EntityID:   "env_id2",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name2",
		},
		{
			EntityID:   "svc_id1",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name1",
		},
		{
			EntityID:   "svc_id1",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id1",
		},
		{
			EntityID:   "svc_id2",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name2",
		},
		{
			EntityID:   "svc_id2",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id1",
		},
		{
			EntityID:   "svc_id3",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name3",
		},
		{
			EntityID:   "svc_id3",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id2",
		},
		{
			EntityID:   "svc_id4",
			EntityType: "service",
			Key:        "name",
			Value:      "svc_name4",
		},
		{
			EntityID:   "svc_id4",
			EntityType: "service",
			Key:        "environment_id",
			Value:      "env_id2",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	clusterARNs := []*string{
		aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id1"),
		aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id2"),
		aws.String("arn:aws:ecs:region:012345678910:cluster/l0-bad-env_id1"),
		aws.String("arn:aws:ecs:region:012345678910:cluster/bad"),
	}

	listClustersPagesFN := func(input *ecs.ListClustersInput, fn func(output *ecs.ListClustersOutput, lastPage bool) bool) error {
		output := &ecs.ListClustersOutput{}
		output.SetClusterArns(clusterARNs)
		fn(output, true)

		return nil
	}

	mockAWS.ECS.EXPECT().
		ListClustersPages(&ecs.ListClustersInput{}, gomock.Any()).
		Do(listClustersPagesFN).
		Return(nil)

	clusterServices := map[string][]*string{
		"l0-test-env_id1": []*string{
			aws.String("arn:aws:ecs:region:012345678910:service/l0-test-svc_id1"),
			aws.String("arn:aws:ecs:region:012345678910:service/l0-test-svc_id2"),
		},
		"l0-test-env_id2": []*string{
			aws.String("arn:aws:ecs:region:012345678910:service/l0-test-svc_id3"),
			aws.String("arn:aws:ecs:region:012345678910:service/l0-test-svc_id4"),
		},
	}

	generateListServicesPagesFN := func(serviceARNs []*string) func(input *ecs.ListServicesInput, fn func(output *ecs.ListServicesOutput, lastPage bool) bool) error {
		listServicesPagesFN := func(input *ecs.ListServicesInput, fn func(output *ecs.ListServicesOutput, lastPage bool) bool) error {
			output := &ecs.ListServicesOutput{}
			output.SetServiceArns(serviceARNs)

			fn(output, true)

			return nil
		}

		return listServicesPagesFN
	}

	for clusterName, serviceARNs := range clusterServices {
		input := &ecs.ListServicesInput{}
		input.SetCluster(clusterName)

		mockAWS.ECS.EXPECT().
			ListServicesPages(input, gomock.Any()).
			Do(generateListServicesPagesFN(serviceARNs)).
			Return(nil)
	}

	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, mockConfig)
	result, err := target.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := []models.ServiceSummary{
		{
			ServiceID:       "svc_id1",
			ServiceName:     "svc_name1",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			ServiceID:       "svc_id2",
			ServiceName:     "svc_name2",
			EnvironmentID:   "env_id1",
			EnvironmentName: "env_name1",
		},
		{
			ServiceID:       "svc_id3",
			ServiceName:     "svc_name3",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
		{
			ServiceID:       "svc_id4",
			ServiceName:     "svc_name4",
			EnvironmentID:   "env_id2",
			EnvironmentName: "env_name2",
		},
	}

	assert.Equal(t, expected, result)
}
