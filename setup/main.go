package main

import (
	"fmt"
	"os"

	"github.com/jawher/mow.cli"
	"github.com/quintilesims/layer0/setup/context"
)

var Version string

var InstanceArg = cli.StringArg{
	Name:  "INSTANCE",
	Value: "",
	Desc:  "Layer0 instance name",
}

func exit(c *context.Context, err error) {
	if err != nil {
		c.Save()
		fmt.Printf("[ERROR] %v\n ", err)
		os.Exit(1)
	}

	if err := c.Save(); err != nil {
		fmt.Printf("[ERROR] %v\n ", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func loadFlags(cmd *cli.Cmd, variableFlags []string) map[string]*string {
	flags := map[string]*string{}
	for _, name := range variableFlags {
		if variable, ok := context.GetTerraformVariable(name); ok {
			// default will be inserted during c.Load
			flags[name] = cmd.StringOpt(name, "", variable.Description)
		} else {
			err := fmt.Errorf("Terraform variable '%s' does not exist", name)
			exit(nil, err)
		}
	}

	return flags
}

type Command func(*context.Context) error

func basicCommand(cmd *cli.Cmd, command Command) {
	flagsCommand(cmd, command, nil)
}

func flagsCommand(cmd *cli.Cmd, command Command, flags map[string]*string) {
	instance := cmd.String(InstanceArg)

	cmd.Action = func() {
		c, err := context.NewContext(*instance, Version, flags)
		if err != nil {
			exit(nil, err)
		}

		exit(c, command(c))
	}
}

func main() {
	if Version == "" {
		Version = "0.0.1x-unset-develop"
	}

	app := cli.App("l0-setup", "Create and manage Layer0 instances")
	app.Version("v version", Version)

	app.Command("apply", "Create/Update a Layer0", func(cmd *cli.Cmd) {
		flags := loadFlags(cmd, []string{"access_key", "secret_key", "region"})

		dockercfg := cmd.StringOpt("dockercfg", "", "Path to valid config.json or dockercfg file")

		force := cmd.BoolOpt("force", false, "Set this flag to skip prompting on a missing dockercfg file")

		vpc := cmd.StringOpt("vpc", "", "VPC id to target.  Will create new VPC if blank.")
		flags["vpc_id"] = vpc

		command := func(c *context.Context) error {
			return context.Apply(c, *force, *dockercfg)
		}

		flagsCommand(cmd, command, flags)
	})

	app.Command("plan", "Plan an update for a Layer0", func(cmd *cli.Cmd) {
		args := cmd.StringsArg("ARGS", nil, "Terraform arguments")
		cmd.Spec = "INSTANCE [-- ARGS...]"

		command := func(c *context.Context) error {
			return context.Plan(c, *args)
		}

		basicCommand(cmd, command)
	})

	app.Command("backup", "Backup Layer0 resource files to S3", func(cmd *cli.Cmd) {
		basicCommand(cmd, context.Backup)
	})

	app.Command("destroy", "Destroy a Layer0", func(cmd *cli.Cmd) {
		force := cmd.BoolOpt("force", false, "Set this flag to skip prompting on destroy")

		command := func(c *context.Context) error {
			return context.Destroy(c, *force)
		}

		basicCommand(cmd, command)
	})

	app.Command("migrate", "Migrate old state files to the current version", func(cmd *cli.Cmd) {
		basicCommand(cmd, context.Migrate)
	})

	app.Command("restore", "Restore Layer0 resource files from S3", func(cmd *cli.Cmd) {
		flags := loadFlags(cmd, []string{"access_key", "secret_key", "region"})
		flagsCommand(cmd, context.Restore, flags)
	})

	app.Command("endpoint", "Configure the endpoint for Layer0 CLI", func(cmd *cli.Cmd) {
		syntax := cmd.StringOpt("s syntax", "bash", "Show commands using the specified syntax (bash, powershell, cmd)")
		insecure := cmd.BoolOpt("i insecure", false, "Allow incomplete SSL configuration. NOT RECOMMENDED FOR PRODUCTION USE!")
		dev := cmd.BoolOpt("d dev", false, "Show configuration variables required for local development")
		quiet := cmd.BoolOpt("q quiet", false, "Silence CLI and API version mismatch warning messages")

		command := func(c *context.Context) error {
			return context.Endpoint(c, *syntax, *insecure, *dev, *quiet)
		}

		basicCommand(cmd, command)
	})

	app.Command("terraform", "Send a command directly to terraform using Layer0 resource files", func(cmd *cli.Cmd) {
		args := cmd.StringsArg("ARGS", nil, "Terraform arguments")
		cmd.Spec = "INSTANCE [-- ARGS...]"

		command := func(c *context.Context) error {
			return context.Terraform(c, *args)
		}

		basicCommand(cmd, command)
	})

	app.Command("vpc", "Lookup details from a vpc", func(cmd *cli.Cmd) {
		flags := loadFlags(cmd, []string{"access_key", "secret_key", "region"})
		vpc := cmd.StringArg("VPC", "", "VPC id to inspect")

		command := func(c *context.Context) error {
			return context.Vpc(c, *vpc)
		}

		flagsCommand(cmd, command, flags)
	})

	app.Run(os.Args)
}
