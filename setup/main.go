package main

import (
	"fmt"
	"log"
	"os"

	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logging"
	"github.com/quintilesims/layer0/setup/command"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

var Version string

func main() {
	if Version == "" {
		Version = "unset/developer"
	}

	app := cli.NewApp()
	app.Name = "Layer0 Setup"
	app.Usage = "Create and manage Layer0 instances"
	app.UsageText = "l0-setup [global options] command [command options] [arguments...]"
	app.Version = Version
	app.Flags = config.SetupFlags()

	commandFactory := command.NewCommandFactory(instance.NewLocalInstance, awsc.NewClient)
	app.Commands = []cli.Command{
		commandFactory.Init(),
		commandFactory.List(),
		commandFactory.Plan(),
		commandFactory.Apply(),
		commandFactory.Destroy(),
		commandFactory.Endpoint(),
		commandFactory.Push(),
		commandFactory.Pull(),
		commandFactory.Set(),
		commandFactory.Unset(),
		commandFactory.Upgrade(),
	}

	app.Before = func(c *cli.Context) error {
		logger := logging.NewLogWriter(c.Bool("debug"))
		log.SetOutput(logger)

		instance.InitializeLayer0ModuleInputs(Version)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
