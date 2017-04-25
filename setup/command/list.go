package command

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local and remote layer0 instances",
		Flags: s3Flags,
		Action: func(c *cli.Context) error {
			s3, err := newS3(c)
			if err != nil {
				return err
			}

			remote, err := instance.ListRemoteInstances(s3)
			if err != nil {
				return err
			}

			local, err := instance.ListLocalInstances()
			if err != nil {
				return err
			}

			catalog := map[string]string{}
			for _, instance := range local {
				catalog[instance] += "l"
			}

			for _, instance := range remote {
				catalog[instance] += "r"
			}

			for instance, token := range catalog {
				fmt.Printf("%2s    %s\n", token, instance)
			}

			return nil
		},
	}
}
