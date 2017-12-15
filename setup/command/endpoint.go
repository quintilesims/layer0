package command

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type output struct {
	TerraformOutput string
	EnvVar          string
}

func (f *CommandFactory) Endpoint() cli.Command {
	return cli.Command{
		Name:      "endpoint",
		Usage:     "show environment variables used to connect to a Layer0 instance",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "i, insecure",
				Usage: "show environment variables that allow for insecure settings",
			},
			cli.BoolFlag{
				Name:  "d, dev",
				Usage: "show environment variables that allow for local development",
			},
			cli.StringFlag{
				Name:  "s, syntax",
				Value: "bash",
				Usage: "choose the syntax to display environment variables (choices: bash, cmd, powershell)",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			instance := f.NewInstance(args["NAME"])
			outputs := []output{
				{"endpoint", config.FlagEndpoint.EnvVar},
				{"token", config.FlagToken.EnvVar},
			}

			// todo: not use hardocded strings in instance.Apply and instance.Push
			// todo: ensure module outputs match
			if c.Bool("dev") {
				devOutputs := []output{
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
				}

				outputs = append(outputs, devOutputs...)
			}

			fmt.Println("# set the following environment variables in your current session: ")
			for _, o := range outputs {
				v, err := instance.Output(o.TerraformOutput)
				if err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), o.EnvVar, v); err != nil {
					return err
				}
			}

			if c.Bool("insecure") {
				if err := printOutput(c.String("syntax"), config.FlagSkipVerifySSL.EnvVar, "true"); err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), config.FlagSkipVerifyVersion.EnvVar, "true"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func printOutput(syntax, envvar, v string) error {
	switch syntax {
	case "bash":
		fmt.Printf("export %s=\"%s\"\n", envvar, v)
	case "cmd":
		fmt.Printf("set %s=%s\n", envvar, v)
	case "powershell":
		fmt.Printf("$env:%s=\"%s\"\n", envvar, v)
	default:
		return fmt.Errorf("Unknown syntax '%s'. Please specify 'bash', 'cmd', or 'powershell'", syntax)
	}

	return nil
}
