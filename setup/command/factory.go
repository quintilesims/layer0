package command

import (
	"github.com/quintilesims/layer0/setup/layer0"
)

type CommandFactory struct {
	context layer0.Context
}

func NewCommandFactory(c layer0.Context) *CommandFactory {
	return &CommandFactory{
		context: c,
	}
}
