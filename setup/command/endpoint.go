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
		Usage:     "todo",
		ArgsUsage: "NAME",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "i, insecure",
			},
			cli.BoolFlag{
				Name: "d, dev",
			},
			cli.StringFlag{
				Name:  "s, syntax",
				Value: "bash",
			},
		},
		Action: func(c *cli.Context) error {
			args, err := extractArgs(c.Args(), "NAME")
			if err != nil {
				return err
			}

			outputEnvvars := map[string]string{
				instance.OUTPUT_ENDPOINT: config.API_ENDPOINT,
				instance.OUTPUT_TOKEN:    config.AUTH_TOKEN,
			}

			if c.Bool("dev") {
				outputEnvvars[instance.OUTPUT_NAME] = config.PREFIX
				outputEnvvars[instance.OUTPUT_ACCOUNT_ID] = config.AWS_ACCOUNT_ID
				outputEnvvars[instance.OUTPUT_ACCESS_KEY] = config.AWS_ACCESS_KEY_ID
				outputEnvvars[instance.OUTPUT_SECRET_KEY] = config.AWS_SECRET_ACCESS_KEY
				outputEnvvars[instance.OUTPUT_VPC_ID] = config.AWS_VPC_ID
				outputEnvvars[instance.OUTPUT_PRIVATE_SUBNETS] = config.AWS_PRIVATE_SUBNETS
				outputEnvvars[instance.OUTPUT_PUBLIC_SUBNETS] = config.AWS_PUBLIC_SUBNETS
				outputEnvvars[instance.OUTPUT_ECS_ROLE] = config.AWS_ECS_ROLE
				outputEnvvars[instance.OUTPUT_KEY_PAIR] = config.AWS_KEY_PAIR
				outputEnvvars[instance.OUTPUT_S3_BUCKET] = config.AWS_S3_BUCKET
				outputEnvvars[instance.OUTPUT_ECS_AGENT_SECURITY_GROUP_ID] = config.AWS_ECS_AGENT_SECURITY_GROUP_ID
				outputEnvvars[instance.OUTPUT_ECS_INSTANCE_PROFILE] = config.AWS_ECS_INSTANCE_PROFILE
				outputEnvvars[instance.OUTPUT_AWS_LINUX_SERVICE_AMI] = config.AWS_LINUX_SERVICE_AMI
				outputEnvvars[instance.OUTPUT_WINDOWS_SERVICE_AMI] = config.AWS_WINDOWS_SERVICE_AMI
			}

			instance := f.NewInstance(args["NAME"])
			for output, envvar := range outputEnvvars {
				v, err := instance.Output(output)
				if err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), envvar, v); err != nil {
					return err
				}
			}

			if c.Bool("insecure") {
				if err := printOutput(c.String("syntax"), config.SKIP_SSL_VERIFY, "1"); err != nil {
					return err
				}

				if err := printOutput(c.String("syntax"), config.SKIP_VERSION_VERIFY, "1"); err != nil {
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
