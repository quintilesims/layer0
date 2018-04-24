package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)

		outputs := []string{
			instance.OUTPUT_ENDPOINT,
			instance.OUTPUT_TOKEN,
		}

		for _, output := range outputs {
			mockInstance.EXPECT().
				Output(output).
				Return("", nil)
		}

		return mockInstance
	}

	input := "l0-setup endpoint name"
	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Endpoint(), input); err != nil {
		t.Fatal(err)
	}
}

func TestEndpointDev(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)

		outputs := []string{
			instance.OUTPUT_ENDPOINT,
			instance.OUTPUT_TOKEN,
			instance.OUTPUT_NAME,
			instance.OUTPUT_ACCOUNT_ID,
			instance.OUTPUT_ACCESS_KEY,
			instance.OUTPUT_SECRET_KEY,
			instance.OUTPUT_SSH_KEY_PAIR,
			instance.OUTPUT_AWS_LOG_GROUP_NAME,
			instance.OUTPUT_VPC_ID,
			instance.OUTPUT_PRIVATE_SUBNETS,
			instance.OUTPUT_PUBLIC_SUBNETS,
			instance.OUTPUT_S3_BUCKET,
			instance.OUTPUT_ECS_INSTANCE_PROFILE,
			instance.OUTPUT_AWS_LINUX_SERVICE_AMI,
			instance.OUTPUT_AWS_DYNAMO_TAG_TABLE,
			instance.OUTPUT_AWS_DYNAMO_LOCK_TABLE,
			instance.OUTPUT_AWS_REGION,
		}

		for _, output := range outputs {
			mockInstance.EXPECT().
				Output(output).
				Return("", nil)
		}

		return mockInstance
	}

	input := "l0-setup endpoint "
	input += "--dev "
	input += "name"

	factory := NewCommandFactory(instanceFactory, nil)
	if err := testutils.RunApp(factory.Endpoint(), input); err != nil {
		t.Fatal(err)
	}
}
