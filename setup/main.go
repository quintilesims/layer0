package main

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/command"
	"github.com/quintilesims/layer0/setup/layer0"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Layer0 Setup"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "profile",
			Usage:  "profile to use for configuration",
			EnvVar: "LAYER0_PROFILE",
		},
	}

	context := layer0.NewLocalContext()
	commandFactory := command.NewCommandFactory(context)
	app.Commands = []cli.Command{
		commandFactory.Init(),
		commandFactory.Config(),
		commandFactory.List(),
		commandFactory.Plan(),
		commandFactory.Apply(),
		commandFactory.Destroy(),
		commandFactory.Endpoint(),
		commandFactory.Push(),
		commandFactory.Pull(),
		commandFactory.Upgrade(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
