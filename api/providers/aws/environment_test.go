package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws/mock_aws"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_createCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	validateCreateClusterInput := func(input *ecs.CreateClusterInput) {
		assert.Equal(t, "e1", aws.StringValue(input.ClusterName))
	}

	mockAWS.ECS.EXPECT().
		CreateCluster(gomock.Any()).
		Do(validateCreateClusterInput).
		Return(&ecs.CreateClusterOutput{}, nil)

	environment := NewEnvironment(mockAWS.Client(), "e1")
	if err := environment.createCluster(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironment_deleteCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	validateDeleteClusterInput := func(input *ecs.DeleteClusterInput) {
		assert.Equal(t, "e1", aws.StringValue(input.Cluster))
	}

	mockAWS.ECS.EXPECT().
		DeleteCluster(gomock.Any()).
		Do(validateDeleteClusterInput).
		Return(&ecs.DeleteClusterOutput{}, nil)

	environment := NewEnvironment(mockAWS.Client(), "e1")
	if err := environment.deleteCluster(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironment_deleteClusterDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	mockAWS.ECS.EXPECT().
		DeleteCluster(gomock.Any()).
		Return(nil, awserr.New("ClusterNotFoundException", "", nil))

	environment := NewEnvironment(mockAWS.Client(), "e1")
	if err := environment.deleteCluster(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironment_readCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	validateDescribeClustersInput := func(input *ecs.DescribeClustersInput) {
		assert.Len(t, input.Clusters, 1)
		assert.Equal(t, "e1", aws.StringValue(input.Clusters[0]))
	}

	output := &ecs.DescribeClustersOutput{
		Clusters: []*ecs.Cluster{{
			ClusterName: aws.String("e1"),
			Status:      aws.String("ACTIVE"),
		}},
	}

	mockAWS.ECS.EXPECT().
		DescribeClusters(gomock.Any()).
		Do(validateDescribeClustersInput).
		Return(output, nil)

	environment := NewEnvironment(mockAWS.Client(), "e1")
	cluster, err := environment.readCluster()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output.Clusters[0], cluster)
}

func TestEnvironment_readClusterDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := mock_aws.NewMockClient(ctrl)
	mockAWS.ECS.EXPECT().
		DescribeClusters(gomock.Any()).
		Return(&ecs.DescribeClustersOutput{}, nil)

	environment := NewEnvironment(mockAWS.Client(), "e1")
	_, err := environment.readCluster()
	assertEqualErrorCode(t, err, errors.EnvironmentDoesNotExist)
}
