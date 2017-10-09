package command

import (
	"github.com/urfave/cli"
)

func (f *CommandFactory) Environment() cli.Command {
	return cli.Command{
		Name:  "environment",
		Usage: "manage layer0 environments",
		Subcommands: []cli.Command{
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

func (f *CommandFactory) deleteEnvironment(c *cli.Context) error {
	// hint: get timeout via c.GlobalDuration(config.FLAG_TIMEOUT)
	// then pass to: client.WaitForJob(f.client, jobID, timeout)
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
