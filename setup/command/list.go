package command

import (
	"sort"
	"fmt"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func (f *CommandFactory) List() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list local and remote layer0 instances",
		Flags: awsFlags,
		Action: func(c *cli.Context) error {
			provider, err := f.newAWSProviderHelper(c)
			if err != nil {
				return err
			}

			remote, err := instance.ListRemoteInstances(provider.S3)
			if err != nil {
				return err
			}

			local, err := instance.ListLocalInstances()
			if err != nil {
				return err
			}

			// print 'l' if local, 'r' if remote, 'lr' if both
			catalog := map[string]string{}
			for _, instance := range local {
				catalog[instance] = "l "
			}

			for _, instance := range remote {
				if _, ok := catalog[instance]; ok {
					catalog[instance] = "lr"
				} else {
					catalog[instance] = " r"
				}
			}

			sortAndIterate(catalog, func(instance, token string) {
				fmt.Printf("%s\t%s\n", token, instance)
			})

			return nil
		},
	}
}

func sortAndIterate(m map[string]string, fn func(string, string)) {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fn(key, m[key])
	}
}
