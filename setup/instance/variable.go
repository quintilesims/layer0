package instance

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/terraform"
	"github.com/urfave/cli"
	"os"
)

type VariableSchema struct {
	Name            string
	Default         interface{}
	CLIFlag         cli.Flag
	LoadFromContext func(c *cli.Context) interface{}
}

var InstanceVariableSchemas = []VariableSchema{
	{
		Name: "aws_access_key",
		LoadFromContext: func(c *cli.Context) interface{} {
			if v := c.String("aws-access-key"); v != "" {
				return v
			}

			return nil
		},
	},
	{
		Name: "aws_secret_key",
		LoadFromContext: func(c *cli.Context) interface{} {
			if v := c.String("aws-secret-key"); v != "" {
				return v
			}

			return nil
		},
	},
	{
		Name:    "aws_region",
		Default: "us-west-2",
		LoadFromContext: func(c *cli.Context) interface{} {
			if v := c.String("aws-region"); v != "" {
				return v
			}

			return nil
		},
	},
}

func (v *VariableSchema) Load(c *cli.Context, path string) (interface{}, error) {
	if val := v.LoadFromContext(c); val != nil {
		return val, nil
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		tfvars, err := terraform.LoadTFVars(path)
		if err != nil {
			return nil, err
		}

		if current, ok := tfvars[v.Name]; ok {
			return v.prompt(v.Name, current)
		}
	}

	return v.prompt(v.Name, v.Default)
}

// may need to make a MustLoad() function later

func (v *VariableSchema) prompt(name, current interface{}) (interface{}, error) {
	prompt := fmt.Sprintf("%s: ", name)
	if current != nil {
		prompt = fmt.Sprintf("%s [%v]: ", name, current)
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

	return nil, fmt.Errorf("Failed to get input for '%s'", name)
}
