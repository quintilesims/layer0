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
			// {
			// Name:      "create",
			// Usage:     "create a new environment",
			// Action:    f.createEnvironment,
			// ArgsUsage: "NAME",
			// Flags: []cli.Flag{
			// cli.StringFlag{
			// Name:  "size",
			// Value: "m3.medium",
			// Usage: "size of the ec2 instances to use in the environment cluster",
			// },
			// cli.IntFlag{
			// Name:  "min-count",
			// Value: 0,
			// Usage: "minimum number of instances allowed in the environment cluster",
			// },
			// cli.StringFlag{
			// Name:  "user-data",
			// Usage: "path to user data file",
			// },
			// cli.StringFlag{
			// Name:  "os",
			// Value: "linux",
			// Usage: "specifies if the environment will run windows or linux containers",
			// },
			// cli.StringFlag{
			// Name:  "ami",
			// Usage: "specifies a custom AMI ID to use in the environment",
			// },
			// cli.BoolFlag{
			// Name:  "nowait",
			// Usage: "don't wait for the job to finish",
			// },
			// },
			// },
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

	if c.Bool("nowait") {
		f.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	timeout := c.GlobalDuration(config.FLAG_TIMEOUT)

	job, err := client.WaitForJob(f.client, jobID, timeout)
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
	// hint: get timeout via c.GlobalDuration(config.FLAG_TIMEOUT)
	// then pass to: client.WaitForJob(f.client, jobID, timeout)

	// get args
	//  - resolve
	//  - entityType

	f.deleteHelper(c, "environment", func(environmentID string) (string, error) {
		return f.client.DeleteEnvironment(environmentID)
	})

	job, err := client.WaitForJob(f.client, jobID, timeout)

	return nil
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

// put this somewhere better for common funcs
func (f *CommandFactory) deleteHelper(c *cli.Context, entityType string, fn func(entityID string) (string, error)) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	entityID, err := resolveSingleID(entityType, args["NAME"])
	if err != nil {
		return err
	}

	jobID, err := fn(entityID)
	if err != nil {
		return err
	}

	if c.Bool("nowait") {
		f.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	f.printer.StartSpinner("Deleting")
	defer f.printer.StopSpinner()

	timeout := c.GlobalDuration(config.FLAG_TIMEOUT)
	if _, err := client.WaitForJob(f.client, jobID, timeout); err != nil {
		return err
	}

	return nil
}
