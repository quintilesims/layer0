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
		Usage:     "endpoint layer0 instances",
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
				/*
					todo: include all required outputs:
					settings["account_id"] = config.AWS_ACCOUNT_ID
					settings["key_pair"] = config.AWS_KEY_PAIR
					settings["agent_security_group_id"] = config.AWS_ECS_AGENT_SECURITY_GROUP_ID
					settings["ecs_instance_profile"] = config.AWS_ECS_INSTANCE_PROFILE
					settings["ecs_role"] = config.AWS_ECS_ROLE
					settings["public_subnets"] = config.AWS_PUBLIC_SUBNETS
					settings["private_subnets"] = config.AWS_PRIVATE_SUBNETS
					settings["access_key"] = config.AWS_ACCESS_KEY_ID
					settings["secret_key"] = config.AWS_SECRET_ACCESS_KEY
					settings["region"] = config.AWS_REGION
					settings["l0_prefix"] = config.PREFIX
					settings["runner_docker_image_tag"] = config.RUNNER_VERSION_TAG
					settings["vpc_id"] = config.AWS_VPC_ID
					settings["s3_bucket"] = config.AWS_S3_BUCKET
					settings["linux_service_ami"] = config.AWS_LINUX_SERVICE_AMI
					settings["windows_service_ami"] = config.AWS_WINDOWS_SERVICE_AMI
					settings["dynamo_tag_table"] = config.AWS_DYNAMO_TAG_TABLE
					settings["dynamo_job_table"] = config.AWS_DYNAMO_JOB_TABLE
				*/

			}

			instance := instance.NewInstance(args["NAME"])
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
