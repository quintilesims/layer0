package instance

import (
	"fmt"
	"github.com/urfave/cli"
)

type ModuleInput struct {
	Name            string
	Default         interface{}
	LoadFromContext func(c *cli.Context) interface{}
}

var MainModuleInputs = []ModuleInput{
	{
		Name:            "source",
		Default:         "github.com/quintilesims/layer0/setup/module",
		LoadFromContext: stringOrNil("module-source"),
	},
	{
		Name:            "aws_access_key",
		LoadFromContext: stringOrNil("aws-access-key"),
	},
	{
		Name:            "aws_secret_key",
		LoadFromContext: stringOrNil("aws-secret-key"),
	},
	{
		Name:            "aws_region",
		Default:         "us-west-2",
		LoadFromContext: stringOrNil("aws-region"),
	},
}

func (m ModuleInput) Load(c *cli.Context, current interface{}) (interface{}, error) {
	if v := m.LoadFromContext(c); v != nil {
		return v, nil
	}

	if current != nil {
		return m.prompt(current)
	}

	return m.prompt(m.Default)
}

func (m ModuleInput) prompt(current interface{}) (interface{}, error) {
	prompt := fmt.Sprintf("%s: ", m.Name)
	if current != nil {
		prompt = fmt.Sprintf("%s [%v]: ", m.Name, current)
	}

	for i := 0; i < 3; i++ {
		fmt.Printf(prompt)
		var input string
		fmt.Scanln(&input)

		if input == "" && current != nil {
			return current, nil
		}

		// todo: input may need to be converted to a different type
		if input != "" {
			return input, nil
		}
	}

	return nil, fmt.Errorf("Failed to get input for '%s'", m.Name)
}

func stringOrNil(key string) func(c *cli.Context) interface{} {
	return func(c *cli.Context) interface{} {
		if v := c.String(key); v != "" {
			return v
		}

		return nil
	}
}
