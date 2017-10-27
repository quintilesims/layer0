package aws

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws/mock_aws"
)

type MockClient struct {
	AutoScaling    *mock_aws.MockAutoScalingAPI
	CloudWatchLogs *mock_aws.MockCloudWatchLogsAPI
	EC2            *mock_aws.MockEC2API
	ECS            *mock_aws.MockECSAPI
	ELB            *mock_aws.MockELBAPI
	IAM            *mock_aws.MockIAMAPI
	S3             *mock_aws.MockS3API
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		AutoScaling:    mock_aws.NewMockAutoScalingAPI(ctrl),
		CloudWatchLogs: mock_aws.NewMockCloudWatchLogsAPI(ctrl),
		EC2:            mock_aws.NewMockEC2API(ctrl),
		ECS:            mock_aws.NewMockECSAPI(ctrl),
		ELB:            mock_aws.NewMockELBAPI(ctrl),
		IAM:            mock_aws.NewMockIAMAPI(ctrl),
		S3:             mock_aws.NewMockS3API(ctrl),
	}
}

func (m *MockClient) Client() *Client {
	return &Client{
		AutoScaling:    m.AutoScaling,
		CloudWatchLogs: m.CloudWatchLogs,
		EC2:            m.EC2,
		ECS:            m.ECS,
		ELB:            m.ELB,
		IAM:            m.IAM,
		S3:             m.S3,
	}
}
