package command

import (
	"testing"

	"github.com/golang/mock/gomock"
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

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Endpoint())

	flags := map[string]interface{}{
		"syntax": "bash",
	}

	c := NewContext(t, []string{"name"}, flags)
	if err := action(c); err != nil {
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
			instance.OUTPUT_WINDOWS_SERVICE_AMI,
			instance.OUTPUT_AWS_DYNAMO_TAG_TABLE,
			instance.OUTPUT_AWS_DYNAMO_JOB_TABLE,
		}

		for _, output := range outputs {
			mockInstance.EXPECT().
				Output(output).
				Return("", nil)
		}

		return mockInstance
	}

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Endpoint())

	flags := map[string]interface{}{
		"syntax": "bash",
		"dev":    "true",
	}

	c := NewContext(t, []string{"name"}, flags)
	if err := action(c); err != nil {
		t.Fatal(err)
	}
}
