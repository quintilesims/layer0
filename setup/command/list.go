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
			// uses 'default' as default
			// if this doesn't exist yet, will prompt for it
			profile, err := f.context.LoadProfile(c.GlobalString("profile"))
			if err != nil {
				return err
			}

			local, err := f.context.ListLocalInstances(profile)
			if err != nil {
				return err
			}

			remote, err := f.context.ListRemoteInstances(profile)
			if err != nil {
				return err
			}

			instances := map[string]string{}
			for _, l := range local {
				instances[l] = "+l"
			}

			for _, r := range remote {
				if instances[r] == "" {
					instances[r] = "+"
				}

				instances[r] += "r"
			}

			for name, location := range instances {
				fmt.Printf("%s\t%s\n", name, location)
			}

			return nil
		},
	}
}
