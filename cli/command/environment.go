package command

import (
	"io/ioutil"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

func (f *CommandFactory) Environment() cli.Command {
	return cli.Command{
		Name:  "environment",
		Usage: "manage layer0 environments",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new environment",
				Action:    f.createEnvironment,
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
				Action:    f.deleteEnvironment,
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
				Action:    f.listEnvironments,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "describe an environment",
				Action:    f.readEnvironment,
				ArgsUsage: "NAME",
			},
		},
	}
}

func (f *CommandFactory) createEnvironment(c *cli.Context) error {
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

	jobID, err := f.client.CreateEnvironment(req)
	if err != nil {
		return err
	}

	if c.GlobalBool(config.FLAG_NO_WAIT) {
		// todo: use common helper
		f.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	f.printer.StartSpinner("creating")
	defer f.printer.StopSpinner()

	job, err := client.WaitForJob(f.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT))
	if err != nil {
		return err
	}

	environmentID := job.Result
	environment, err := f.client.ReadEnvironment(environmentID)
	if err != nil {
		return err
	}

	return f.printer.PrintEnvironments(environment)
}

func (f *CommandFactory) deleteEnvironment(c *cli.Context) error {
	return f.deleteHelper(c, "environment", func(environmentID string) (string, error) {
		return f.client.DeleteEnvironment(environmentID)
	})
}

func (f *CommandFactory) listEnvironments(c *cli.Context) error {
	environmentSummaries, err := f.client.ListEnvironments()
	if err != nil {
		return err
	}

	return f.printer.PrintEnvironmentSummaries(environmentSummaries...)
}

func (f *CommandFactory) readEnvironment(c *cli.Context) error {
	return nil
}
