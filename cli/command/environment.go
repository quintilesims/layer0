package command

import (
	"fmt"
	"io/ioutil"
	"strconv"

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
						Name:  "scale",
						Value: 0,
						Usage: "specifies the number of instances in the cluster. Setting this value will create a static environment",
					},
					cli.StringFlag{
						Name:  "user-data",
						Usage: "path to user data file",
					},
					cli.StringFlag{
						Name:  "os",
						Value: config.DefaultEnvironmentOS,
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
				Name:      "logs",
				Usage:     "get the logs for an environment",
				Action:    e.logs,
				ArgsUsage: "ENVIRONMENT_NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "tail",
						Usage: "number of lines from the end to return (default: 0)",
					},
					cli.StringFlag{
						Name:  "start",
						Usage: "the start of the time range to fetch logs (format: YYYY-MM-DD HH:MM)",
					},
					cli.StringFlag{
						Name:  "end",
						Usage: "the end of the time range to fetch logs (format: YYYY-MM-DD HH:MM)",
					},
				},
			},
			{
				Name:      "get",
				Usage:     "describe an environment",
				Action:    e.read,
				ArgsUsage: "ENVIRONMENT_NAME",
			},
			{
				Name:      "scale",
				Usage:     "update the scale of a static environment cluster",
				Action:    e.setScale,
				ArgsUsage: "ENVIRONMENT_NAME SCALE",
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
		Scale:            c.Int("scale"),
		EnvironmentType:  config.DefaultEnvironmentType,
		UserDataTemplate: userData,
		OperatingSystem:  c.String("os"),
		AMIID:            c.String("ami"),
	}

	if c.IsSet("scale") {
		req.EnvironmentType = models.EnvironmentTypeStatic
	}

	if !c.IsSet("os") {
		req.OperatingSystem = config.DefaultEnvironmentOS
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
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	environmentID, err := e.CommandBase.resolveSingleEntityIDHelper("environment", args["NAME"])
	if err != nil {
		return err
	}

	if c.Bool("recursive") {
		loadBalancers, err := e.client.ListLoadBalancers()
		if err != nil {
			return err
		}

		for _, loadBalancer := range loadBalancers {
			if loadBalancer.EnvironmentID == environmentID {
				if err := e.client.DeleteLoadBalancer(loadBalancer.LoadBalancerID); err != nil {
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
				if err := e.client.DeleteTask(task.TaskID); err != nil {
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
				if err := e.client.DeleteService(service.ServiceID); err != nil {
					return err
				}
			}
		}
	}

	return e.client.DeleteEnvironment(environmentID)
}

func (e *EnvironmentCommand) list(c *cli.Context) error {
	environmentSummaries, err := e.client.ListEnvironments()
	if err != nil {
		return err
	}

	return e.printer.PrintEnvironmentSummaries(environmentSummaries...)
}

func (e *EnvironmentCommand) logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME")
	if err != nil {
		return err
	}

	environmentID, err := e.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	query := buildLogQueryHelper(c.String("start"), c.String("end"), c.Int("tail"))

	logs, err := e.client.ReadEnvironmentLogs(environmentID, query)
	if err != nil {
		return err
	}

	return e.printer.PrintLogs(logs...)
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
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME", "SCALE")
	if err != nil {
		return err
	}

	scale, err := strconv.Atoi(args["SCALE"])
	if err != nil {
		return err
	}

	req := models.UpdateEnvironmentRequest{
		Scale: &scale,
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
