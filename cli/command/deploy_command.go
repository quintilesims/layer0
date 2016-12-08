package command

import (
	"github.com/quintilesims/layer0/cli/entity"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"io/ioutil"
)

type DeployCommand struct {
	*Command
}

func NewDeployCommand(command *Command) *DeployCommand {
	return &DeployCommand{command}
}

func (d *DeployCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "manage layer0 deploys",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new deploy",
				Action:    wrapAction(d.Command, d.Create),
				ArgsUsage: "PATH NAME",
			},
			{
				Name:      "delete",
				Usage:     "delete a deploy",
				ArgsUsage: "NAME",
				Action:    wrapAction(d.Command, d.Delete),
			},
			{
				Name:      "get",
				Usage:     "describe a deploy",
				Action:    wrapAction(d.Command, d.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all deploys (only the latest versions of each family will be shown)",
				Action:    wrapAction(d.Command, d.List),
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

func (d *DeployCommand) Create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "PATH", "NAME")
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(args["PATH"])
	if err != nil {
		return err
	}

	deploy, err := d.Client.CreateDeploy(args["NAME"], content)
	if err != nil {
		return err
	}

	return d.printDeploy(deploy)
}

func (d *DeployCommand) Delete(c *cli.Context) error {
	return d.delete(c, "deploy", d.Client.DeleteDeploy)
}

func (d *DeployCommand) Get(c *cli.Context) error {
	return d.get(c, "deploy", func(id string) (entity.Entity, error) {
		deploy, err := d.Client.GetDeploy(id)
		if err != nil {
			return nil, err
		}

		return entity.NewDeploy(deploy), nil
	})
}

func (d *DeployCommand) List(c *cli.Context) error {
	deploys, err := d.Client.ListDeploys()
	if err != nil {
		return err
	}

	if !c.Bool("all") {
		deploys = filterDeploys(deploys)
	}

	return d.printDeploys(deploys)
}

func (d *DeployCommand) printDeploy(deploy *models.Deploy) error {
	entity := entity.NewDeploy(deploy)
	return d.Printer.PrintEntity(entity)
}

func (d *DeployCommand) printDeploys(deploys []*models.Deploy) error {
	entities := []entity.Entity{}
	for _, deploy := range deploys {
		entities = append(entities, entity.NewDeploy(deploy))
	}

	return d.Printer.PrintEntities(entities)
}

func filterDeploys(deploys []*models.Deploy) []*models.Deploy {
	catalog := map[string]*models.Deploy{}

	for _, deploy := range deploys {
		if name := deploy.DeployName; name != "" {
			if _, exists := catalog[name]; !exists {
				catalog[name] = deploy
				continue
			}

			if deploy.Version > catalog[name].Version {
				catalog[name] = deploy
			}
		}
	}

	filtered := []*models.Deploy{}
	for _, deploy := range catalog {
		filtered = append(filtered, deploy)
	}

	return filtered
}
