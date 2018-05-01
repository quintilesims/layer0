package command

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
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

			outputEnvvars := map[string]string{
				instance.OUTPUT_ENDPOINT: config.ENVVAR_ENDPOINT,
				instance.OUTPUT_TOKEN:    config.ENVVAR_TOKEN,
			}

			if c.Bool("dev") {
				outputEnvvars[instance.OUTPUT_NAME] = config.ENVVAR_INSTANCE
				outputEnvvars[instance.OUTPUT_ACCOUNT_ID] = config.ENVVAR_AWS_ACCOUNT_ID
				outputEnvvars[instance.OUTPUT_ACCESS_KEY] = config.ENVVAR_AWS_ACCESS_KEY
				outputEnvvars[instance.OUTPUT_SECRET_KEY] = config.ENVVAR_AWS_SECRET_KEY
				outputEnvvars[instance.OUTPUT_VPC_ID] = config.ENVVAR_AWS_VPC
				outputEnvvars[instance.OUTPUT_PRIVATE_SUBNETS] = config.ENVVAR_AWS_PRIVATE_SUBNETS
				outputEnvvars[instance.OUTPUT_PUBLIC_SUBNETS] = config.ENVVAR_AWS_PUBLIC_SUBNETS
				outputEnvvars[instance.OUTPUT_S3_BUCKET] = config.ENVVAR_AWS_S3_BUCKET
				outputEnvvars[instance.OUTPUT_SSH_KEY_PAIR] = config.ENVVAR_AWS_SSH_KEY_PAIR
				outputEnvvars[instance.OUTPUT_AWS_LOG_GROUP_NAME] = config.ENVVAR_AWS_LOG_GROUP_NAME
				outputEnvvars[instance.OUTPUT_ECS_INSTANCE_PROFILE] = config.ENVVAR_AWS_INSTANCE_PROFILE
				outputEnvvars[instance.OUTPUT_AWS_LINUX_SERVICE_AMI] = config.ENVVAR_AWS_LINUX_AMI
				outputEnvvars[instance.OUTPUT_AWS_DYNAMO_TAG_TABLE] = config.ENVVAR_AWS_DYNAMO_TAG_TABLE
				outputEnvvars[instance.OUTPUT_AWS_DYNAMO_LOCK_TABLE] = config.ENVVAR_AWS_DYNAMO_LOCK_TABLE
				outputEnvvars[instance.OUTPUT_AWS_REGION] = config.ENVVAR_AWS_REGION
			}

			fmt.Println("# set the following environment variables in your current session: ")

			instance := f.NewInstance(args["NAME"])
			for output, envvar := range outputEnvvars {
				v, err := instance.Output(output)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if err := printOutput(c.String("syntax"), envvar, v); err != nil {
					return err
				}
			}

			if c.Bool("insecure") {
				if err := printOutput(c.String("syntax"), config.ENVVAR_SKIP_VERIFY_SSL, "1"); err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), config.ENVVAR_SKIP_VERIFY_VERSION, "1"); err != nil {
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
