package command

import (
	"fmt"
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local layer0 instances",
		Action: func(c *cli.Context) error {
			instances, err := f.InstanceFactory.ListInstances()
			if err != nil {
				return err
			}

			if len(instances) == 0 {
				fmt.Println("You don't have any local Layer0 instances")
			}

			for _, instance := range instances {
				fmt.Println(instance)
			}

			return nil
		},
	}
}
