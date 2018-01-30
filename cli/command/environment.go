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
						Usage: "cluster size. setting scale explicitly, results in a static environment",
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

	if err := req.Validate(); err != nil {
		return err
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

	req := models.UpdateEnvironmentRequest{}
	scale, err := strconv.Atoi(args["SCALE"])
	if err == nil {
		req.Scale = &scale
	}

	if err := req.Validate(); err != nil {
		return err
	}

	environmentID, err := e.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	jobID, err := e.client.UpdateEnvironment(environmentID, req)
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
	generateAddLinkRequest := func(sourceEnvironmentID, destEnvironmentID string) (models.UpdateEnvironmentRequest, error) {
		env, err := e.client.ReadEnvironment(sourceEnvironmentID)
		if err != nil {
			return models.UpdateEnvironmentRequest{}, err
		}

		appendedLinks := []string{}
		envLinks := models.LinkTags(env.Links)
		if !envLinks.Contains(destEnvironmentID) {
			appendedLinks = append(env.Links, destEnvironmentID)
		}

		req := models.UpdateEnvironmentRequest{Links: &appendedLinks}
		return req, nil
	}

	return e.updateEnvironmentLinksHelper(c, generateAddLinkRequest)
}

func (e *EnvironmentCommand) unlink(c *cli.Context) error {
	generateRemoveLinkRequest := func(sourceEnvironmentID, destEnvironmentID string) (models.UpdateEnvironmentRequest, error) {
		env, err := e.client.ReadEnvironment(sourceEnvironmentID)
		if err != nil {
			return models.UpdateEnvironmentRequest{}, err
		}

		updatedLinks := []string{}
		for _, link := range env.Links {
			if link != destEnvironmentID {
				updatedLinks = append(updatedLinks, link)
			}
		}

		req := models.UpdateEnvironmentRequest{Links: &updatedLinks}

		return req, nil
	}

	return e.updateEnvironmentLinksHelper(c, generateRemoveLinkRequest)
}

func (e *EnvironmentCommand) updateEnvironmentLinksHelper(
	c *cli.Context,
	generateReq func(string, string) (models.UpdateEnvironmentRequest, error),
) error {
	args, err := extractArgs(c.Args(), "SOURCE_ENVIRONMENT_NAME", "DESTINATION_ENVIRONMENT_NAME")
	if err != nil {
		return err
	}

	sourceEnvironmentID, err := e.resolveSingleEntityIDHelper("environment", args["SOURCE_ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	destEnvironmentID, err := e.resolveSingleEntityIDHelper("environment", args["DESTINATION_ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	if sourceEnvironmentID == destEnvironmentID {
		return fmt.Errorf("Cannot unlink an environment from itself")
	}

	updateLinkFN := func(sourceEnvID, destEnvID string) error {
		updateEnvReq, err := generateReq(sourceEnvID, destEnvID)
		if err != nil {
			return err
		}

		jobID, err := e.client.UpdateEnvironment(sourceEnvID, updateEnvReq)
		if err != nil {
			return err
		}

		return e.waitOnJobHelper(c, jobID, "updating", func(environmentID string) error {
			e.printer.Printf("Environment update successfull")
			return nil
		})
	}

	updateLinkFN(sourceEnvironmentID, destEnvironmentID)

	if !c.Bool("bi-directional") {
		return nil
	}

	return updateLinkFN(destEnvironmentID, sourceEnvironmentID)
}
