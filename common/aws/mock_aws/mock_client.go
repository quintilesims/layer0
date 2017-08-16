package mock_aws

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws"
)

type MockClient struct {
	AutoScaling *MockAutoScalingAPI
	EC2         *MockEC2API
	ECS         *MockECSAPI
	S3          *MockS3API
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		AutoScaling: NewMockAutoScalingAPI(ctrl),
		EC2:         NewMockEC2API(ctrl),
		ECS:         NewMockECSAPI(ctrl),
		S3:          NewMockS3API(ctrl),
	}
}

func (m *MockClient) Client() *aws.Client {
	return &aws.Client{
		AutoScaling: m.AutoScaling,
		EC2:         m.EC2,
		ECS:         m.ECS,
		S3:          m.S3,
	}
}
