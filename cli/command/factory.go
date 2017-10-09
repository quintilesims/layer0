package command

import (
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type CommandFactory struct {
	client   client.Client
	printer  printer.Printer
	resolver resolver.Resolver
}

func NewCommandFactory(c client.Client, p printer.Printer, r resolver.Resolver) *CommandFactory {
	return &CommandFactory{
		client:   c,
		printer:  p,
		resolver: r,
	}
}

func (f *CommandFactory) SetClient(c client.Client) {
	f.client = c
}

func (f *CommandFactory) SetPrinter(p printer.Printer) {
	f.printer = p
}

func (f *CommandFactory) SetResolver(r resolver.Resolver) {
	f.resolver = r
}

func (f *CommandFactory) deleteHelper(c *cli.Context, entityType string, deleteFN func(entityID string) (string, error)) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	entityID, err := resolveSingleEntityID(f.resolver, entityType, args["NAME"])
	if err != nil {
		return err
	}

	jobID, err := deleteFN(entityID)
	if err != nil {
		return err
	}

	if c.GlobalBool(config.FLAG_NO_WAIT) {
		// todo: use single 'running as job' helper printer
		f.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	f.printer.StartSpinner("deleting")
	defer f.printer.StopSpinner()

	if _, err := client.WaitForJob(f.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT)); err != nil {
		return err
	}

	return nil
}
