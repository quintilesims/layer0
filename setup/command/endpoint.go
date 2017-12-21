package command

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

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
			outputEnvVars := map[string]string{
				config.FlagEndpoint.GetName(): config.FlagEndpoint.EnvVar,
				config.FlagToken.GetName():    config.FlagToken.EnvVar,
			}

			if c.Bool("dev") {
				outputEnvVars[config.FlagInstance.GetName()] = config.FlagInstance.EnvVar
				outputEnvVars[config.FlagAWSAccountID.GetName()] = config.FlagAWSAccountID.EnvVar
				outputEnvVars[config.FlagAWSAccessKey.GetName()] = config.FlagAWSAccessKey.EnvVar
				outputEnvVars[config.FlagAWSSecretKey.GetName()] = config.FlagAWSSecretKey.EnvVar
				outputEnvVars[config.FlagAWSVPC.GetName()] = config.FlagAWSVPC.EnvVar
				outputEnvVars[config.FlagAWSLinuxAMI.GetName()] = config.FlagAWSLinuxAMI.EnvVar
				outputEnvVars[config.FlagAWSWindowsAMI.GetName()] = config.FlagAWSWindowsAMI.EnvVar
				outputEnvVars[config.FlagAWSInstanceProfile.GetName()] = config.FlagAWSInstanceProfile.EnvVar
				outputEnvVars[config.FlagAWSJobTable.GetName()] = config.FlagAWSJobTable.EnvVar
				outputEnvVars[config.FlagAWSTagTable.GetName()] = config.FlagAWSTagTable.EnvVar
				outputEnvVars[config.FlagAWSLockTable.GetName()] = config.FlagAWSLockTable.EnvVar
				outputEnvVars[config.FlagAWSPublicSubnets.GetName()] = config.FlagAWSPublicSubnets.EnvVar
				outputEnvVars[config.FlagAWSPrivateSubnets.GetName()] = config.FlagAWSPrivateSubnets.EnvVar
				outputEnvVars[config.FlagAWSLogGroup.GetName()] = config.FlagAWSLogGroup.EnvVar
				outputEnvVars[config.FlagAWSSSHKey.GetName()] = config.FlagAWSSSHKey.EnvVar
				outputEnvVars[config.FlagAWSS3Bucket.GetName()] = config.FlagAWSS3Bucket.EnvVar
			}

			// note that this requires the following relationship:
			// layer0 terraform module output name == layer0 api flag name

			fmt.Println("# set the following environment variables in your current session: ")
			for output, envVar := range outputEnvVars {
				v, err := instance.Output(output)
				if err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), envVar, v); err != nil {
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
