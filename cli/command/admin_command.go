package command

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type AdminCommand struct {
	*Command
}

func NewAdminCommand(command *Command) *AdminCommand {
	return &AdminCommand{command}
}

func (a *AdminCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:        "admin",
		Usage:       "manage the layer0 api",
		Description: "manage the Layer0 API",
		Subcommands: []cli.Command{
			{
				Name:      "debug",
				Usage:     "generate debug information",
				Action:    wrapAction(a.Command, a.Debug),
				ArgsUsage: " ",
			},
			{
				Name:      "sql",
				Usage:     "initialize sql settings on the layer0 api",
				Action:    wrapAction(a.Command, a.SQL),
				ArgsUsage: " ",
			},
			{
				Name:      "version",
				Usage:     "show the current version of the layer0 api",
				Action:    wrapAction(a.Command, a.Version),
				ArgsUsage: " ",
			},
			{
				Name:      "rightsizer",
				Usage:     "Run the right sizer on the layer0 api",
				Action:    wrapAction(a.Command, a.RightSizer),
				ArgsUsage: " ",
			},
		},
	}
}

func (a *AdminCommand) Debug(c *cli.Context) error {
	apiEndpoint := config.APIEndpoint()
	cliAuth := config.AuthToken()
	cliVersion := config.CLIVersion()

	apiVersion, err := a.Client.GetVersion()
	if err != nil {
		return err
	}

	sslVerify := "disabled"
	if config.ShouldVerifySSL() {
		sslVerify = "enabled"
	}

	versionVerify := "disabled"
	if config.ShouldVerifyVersion() {
		versionVerify = "enabled"
	}

	a.Printer.Printf("DEBUG REPORT \n")
	a.Printer.Printf("------------ \n")
	a.Printer.Printf("API Endpoint:   %v\n", apiEndpoint)
	a.Printer.Printf("API Version:    %v\n", apiVersion)
	a.Printer.Printf("CLI Version:    %v\n", cliVersion)
	a.Printer.Printf("CLI Auth:       %v\n", cliAuth)
	a.Printer.Printf("SSL Verify::    %v\n", sslVerify)
	a.Printer.Printf("Version Verify: %v\n", versionVerify)

	return nil
}

func (a *AdminCommand) SQL(c *cli.Context) error {
	if err := a.Client.UpdateSQL(); err != nil {
		return err
	}

	a.Printer.Printf("Successfully Configured SQL\n")
	return nil

}

func (a *AdminCommand) Version(c *cli.Context) error {
	version, err := a.Client.GetVersion()
	if err != nil {
		return err
	}

	a.Printer.Printf("%s\n", version)
	return nil
}

func (a *AdminCommand) RightSizer(c *cli.Context) error {
	if err := a.Client.RunRightSizer(); err != nil {
		return err
	}

	a.Printer.Printf("Right Sizer finished successfully\n")
	return nil
}
