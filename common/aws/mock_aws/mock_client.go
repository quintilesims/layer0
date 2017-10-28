package mock_aws

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws"
)

type MockClient struct {
	AutoScaling    *MockAutoScalingAPI
	CloudTrail     *MockCloudTrailAPI
	CloudWatchLogs *MockCloudWatchLogsAPI
	EC2            *MockEC2API
	ECS            *MockECSAPI
	ELB            *MockELBAPI
	IAM            *MockIAMAPI
	S3             *MockS3API
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		AutoScaling:    NewMockAutoScalingAPI(ctrl),
		CloudTrail:     NewMockCloudTrailAPI(ctrl),
		CloudWatchLogs: NewMockCloudWatchLogsAPI(ctrl),
		EC2:            NewMockEC2API(ctrl),
		ECS:            NewMockECSAPI(ctrl),
		ELB:            NewMockELBAPI(ctrl),
		IAM:            NewMockIAMAPI(ctrl),
		S3:             NewMockS3API(ctrl),
	}
}

func (m *MockClient) Client() *aws.Client {
	return &aws.Client{
		AutoScaling:    m.AutoScaling,
		CloudTrail:     m.CloudTrail,
		CloudWatchLogs: m.CloudWatchLogs,
		EC2:            m.EC2,
		ECS:            m.ECS,
		ELB:            m.ELB,
		IAM:            m.IAM,
		S3:             m.S3,
	}
}
