package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
	"github.com/urfave/cli"
)

func TestEndpoint(t *testing.T) {
	testEndpointHelper(t, map[string]interface{}{"syntax": "bash"}, []cli.Flag{
		config.FlagEndpoint,
		config.FlagToken,
	})
}

func TestEndpointDev(t *testing.T) {
	testEndpointHelper(t, map[string]interface{}{"syntax": "bash", "dev": true}, []cli.Flag{
		config.FlagEndpoint,
		config.FlagToken,
		config.FlagInstance,
		config.FlagAWSAccountID,
		config.FlagAWSAccessKey,
		config.FlagAWSSecretKey,
		config.FlagAWSVPC,
		config.FlagAWSLinuxAMI,
		config.FlagAWSWindowsAMI,
		config.FlagAWSInstanceProfile,
		config.FlagAWSJobTable,
		config.FlagAWSTagTable,
		config.FlagAWSLockTable,
		config.FlagAWSPublicSubnets,
		config.FlagAWSPrivateSubnets,
		config.FlagAWSLogGroup,
		config.FlagAWSSSHKey,
		config.FlagAWSS3Bucket,
	})
}

func testEndpointHelper(t *testing.T, flags map[string]interface{}, expectedOutputFlags []cli.Flag) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)
		for _, flag := range expectedOutputFlags {
			mockInstance.EXPECT().
				Output(flag.GetName()).
				Return("", nil)
		}

		return mockInstance
	}

	commandFactory := NewCommandFactory(instanceFactory, nil)
	action := extractAction(t, commandFactory.Endpoint())
	c := config.NewTestContext(t, []string{"name"}, flags)
	if err := action(c); err != nil {
		t.Fatal(err)
	}
}
