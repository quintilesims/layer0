package command

import (
	"io/ioutil"
	"strconv"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type EnvironmentCommand struct {
	*Command
}

func NewEnvironmentCommand(command *Command) *EnvironmentCommand {
	return &EnvironmentCommand{command}
}

func (e *EnvironmentCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "environment",
		Usage: "manage layer0 environments",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new environment",
				Action:    wrapAction(e.Command, e.Create),
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
						Usage: "specifies windows or linux container environment cluster",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete an environment",
				ArgsUsage: "NAME",
				Action:    wrapAction(e.Command, e.Delete),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "get",
				Usage:     "describe an environment",
				Action:    wrapAction(e.Command, e.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all environments",
				Action:    wrapAction(e.Command, e.List),
				ArgsUsage: " ",
			},
			{
				Name:      "setmincount",
				Usage:     "set the minimum instance count for an environment cluster",
				Action:    wrapAction(e.Command, e.SetMinCount),
				ArgsUsage: "NAME COUNT",
			},
		},
	}
}

func (e *EnvironmentCommand) Create(c *cli.Context) error {
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

	//todo: add it here
	environment, err := e.Client.CreateEnvironment(args["NAME"], c.String("size"), c.Int("min-count"), userData, c.String("os"))
	if err != nil {
		return err
	}

	return e.Printer.PrintEnvironments(environment)
}

func (e *EnvironmentCommand) Delete(c *cli.Context) error {
	return e.deleteWithJob(c, "environment", e.Client.DeleteEnvironment)
}

func (e *EnvironmentCommand) Get(c *cli.Context) error {
	environments := []*models.Environment{}
	getEnvironmentf := func(id string) error {
		environment, err := e.Client.GetEnvironment(id)
		if err != nil {
			return err
		}

		environments = append(environments, environment)
		return nil
	}

	if err := e.get(c, "environment", getEnvironmentf); err != nil {
		return err
	}

	return e.Printer.PrintEnvironments(environments...)
}

func (e *EnvironmentCommand) List(c *cli.Context) error {
	environmentSummaries, err := e.Client.ListEnvironments()
	if err != nil {
		return err
	}

	return e.Printer.PrintEnvironmentSummaries(environmentSummaries...)
}

func (e *EnvironmentCommand) SetMinCount(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "COUNT")
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(args["COUNT"], 10, 64)
	if err != nil {
		return NewUsageError("'%s' is not a valid integer", args["COUNT"])
	}

	id, err := e.resolveSingleID("environment", args["NAME"])
	if err != nil {
		return err
	}

	environment, err := e.Client.UpdateEnvironment(id, int(count))
	if err != nil {
		return err
	}

	return e.Printer.PrintEnvironments(environment)
}
