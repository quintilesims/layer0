package command

import "github.com/urfave/cli"

type LoadBalancerCommand struct {
	*CommandBase
}

func NewLoadBalancerCommand(b *CommandBase) *LoadBalancerCommand {
	return &LoadBalancerCommand{b}
}

func (e *LoadBalancerCommand) Command() cli.Command {
	return cli.Command{}
}
