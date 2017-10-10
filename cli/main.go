package main

import (
	"fmt"
	"log"
	"os"

	"github.com/quintilesims/layer0/cli/command"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logging"
	"github.com/urfave/cli"
)

var Version string

func main() {
	if Version == "" {
		Version = "unset/developer"
	}

	app := cli.NewApp()
	app.Name = "l0"
	app.Usage = "Manage Layer0"
	app.UsageText = "l0 [global options] command [command options] [arguments...]"
	app.Version = Version
	app.Flags = config.CLIFlags()

<<<<<<< 058026bb31c76f9af4c89a26f855677f76370c25
	base := &command.CommandBase{}
	app.Commands = []cli.Command{
		command.NewAdminCommand(base).Command(),
		command.NewDeployCommand(base).Command(),
		command.NewEnvironmentCommand(base).Command(),
		command.NewJobCommand(base).Command(),
		command.NewLoadBalancerCommand(base).Command(),
		command.NewServiceCommand(base).Command(),
		command.NewTaskCommand(base).Command(),
=======
	mediator := &command.CommandMediator{}
	app.Commands = []cli.Command{
		command.NewDeploy(mediator),
		command.NewEnvironment(mediator),
		// todo: other entities
>>>>>>> CommandFactory to CommandMediator
	}

	app.Before = func(c *cli.Context) error {
		if err := config.ValidateCLIContext(c); err != nil {
			return err
		}

		logger := logging.NewLogWriter(c.GlobalBool(config.FLAG_DEBUG))
		log.SetOutput(logger)

		apiClient := client.NewAPIClient(client.Config{
			Endpoint:  c.GlobalString(config.FLAG_ENDPOINT),
			Token:     c.GlobalString(config.FLAG_TOKEN),
			VerifySSL: !c.GlobalBool(config.FLAG_SKIP_VERIFY_SSL),
		})

<<<<<<< 058026bb31c76f9af4c89a26f855677f76370c25
		tagResolver := resolver.NewTagResolver(apiClient)
=======
		mediator.SetClient(apiClient)

		// inject the resolver
		tagResolver := resolver.NewTagResolver(apiClient)
		mediator.SetResolver(tagResolver)
>>>>>>> CommandFactory to CommandMediator

		var p printer.Printer
		switch format := c.GlobalString(config.FLAG_OUTPUT); format {
		case "text":
			p = &printer.TextPrinter{}
		case "json":
			p = &printer.JSONPrinter{}
		default:
			return fmt.Errorf("Unrecognized output format '%s'", format)
		}

<<<<<<< 058026bb31c76f9af4c89a26f855677f76370c25
		base.SetClient(apiClient)
		base.SetResolver(tagResolver)
		base.SetPrinter(p)
=======
		mediator.SetPrinter(p)
>>>>>>> CommandFactory to CommandMediator

		// todo: verify version

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
