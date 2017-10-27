package aws_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws/mock_aws"
)

func TestEnvironmentList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_aws.NewMockClient(ctrl)
	print(mockClient)
}
