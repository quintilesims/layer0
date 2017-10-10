package command

import "github.com/urfave/cli"

type ServiceCommand struct {
	*CommandBase
}

func NewServiceCommand(b *CommandBase) *ServiceCommand {
	return &ServiceCommand{b}
}

func (e *ServiceCommand) Command() cli.Command {
	return cli.Command{}
}
