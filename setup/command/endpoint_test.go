package command

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/quintilesims/layer0/setup/instance/mock_instance"
)

func TestEndpoint(t *testing.T) {
	cases := []struct {
		Name            string
		Flags           map[string]interface{}
		ExpectedOutputs []output
	}{
		{
			Name:  "no flags",
			Flags: map[string]interface{}{},
			ExpectedOutputs: []output{
				{"endpoint", config.FlagEndpoint.EnvVar},
				{"token", config.FlagToken.EnvVar},
			},
		},
		{
			Name:  "dev flag",
			Flags: map[string]interface{}{"dev": true},
			ExpectedOutputs: []output{
				{"endpoint", config.FlagEndpoint.EnvVar},
				{"token", config.FlagToken.EnvVar},
				{"instance", config.FlagInstance.EnvVar},
				{"aws_account_id", config.FlagAWSAccountID.EnvVar},
				{"aws_access_key", config.FlagAWSAccessKey.EnvVar},
				{"aws_secret_key", config.FlagAWSSecretKey.EnvVar},
				{"aws_vpc", config.FlagAWSVPC.EnvVar},
				{"aws_linux_ami", config.FlagAWSLinuxAMI.EnvVar},
				{"aws_windows_ami", config.FlagAWSWindowsAMI.EnvVar},
				{"aws_s3_bucket", config.FlagAWSS3Bucket.EnvVar},
				{"aws_instance_profile", config.FlagAWSInstanceProfile.EnvVar},
				{"aws_job_table", config.FlagAWSJobTable.EnvVar},
				{"aws_tag_table", config.FlagAWSTagTable.EnvVar},
				{"aws_lock_table", config.FlagAWSLockTable.EnvVar},
				{"aws_public_subnets", config.FlagAWSPublicSubnets.EnvVar},
				{"aws_private_subnets", config.FlagAWSPrivateSubnets.EnvVar},
				{"aws_log_group", config.FlagAWSLogGroup.EnvVar},
				{"aws_ssh_key", config.FlagAWSSSHKey.EnvVar},
			},
		},
	}

	for _, c := range cases {
		c.Flags["syntax"] = "bash"

		t.Run(c.Name, func(t *testing.T) {
			testEndpointHelper(t, c.ExpectedOutputs, c.Flags)
		})
	}
}

func testEndpointHelper(t *testing.T, expectedOutputs []output, flags map[string]interface{}) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instanceFactory := func(name string) instance.Instance {
		mockInstance := mock_instance.NewMockInstance(ctrl)

		for _, o := range expectedOutputs {
			mockInstance.EXPECT().
				Output(o.TerraformOutput).
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
