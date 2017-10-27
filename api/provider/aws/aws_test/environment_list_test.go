package aws_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws"
)

func TestEnvironmentList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := aws.NewMockClient(ctrl)
	print(mockClient)
}
