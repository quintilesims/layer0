package command

import "github.com/urfave/cli"

type TaskCommand struct {
	*CommandBase
}

func NewTaskCommand(b *CommandBase) *TaskCommand {
	return &TaskCommand{b}
}

func (e *TaskCommand) Command() cli.Command {
	return cli.Command{}
}
