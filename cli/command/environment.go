package command

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type EnvironmentCommand struct {
	*CommandBase
}

func NewEnvironmentCommand(b *CommandBase) *EnvironmentCommand {
	return &EnvironmentCommand{b}
}

func (e *EnvironmentCommand) Command() cli.Command {
	return cli.Command{
		Name:  "environment",
		Usage: "manage layer0 environments",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new environment",
				Action:    e.create,
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
				Action:    e.delete,
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
				Action:    e.list,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "describe an environment",
				Action:    e.read,
				ArgsUsage: "NAME",
			},
			{
				Name:      "setmincount",
				Usage:     "set the minimum instance count for an environment cluster",
				Action:    e.update,
				ArgsUsage: "NAME COUNT",
			},
			{
				Name:      "link",
				Usage:     "links two environments together",
				Action:    e.link,
				ArgsUsage: "SOURCE DESTINATION",
			},
			{
				Name:      "unlink",
				Usage:     "unlinks two previously linked environments",
				Action:    e.unlink,
				ArgsUsage: "SOURCE DESTINATION",
			},
		},
	}
}

func (e *EnvironmentCommand) create(c *cli.Context) error {
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

	return e.waitOnJobHelper(c, jobID, "creating", func(environmentID string) error {
		environment, err := e.client.ReadEnvironment(environmentID)
		if err != nil {
			return err
		}

		return e.printer.PrintEnvironments(environment)
	})
}

func (e *EnvironmentCommand) delete(c *cli.Context) error {
	return e.deleteHelper(c, "environment", func(environmentID string) (string, error) {
		return e.client.DeleteEnvironment(environmentID)
	})
}

func (e *EnvironmentCommand) list(c *cli.Context) error {
	environmentSummaries, err := e.client.ListEnvironments()
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironmentSummaries(environmentSummaries...)
}

func (e *EnvironmentCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	environment, err := e.client.ReadEnvironment(args["NAME"])
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironments(environment)
}

func (e *EnvironmentCommand) update(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(args["COUNT"], 10, 64)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid integer", args["COUNT"])
	}
	minClusterCount := int(count)

	id, err := resolveSingleEntityID(e.resolver, "environment", args["NAME"])
	if err != nil {
		return err
	}

	req := models.UpdateEnvironmentRequest{
		EnvironmentID:   id,
		MinClusterCount: &minClusterCount,
	}

	jobID, err := e.client.UpdateEnvironment(req)
	if err != nil {
		return err
	}

	return e.waitOnJobHelper(c, jobID, "updating", func(environmentID string) error {
		environment, err := e.client.ReadEnvironment(environmentID)
		if err != nil {
			return err
		}

		return e.printer.PrintEnvironments(environment)
	})
}

func (e *EnvironmentCommand) link(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SOURCE", "DESTINATION")
	if err != nil {
		return err
	}

	id1, err := resolveSingleEntityID(e.resolver, "environment", args["SOURCE"])
	if err != nil {
		return err
	}

	id2, err := resolveSingleEntityID(e.resolver, "environment", args["DESTINATION"])
	if err != nil {
		return err
	}

	if id1 == id2 {
		return fmt.Errorf("Cannot link an environment to itself")
	}

	if err := e.client.CreateLink(id1, id2); err != nil {
		return err
	}

	e.printer.Printf("Environment successfully linked\n")

	return nil
}

func (e *EnvironmentCommand) unlink(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SOURCE", "DESTINATION")
	if err != nil {
		return err
	}

	id1, err := resolveSingleEntityID(e.resolver, "environment", args["SOURCE"])
	if err != nil {
		return err
	}

	id2, err := resolveSingleEntityID(e.resolver, "environment", args["DESTINATION"])
	if err != nil {
		return err
	}

	if id1 == id2 {
		return fmt.Errorf("Cannot unlink an environment from itself")
	}

	if err := e.client.DeleteLink(id1, id2); err != nil {
		return err
	}

	e.printer.Printf("Environment successfully unlinked\n")

	return nil
}
