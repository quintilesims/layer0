package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/cli/command"
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/waitutils"
	"github.com/urfave/cli"
	"os"
	"sync"
	"time"
)

var Version string

func main() {
	if Version == "" {
		Version = "unset/developer"
	}

	config.SetCLIVersion(Version)
	RunApp()
}

func RunApp() {
	app := cli.NewApp()
	app.Name = "l0"
	app.Usage = "Manage Layer0"
	app.UsageText = "l0 [global options] command [command options] [arguments...]"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "o, output",
			Value: "text",
			Usage: "output format [text,json]",
		},
		cli.StringFlag{
			Name:  "t, timeout",
			Value: "15m",
			Usage: "timeout [h,m,s,ms]",
		},
		cli.BoolFlag{
			Name:  "d, debug",
			Usage: "Print debug statements",
		},
	}

	apiClient := client.NewAPIClient(client.Config{
		Endpoint:      config.APIEndpoint(),
		Token:         fmt.Sprintf("Basic %s", config.AuthToken()),
		VerifySSL:     config.ShouldVerifySSL(),
		VerifyVersion: config.ShouldVerifyVersion(),
		Clock:         waitutils.RealClock{},
	})

	commands := getCommands(apiClient)
	for _, cmd := range commands {
		app.Commands = append(app.Commands, cmd.GetCommand())
	}

	app.Commands = append(app.Commands, command.HelpCommand(app))

	var timeout time.Duration
	var wg sync.WaitGroup
	wg.Add(1)

	app.Before = func(c *cli.Context) error {
		defer wg.Done()

		var p printer.Printer
		format := c.GlobalString("output")
		switch format {
		case "text":
			p = &printer.TextPrinter{}
		case "json":
			p = &printer.JSONPrinter{}
		default:
			return fmt.Errorf("Unrecognized output format '%s'", format)
		}

		for _, cmd := range commands {
			cmd.SetPrinter(p)
		}

		t, err := time.ParseDuration(c.GlobalString("timeout"))
		if err != nil {
			return err
		}

		timeout = t

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}

	go func() {
		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	// wait for timeout to be set in app.Before
	// if the user inputs invalid syntax, app.Before won't run
	wg.Wait()
	<-time.After(timeout)
	log.Fatalf("Timeout after %v", timeout)
}

func getCommands(client *client.APIClient) []command.CommandGroup {
	cmd := &command.Command{
		Client:   client,
		Resolver: command.NewTagResolver(client),
		Printer:  nil,
	}

	return []command.CommandGroup{
		command.NewAdminCommand(cmd),
		command.NewDeployCommand(cmd),
		command.NewEnvironmentCommand(cmd),
		command.NewJobCommand(cmd),
		command.NewLoadBalancerCommand(cmd),
		command.NewServiceCommand(cmd),
		command.NewTaskCommand(cmd),
	}
}
