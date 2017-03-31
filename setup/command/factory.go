package command

import (
	"github.com/quintilesims/layer0/setup/instance"
)

type CommandFactory struct {
	InstanceFactory instance.Factory
}

func NewCommandFactory(f instance.Factory) *CommandFactory {
	return &CommandFactory{
		InstanceFactory: f,
	}
}
