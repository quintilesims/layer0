package command

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type DeployCommand struct {
	*CommandBase
}

func NewDeployCommand(b *CommandBase) *DeployCommand {
	return &DeployCommand{
		CommandBase: b,
	}
}

func (d *DeployCommand) Command() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "manage layer0 deploys",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new deploy",
				Action:    d.create,
				ArgsUsage: "PATH DEPLOY_NAME",
			},
			{
				Name:      "delete",
				Usage:     "delete a deploy",
				ArgsUsage: "DEPLOY_NAME",
				Action:    d.delete,
			},
			{
				Name:      "get",
				Usage:     "describe a deploy",
				Action:    d.read,
				ArgsUsage: "DEPLOY_NAME",
			},
			{
				Name:      "list",
				Usage:     "list all deploys (only the latest versions of each family will be shown)",
				Action:    d.list,
				ArgsUsage: " ",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "all",
						Usage: "list all versions of all deploys",
					},
				},
			},
		},
	}
}

func (d *DeployCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "PATH", "DEPLOY_NAME")
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(args["PATH"])
	if err != nil {
		return err
	}

	req := models.CreateDeployRequest{
		DeployName: args["DEPLOY_NAME"],
		DeployFile: content,
	}

	jobID, err := d.client.CreateDeploy(req)
	if err != nil {
		return err
	}

	return d.waitOnJobHelper(c, jobID, "creating", func(deployID string) error {
		deploy, err := d.client.ReadDeploy(deployID)
		if err != nil {
			return err
		}

		return d.printer.PrintDeploys(deploy)
	})
}

func (d *DeployCommand) delete(c *cli.Context) error {
	return d.deleteHelper(c, "deploy", func(deployID string) (string, error) {
		return d.client.DeleteDeploy(deployID)
	})
}

func (d *DeployCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "DEPLOY_NAME")
	if err != nil {
		return err
	}

	deployIDs, err := d.resolver.Resolve("deploy", args["DEPLOY_NAME"])
	if err != nil {
		return err
	}

	deploys := make([]*models.Deploy, len(deployIDs))
	for i, deployID := range deployIDs {
		deploy, err := d.client.ReadDeploy(deployID)
		if err != nil {
			return err
		}

		deploys[i] = deploy
	}

	return d.printer.PrintDeploys(deploys...)
}

func (d *DeployCommand) list(c *cli.Context) error {
	deploySummaries, err := d.client.ListDeploys()
	if err != nil {
		return err
	}

	if !c.Bool("all") {
		deploySummaries, err = filterDeploySummaries(deploySummaries)
		if err != nil {
			return err
		}
	}

	return d.printer.PrintDeploySummaries(deploySummaries...)
}

func filterDeploySummaries(deploys []*models.DeploySummary) ([]*models.DeploySummary, error) {
	catalog := map[string]*models.DeploySummary{}

	for _, deploy := range deploys {
		if name := deploy.DeployName; name != "" {
			if _, exists := catalog[name]; !exists {
				catalog[name] = deploy
				continue
			}

			max, err := strconv.Atoi(catalog[name].Version)
			if err != nil {
				return nil, fmt.Errorf("Deploy %s has an invalid version tag: %v", catalog[name].DeployID, err)
			}

			current, err := strconv.Atoi(deploy.Version)
			if err != nil {
				return nil, fmt.Errorf("Deploy %s has an invalid version tag: %v", deploy.Version, err)
			}

			if current > max {
				catalog[name] = deploy
			}
		}
	}

	filtered := []*models.DeploySummary{}
	for _, deploy := range catalog {
		filtered = append(filtered, deploy)
	}

	return filtered, nil
}
