package command

import "github.com/urfave/cli"

type JobCommand struct {
	*CommandBase
}

func NewJobCommand(b *CommandBase) *JobCommand {
	return &JobCommand{b}
}

func (e *JobCommand) Command() cli.Command {
	return cli.Command{}
}
