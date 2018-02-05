package command

import (
	"github.com/quintilesims/layer0/cli/printer"
	"github.com/quintilesims/layer0/cli/resolver"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/errors"
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
