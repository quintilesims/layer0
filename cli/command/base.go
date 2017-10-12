package command

import (
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/urfave/cli"
)

type CommandBase struct {
	client   client.Client
	printer  printer.Printer
	resolver resolver.Resolver
}

func (b *CommandBase) SetClient(c client.Client) {
	b.client = c
}

func (b *CommandBase) SetPrinter(p printer.Printer) {
	b.printer = p
}

func (b *CommandBase) SetResolver(r resolver.Resolver) {
	b.resolver = r
}

func (b *CommandBase) waitOnJobHelper(c *cli.Context, jobID, spinnerText string, onCompleteFN func(entityID string) error) error {
	waitFlag := c.GlobalBool(config.FLAG_NO_WAIT)
	waitTimeout := c.GlobalDuration(config.FLAG_TIMEOUT)

	if waitFlag {
		b.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	b.printer.StartSpinner(spinnerText)
	defer b.printer.StopSpinner()

	job, err := client.WaitForJob(b.client, jobID, waitTimeout)
	if err != nil {
		return nil
	}

	return onCompleteFN(job.Result)
}

func (b *CommandBase) resolveSingleEntityIDHelper(entityType, target string) (string, error) {
	entityIDs, err := b.resolver.Resolve(entityType, target)
	if err != nil {
		return "", err
	}

	switch len(entityIDs) {
	case 0:
		return "", errors.NoMatchesError(entityType, target)
	case 1:
		return entityIDs[0], nil
	default:
		return "", errors.MultipleMatchesError(entityType, target, entityIDs)
	}
}

func (b *CommandBase) deleteHelper(c *cli.Context, entityType string, deleteFN func(entityID string) (string, error)) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	entityID, err := b.resolveSingleEntityIDHelper(entityType, args["NAME"])
	if err != nil {
		return err
	}

	jobID, err := deleteFN(entityID)
	if err != nil {
		return err
	}

	if c.GlobalBool(config.FLAG_NO_WAIT) {
		b.printJobResponse(jobID)
		return nil
	}

	b.printer.StartSpinner("deleting")
	defer b.printer.StopSpinner()

	if _, err := client.WaitForJob(b.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT)); err != nil {
		return err
	}

	b.printer.Printf("done")
	return nil
}

func (b *CommandBase) waitOnJobHelper(c *cli.Context, jobID, spinnerText string, onCompleteFN func(entityID string) error) error {
	if c.GlobalBool(config.FLAG_NO_WAIT) {
		b.printJobResponse(jobID)
		return nil
	}

	b.printer.StartSpinner(spinnerText)
	defer b.printer.StopSpinner()

	job, err := client.WaitForJob(b.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT))
	if err != nil {
		return err
	}

	return onCompleteFN(job.Result)
}

func (b *CommandBase) printJobResponse(jobID string) {
	b.printer.Printf("Operation is running as job '%s'", jobID)
}
