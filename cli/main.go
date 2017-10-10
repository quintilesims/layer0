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

	mediator := &command.CommandMediator{}
	app.Commands = []cli.Command{
		command.NewDeployCommand(mediator).Command(),
		command.NewEnvironmentCommand(mediator).Command(),
		command.NewLoadBalancerCommand(mediator).Command(),
		// todo: other entities
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

		mediator.SetClient(apiClient)

		// inject the resolver
		tagResolver := resolver.NewTagResolver(apiClient)
		mediator.SetResolver(tagResolver)

		// inject the printer
		var p printer.Printer
		switch format := c.GlobalString(config.FLAG_OUTPUT); format {
		case "text":
			p = &printer.TextPrinter{}
		case "json":
			p = &printer.JSONPrinter{}
		default:
			return fmt.Errorf("Unrecognized output format '%s'", format)
		}

		mediator.SetPrinter(p)

		// todo: verify version

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
