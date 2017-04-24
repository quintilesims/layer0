package instance

import (
	"fmt"
)

const (
	INPUT_SOURCE         = "source"
	INPUT_AWS_ACCESS_KEY = "aws_access_key"
	INPUT_AWS_SECRET_KEY = "aws_secret_key"
	INPUT_AWS_REGION     = "aws_region"
	INPUT_AWS_KEY_PAIR   = "aws_key_pair"
)

type ModuleInput struct {
	Name    string
	Default interface{}
}

var MainModuleInputs = []ModuleInput{
	{
		Name:    INPUT_SOURCE,
		Default: "github.com/quintilesims/layer0/setup/module",
	},
	{
		Name: "aws_access_key",
	},
	{
		Name: "aws_secret_key",
	},
	{
		Name:    "aws_region",
		Default: "us-west-2",
	},
	{
		Name: "aws_key_pair",
	},
}

func (m ModuleInput) Prompt(current interface{}) (interface{}, error) {
	if current != nil {
		return m.prompt(current)
	}

	if m.Default != "" {
		return m.prompt(m.Default)
	}

	return m.prompt(nil)
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
