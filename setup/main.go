package main

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/command"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Layer0 Setup"

	instanceFactory := instance.NewLayer0Factory()
	commandFactory := command.NewCommandFactory(instanceFactory)
	app.Commands = []cli.Command{
		commandFactory.Init(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
