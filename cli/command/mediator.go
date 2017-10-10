package command

import (
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
)

type CommandMediator struct {
	client   client.Client
	printer  printer.Printer
	resolver resolver.Resolver
}

func (m *CommandMediator) SetClient(c client.Client) {
	m.client = c
}

func (m *CommandMediator) SetPrinter(p printer.Printer) {
	m.printer = p
}

func (m *CommandMediator) SetResolver(r resolver.Resolver) {
	m.resolver = r
}

func (m *CommandMediator) deleteHelper(c *cli.Context, entityType string, deleteFN func(entityID string) (string, error)) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	entityID, err := resolveSingleEntityID(m.resolver, entityType, args["NAME"])
	if err != nil {
		return err
	}

	jobID, err := deleteFN(entityID)
	if err != nil {
		return err
	}

	if c.GlobalBool(config.FLAG_NO_WAIT) {
		// todo: use single 'running as job' helper printer
		m.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	m.printer.StartSpinner("deleting")
	defer m.printer.StopSpinner()

	if _, err := client.WaitForJob(m.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT)); err != nil {
		return err
	}

	return nil
}

func (m *CommandMediator) waitOnJobHelper(c *cli.Context, jobID, spinnerText string, onCompleteFN func(entityID string) error) error {
	waitFlag := c.GlobalBool(config.FLAG_NO_WAIT)
	waitTimeout := c.GlobalDuration(config.FLAG_TIMEOUT)

	if waitFlag {
		m.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	m.printer.StartSpinner(spinnerText)
	defer m.printer.StopSpinner()

	job, err := client.WaitForJob(m.client, jobID, waitTimeout)
	if err != nil {
		return err
	}

	return onCompleteFN(job.Result)
}
