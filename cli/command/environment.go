package command

import (
	"fmt"
	"io/ioutil"

	"github.com/quintilesims/layer0/common/config"
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
				ArgsUsage: "ENVIRONMENT_NAME",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "type",
						Value: config.DefaultEnvironmentInstanceType,
						Usage: "type of the ec2 instances to use in the environment cluster",
					},
					cli.IntFlag{
						Name:  "min-scale",
						Value: 0,
						Usage: "minimum allowed scale of the environment cluster",
					},
					cli.IntFlag{
						Name:  "max-scale",
						Value: config.DefaultEnvironmentMaxScale,
						Usage: "maximum allowed scale of the environment cluster",
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
				},
			},
			{
				Name:      "delete",
				Usage:     "delete an environment",
				ArgsUsage: "ENVIRONMENT_NAME",
				Action:    e.delete,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "r, recursive",
						Usage: "delete all dependecies (services, load balancers, and tasks)",
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
				Name:      "get",
				Usage:     "describe an environment",
				Action:    e.read,
				ArgsUsage: "ENVIRONMENT_NAME",
			},
			{
				Name:      "set-scale",
				Usage:     "update the min/max scale of an environment cluster",
				Action:    e.setScale,
				ArgsUsage: "ENVIRONMENT_NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "min-scale",
						Usage: "minimum allowed scale of the environment cluster",
					},
					cli.IntFlag{
						Name:  "max-scale",
						Usage: "maximum allowed scale of the environment cluster",
					},
				},
			},
			{
				Name:      "link",
				Usage:     "link two environments together",
				Action:    e.link,
				ArgsUsage: "SOURCE_ENVIRONMENT_NAME DESTINATION_ENVIRONMENT_NAME",
				Flags: []cli.Flag{
					cli.BoolTFlag{
						Name:  "bi-directional",
						Usage: "specifies whether the link should be bi-directional",
					},
				},
			},
			{
				Name:      "unlink",
				Usage:     "unlinks two previously linked environments",
				Action:    e.unlink,
				ArgsUsage: "SOURCE_ENVIRONMENT_NAME DESTINATION_ENVIRONMENT_NAME",
				Flags: []cli.Flag{
					cli.BoolTFlag{
						Name:  "bi-directional",
						Usage: "specifies whether the link should be direcional",
					},
				},
			},
		},
	}
}

func (e *EnvironmentCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME")
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
		EnvironmentName:  args["ENVIRONMENT_NAME"],
		InstanceType:     c.String("type"),
		MinScale:         c.Int("min-scale"),
		MaxScale:         c.Int("max-scale"),
		UserDataTemplate: userData,
		OperatingSystem:  c.String("os"),
		AMIID:            c.String("ami"),
	}

	if err := req.Validate(); err != nil {
		return err
	}

	environmentID, err := e.client.CreateEnvironment(req)
	if err != nil {
		return err
	}

	environment, err := e.client.ReadEnvironment(environmentID)
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironments(environment)
}

func (e *EnvironmentCommand) delete(c *cli.Context) error {
	if c.Bool("recursive") {
		args, err := extractArgs(c.Args(), "NAME")
		if err != nil {
			return err
		}

		environmentID, err := e.CommandBase.resolveSingleEntityIDHelper("environment", args["NAME"])
		if err != nil {
			return err
		}

		loadBalancers, err := e.client.ListLoadBalancers()
		if err != nil {
			return err
		}

		for _, loadBalancer := range loadBalancers {
			if loadBalancer.EnvironmentID == environmentID {
				if _, err := e.client.DeleteLoadBalancer(loadBalancer.LoadBalancerID); err != nil {
					return err
				}
			}
		}

		tasks, err := e.client.ListTasks()
		if err != nil {
			return err
		}

		for _, task := range tasks {
			if task.EnvironmentID == environmentID {
				if _, err := e.client.DeleteTask(task.TaskID); err != nil {
					return err
				}
			}
		}

		services, err := e.client.ListServices()
		if err != nil {
			return err
		}

		for _, service := range services {
			if service.EnvironmentID == environmentID {
				if _, err := e.client.DeleteService(service.ServiceID); err != nil {
					return err
				}
			}
		}
	}

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
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME")
	if err != nil {
		return err
	}

	environmentIDs, err := e.resolver.Resolve("environment", args["ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	environments := make([]*models.Environment, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		environment, err := e.client.ReadEnvironment(environmentID)
		if err != nil {
			return err
		}

		environments[i] = environment
	}

	return e.printer.PrintEnvironments(environments...)
}

func (e *EnvironmentCommand) setScale(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME")
	if err != nil {
		return err
	}

	req := models.UpdateEnvironmentRequest{}
	if c.IsSet("min-scale") {
		minScale := c.Int("min-scale")
		req.MinScale = &minScale
	}

	if c.IsSet("max-scale") {
		maxScale := c.Int("max-scale")
		req.MaxScale = &maxScale
	}

	if err := req.Validate(); err != nil {
		return err
	}

	environmentID, err := e.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	if err := e.client.UpdateEnvironment(environmentID, req); err != nil {
		return err
	}

	environment, err := e.client.ReadEnvironment(environmentID)
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironments(environment)
}

func (e *EnvironmentCommand) link(c *cli.Context) error {
	fn := func(src, dst *models.Environment) models.UpdateEnvironmentRequest {
		src.Links = append(src.Links, dst.EnvironmentID)
		return models.UpdateEnvironmentRequest{Links: &src.Links}
	}

	return e.updateLinksHelper(c, fn)
}

func (e *EnvironmentCommand) unlink(c *cli.Context) error {
	fn := func(src, dst *models.Environment) models.UpdateEnvironmentRequest {
		for i, environmentID := range src.Links {
			if environmentID == dst.EnvironmentID {
				src.Links = append(src.Links[:i], src.Links[i+1:]...)
				i--
			}
		}

		return models.UpdateEnvironmentRequest{Links: &src.Links}
	}

	return e.updateLinksHelper(c, fn)
}

func (e *EnvironmentCommand) updateLinksHelper(
	c *cli.Context,
	createReqFN func(src, dst *models.Environment) models.UpdateEnvironmentRequest,
) error {
	args, err := extractArgs(c.Args(), "SOURCE_ENVIRONMENT_NAME", "DESTINATION_ENVIRONMENT_NAME")
	if err != nil {
		return err
	}

	srcEnvironmentID, err := e.resolveSingleEntityIDHelper("environment", args["SOURCE_ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	dstEnvironmentID, err := e.resolveSingleEntityIDHelper("environment", args["DESTINATION_ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	if srcEnvironmentID == dstEnvironmentID {
		return fmt.Errorf("Cannot link/unlink an environment to/from itself")
	}

	srcEnvironment, err := e.client.ReadEnvironment(srcEnvironmentID)
	if err != nil {
		return err
	}

	dstEnvironment, err := e.client.ReadEnvironment(dstEnvironmentID)
	if err != nil {
		return err
	}

	updateLinkFN := func(src, dst *models.Environment) error {
		req := createReqFN(src, dst)
		if err := e.client.UpdateEnvironment(src.EnvironmentID, req); err != nil {
			return err
		}

		return nil
	}

	if err := updateLinkFN(srcEnvironment, dstEnvironment); err != nil {
		return err
	}

	if !c.Bool("bi-directional") {
		return nil
	}

	return updateLinkFN(dstEnvironment, srcEnvironment)
}
