package main

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/command"
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

	commandFactory := command.NewCommandFactory()
	app.Commands = []cli.Command{
		commandFactory.Init(),
		commandFactory.List(),
		commandFactory.Plan(),
		commandFactory.Apply(),
		commandFactory.Destroy(),
		commandFactory.Endpoint(),
		commandFactory.Push(),
		commandFactory.Pull(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
