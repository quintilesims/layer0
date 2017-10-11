package command

import "github.com/urfave/cli"

type AdminCommand struct {
	*CommandBase
}

func NewAdminCommand(b *CommandBase) *AdminCommand {
	return &AdminCommand{b}
}

func (e *AdminCommand) Command() cli.Command {
	return cli.Command{}
}
