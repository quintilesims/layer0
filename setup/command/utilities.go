package command

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/instance"
	"github.com/urfave/cli"
)

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, fmt.Errorf("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func getInstance(f instance.Factory, c *cli.Context) (instance.Instance, error) {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return nil, err
	}

	return f.NewInstance(args["NAME"])
}
