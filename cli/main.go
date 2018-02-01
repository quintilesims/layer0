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

	base := &command.CommandBase{}
	app.Commands = []cli.Command{
		command.NewAdminCommand(base).Command(),
		command.NewDeployCommand(base).Command(),
		command.NewEnvironmentCommand(base).Command(),
		command.NewLoadBalancerCommand(base).Command(),
		command.NewServiceCommand(base).Command(),
		command.NewTaskCommand(base).Command(),
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

		tagResolver := resolver.NewTagResolver(apiClient)

		var p printer.Printer
		switch format := c.GlobalString(config.FLAG_OUTPUT); format {
		case "text":
			p = &printer.TextPrinter{}
		case "json":
			p = &printer.JSONPrinter{}
		default:
			return fmt.Errorf("Unrecognized output format '%s'", format)
		}

		base.SetClient(apiClient)
		base.SetResolver(tagResolver)
		base.SetPrinter(p)

		if !c.GlobalBool(config.FLAG_SKIP_VERIFY_VERSION) {
			apiConfig, err := apiClient.ReadConfig()
			if err != nil {
				return err
			}

			if apiConfig.Version != Version {
				text := fmt.Sprintf("API and CLI version mismatch (CLI: '%s', API: '%s')\n", Version, apiConfig.Version)
				text += fmt.Sprintf("To disable this warning, set %s=\"1\"", config.ENVVAR_SKIP_VERIFY_VERSION)
				return fmt.Errorf(text)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
