package command

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type AdminCommand struct {
	*CommandBase
}

func NewAdminCommand(b *CommandBase) *AdminCommand {
	return &AdminCommand{b}
}

func (a *AdminCommand) Command() cli.Command {
	return cli.Command{
		Name:        "admin",
		Usage:       "manage the layer0 api",
		Description: "manage the Layer0 API",
		Subcommands: []cli.Command{
			{
				Name:      "debug",
				Usage:     "generate debug information",
				Action:    a.debug,
				ArgsUsage: " ",
			},
		},
	}
}

func (a *AdminCommand) debug(c *cli.Context) error {
	apiEndpoint := c.GlobalString(config.FLAG_ENDPOINT)
	cliAuth := c.GlobalString(config.FLAG_TOKEN)
	cliVersion := c.App.Version

	debugInfo, err := a.client.ReadConfig()
	if err != nil {
		return err
	}

	sslVerify := "enabled"
	if c.GlobalBool(config.FLAG_SKIP_VERIFY_SSL) {
		sslVerify = "disabled"
	}

	versionVerify := "enabled"
	if c.GlobalBool(config.FLAG_SKIP_VERIFY_VERSION) {
		versionVerify = "disabled"
	}

	a.printer.Printf("DEBUG REPORT \n")
	a.printer.Printf("------------ \n")
	a.printer.Printf("Instance Name:    %v\n", debugInfo.Instance)
	a.printer.Printf("API Endpoint:     %v\n", apiEndpoint)
	a.printer.Printf("API Version:      %v\n", debugInfo.Version)
	a.printer.Printf("CLI Version:      %v\n", cliVersion)
	a.printer.Printf("SSL Verify:       %v\n", sslVerify)
	a.printer.Printf("Version Verify:   %v\n", versionVerify)
	a.printer.Printf("------------ \n")
	a.printer.Printf("VPC ID:           %v\n", debugInfo.VPCID)
	a.printer.Printf("Public Subnets:   %v\n", debugInfo.PublicSubnets)
	a.printer.Printf("Private Subnets:  %v\n", debugInfo.PrivateSubnets)

	return nil
}
