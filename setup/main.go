package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/command"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "Layer0 Setup"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "l, log",
			Value:  "info",
			EnvVar: config.SETUP_LOG_LEVEL,
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

	app.Before = func(c *cli.Context) error {
		switch level := strings.ToLower(c.String("log")); level {
		case "0", "debug":
			logrus.SetLevel(logrus.DebugLevel)
		case "1", "info":
			logrus.SetLevel(logrus.InfoLevel)
		case "2", "warning":
			logrus.SetLevel(logrus.WarnLevel)
		case "3", "error":
			logrus.SetLevel(logrus.ErrorLevel)
		case "4", "fatal":
			logrus.SetLevel(logrus.FatalLevel)
		default:
			return fmt.Errorf("Unrecognized log level '%s'", level)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
