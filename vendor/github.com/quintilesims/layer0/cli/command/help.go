package command

import (
	"log"

	"github.com/urfave/cli"
)

func HelpCommand(app *cli.App) cli.Command {
	return cli.Command{
		Name:      "help",
		ShortName: "h",
		Usage:     "show help for a command or subcommand",
		ArgsUsage: "COMMAND [SUBCOMMAND]",
		Action: func(c *cli.Context) {
			args := c.Args()
			if len(args) == 0 {
				cli.ShowAppHelp(c)
				return
			}

			a := &cli.App{
				Name:                 app.Name,
				HelpName:             app.HelpName,
				Usage:                app.Usage,
				UsageText:            app.UsageText,
				ArgsUsage:            app.ArgsUsage,
				Version:              app.Version,
				Commands:             app.Commands,
				Flags:                app.Flags,
				EnableBashCompletion: app.EnableBashCompletion,
				HideHelp:             app.HideHelp,
				HideVersion:          app.HideVersion,
				BashComplete:         app.BashComplete,
				Action:               app.Action,
				CommandNotFound:      app.CommandNotFound,
				OnUsageError:         app.OnUsageError,
				Compiled:             app.Compiled,
				Authors:              app.Authors,
				Copyright:            app.Copyright,
				Author:               app.Author,
				Email:                app.Email,
				Writer:               app.Writer,
				ErrWriter:            app.ErrWriter,
				Metadata:             app.Metadata,
			}

			args = append([]string{"l0"}, args...)
			args = append(args, "-h")

			if err := a.Run(args); err != nil {
				log.Fatal(err)
			}
		},
	}
}
