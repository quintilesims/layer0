package command

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local layer0 instances",
		Action: func(c *cli.Context) error {
			instances, err := instance.ListLocalInstances()
			if err != nil {
				return err
			}

			for _, instance := range instances {
				fmt.Println(instance)
			}

			return nil
		},
	}
}
