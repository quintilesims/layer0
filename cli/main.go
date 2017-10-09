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

	commandFactory := command.NewCommandFactory(nil, nil, nil)
	app.Commands = []cli.Command{
		commandFactory.Deploy(),
		commandFactory.Environment(),
		// todo: other entities
	}

	app.Before = func(c *cli.Context) error {
		cfg := config.NewContextCLIConfig(c)
		if err := cfg.Validate(); err != nil {
			return err
		}

		logger := logging.NewLogWriter(cfg.Debug())
		log.SetOutput(logger)

		// inject the api client
		apiClient := client.NewAPIClient(client.Config{
			Endpoint:  cfg.Endpoint(),
			Token:     cfg.Token(),
			VerifySSL: cfg.VerifySSL(),
		})

		commandFactory.SetClient(apiClient)

		// inject the resolver
		tagResolver := resolver.NewTagResolver(apiClient)
		commandFactory.SetResolver(tagResolver)

		// inject the printer
		var p printer.Printer
		switch format := cfg.Output(); format {
		case "text":
			p = &printer.TextPrinter{}
		case "json":
			p = &printer.JSONPrinter{}
		default:
			return fmt.Errorf("Unrecognized output format '%s'", format)
		}

		commandFactory.SetPrinter(p)
		
		// todo: verify version

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
