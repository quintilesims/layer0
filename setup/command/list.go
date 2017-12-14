package command

import (
	"fmt"
	"sort"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

type status struct {
	Local  bool
	Remote bool
}

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local and/or remote Layer0 instances",
		Flags: []cli.Flag{
			config.FlagAWSAccessKey,
			config.FlagAWSSecretKey,
			config.FlagAWSRegion,
			cli.BoolFlag{
				Name:  "l, local",
				Usage: "only show local Layer0 instances, denoted by 'l'",
			},
			cli.BoolFlag{
				Name:  "r, remote",
				Usage: "only show remote Layer0 instances, denoted by 'r'",
			},
		},
		Action: func(c *cli.Context) error {
			instances := map[string]status{}
			if !c.Bool("local") {
				if err := f.addRemoteInstances(c, instances); err != nil {
					return err
				}
			}

			if !c.Bool("remote") {
				if err := f.addLocalInstances(instances); err != nil {
					return err
				}
			}

			fmt.Println("STATUS \t NAME")
			sortAndIterate(instances, func(name string, status status) {
				switch {
				case status.Local && !status.Remote:
					fmt.Printf("l \t %s\n", name)
				case !status.Local && status.Remote:
					fmt.Printf("r \t %s\n", name)
				default:
					fmt.Printf("lr \t %s\n", name)
				}
			})

			return nil
		},
	}
}

func (f *CommandFactory) addRemoteInstances(c *cli.Context, current map[string]status) error {
	provider, err := f.newAWSClientHelper(c)
	if err != nil {
		return err
	}

	remote, err := instance.ListRemoteInstances(provider.S3)
	if err != nil {
		return err
	}

	for _, r := range remote {
		v, ok := current[r]
		if !ok {
			current[r] = status{Remote: true}
			continue
		}

		current[r] = status{Remote: true, Local: v.Local}
	}

	return nil
}

func (f *CommandFactory) addLocalInstances(current map[string]status) error {
	local, err := instance.ListLocalInstances()
	if err != nil {
		return err
	}

	for _, l := range local {
		v, ok := current[l]
		if !ok {
			current[l] = status{Local: true}
			continue
		}

		current[l] = status{Remote: v.Remote, Local: true}
	}

	return nil
}

func sortAndIterate(instances map[string]status, fn func(string, status)) {
	names := []string{}
	for name := range instances {
		names = append(names, name)
	}

	sort.Strings(names)
	for _, name := range names {
		fn(name, instances[name])
	}
}
