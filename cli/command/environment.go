package command

import (
	"io/ioutil"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type EnvironmentCommand struct {
	*CommandMediator
}

func NewEnvironment(m *CommandMediator) cli.Command {
	e := &EnvironmentCommand{
		CommandMediator: m,
	}

	return cli.Command{
		Name:  "environment",
		Usage: "manage layer0 environments",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new environment",
				Action:    e.createEnvironment,
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "size",
						Value: "m3.medium",
						Usage: "size of the ec2 instances to use in the environment cluster",
					},
					cli.IntFlag{
						Name:  "min-count",
						Value: 0,
						Usage: "minimum number of instances allowed in the environment cluster",
					},
					cli.StringFlag{
						Name:  "user-data",
						Usage: "path to user data file",
					},
					cli.StringFlag{
						Name:  "os",
						Value: "linux",
						Usage: "specifies if the environment will run windows or linux containers",
					},
					cli.StringFlag{
						Name:  "ami",
						Usage: "specifies a custom AMI ID to use in the environment",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for the job to finish",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete an environment",
				ArgsUsage: "NAME",
				Action:    e.deleteEnvironment,
				Flags: []cli.Flag{
					cli.BoolTFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "list",
				Usage:     "list all environments",
				Action:    e.listEnvironments,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "describe an environment",
				Action:    e.readEnvironment,
				ArgsUsage: "NAME",
			},
		},
	}
}

func (e *EnvironmentCommand) createEnvironment(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	var userData []byte
	if path := c.String("user-data"); path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		userData = content
	}

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  args["NAME"],
		InstanceSize:     c.String("size"),
		MinClusterCount:  c.Int("min-count"),
		UserDataTemplate: userData,
		OperatingSystem:  c.String("os"),
		AMIID:            c.String("ami"),
	}

	jobID, err := e.client.CreateEnvironment(req)
	if err != nil {
		return err
	}

	if c.GlobalBool(config.FLAG_NO_WAIT) {
		// todo: use common helper
		e.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	e.printer.StartSpinner("creating")
	defer e.printer.StopSpinner()

	job, err := client.WaitForJob(e.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT))
	if err != nil {
		return err
	}

	environmentID := job.Result
	environment, err := e.client.ReadEnvironment(environmentID)
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironments(environment)
}

func (e *EnvironmentCommand) deleteEnvironment(c *cli.Context) error {
	return e.deleteHelper(c, "environment", func(environmentID string) (string, error) {
		return e.client.DeleteEnvironment(environmentID)
	})
}

func (e *EnvironmentCommand) listEnvironments(c *cli.Context) error {
	environmentSummaries, err := e.client.ListEnvironments()
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironmentSummaries(environmentSummaries...)
}

func (e *EnvironmentCommand) readEnvironment(c *cli.Context) error {
	return nil
}
