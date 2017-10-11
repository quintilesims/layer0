package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"io/ioutil"
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
				ArgsUsage: "PATH NAME",
			},
			{
				Name:      "delete",
				Usage:     "delete a deploy",
				ArgsUsage: "NAME",
				Action:    d.delete,
			},
			{
				Name:      "get",
				Usage:     "describe a deploy",
				Action:    d.get,
				ArgsUsage: "NAME",
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
	args, err := extractArgs(c.Args(), "PATH", "NAME")
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(args["PATH"])
	if err != nil {
		return err
	}

	req := models.CreateDeployRequest{
		DeployName: args["NAME"],
		DeployFile: content,
	}

	jobID, err := d.client.CreateDeploy(req)
	if err != nil {
		return err
	}

	// return d.Printer.PrintDeploys(deploy)
	return d.waitOnJobHelper(c, jobID, "creating", func(deployID string) error {
		environment, err := d.client.ReadDeploy(deployID)
		if err != nil {
			return err
		}

		return d.printer.PrintEnvironments(environment)
	})
}

func (d *DeployCommand) delete(c *cli.Context) error {
	return d.delete(c, "deploy", d.Client.DeleteDeploy)
}

func (d *DeployCommand) get(c *cli.Context) error {
	deploys := []*models.Deploy{}
	getDeployf := func(id string) error {
		deploy, err := d.client.GetDeploy(id)
		if err != nil {
			return err
		}

		deploys = append(deploys, deploy)
		return nil
	}

	if err := d.get(c, "deploy", getDeployf); err != nil {
		return err
	}

	return d.Printer.PrintDeploys(deploys...)
}

func (d *DeployCommand) list(c *cli.Context) error {
	deploySummaries, err := d.Client.ListDeploys()
	if err != nil {
		return err
	}

	if !c.Bool("all") {
		deploySummaries, err = filterDeploySummaries(deploySummaries)
		if err != nil {
			return err
		}
	}

	return d.Printer.PrintDeploySummaries(deploySummaries...)
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
				return nil, err
			}

			current, err := strconv.Atoi(deploy.Version)
			if err != nil {
				return nil, err
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
