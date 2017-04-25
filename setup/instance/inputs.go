package instance

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

const (
	INPUT_SOURCE         = "source"
	INPUT_AWS_ACCESS_KEY = "aws_access_key"
	INPUT_AWS_SECRET_KEY = "aws_secret_key"
	INPUT_AWS_REGION     = "aws_region"
	INPUT_AWS_KEY_PAIR   = "aws_key_pair"
	INPUT_VERSION        = "version"
	INPUT_DOCKERCFG      = "dockercfg"
)

type ModuleInput struct {
	Name        string
	Default     interface{}
	StaticValue interface{}
	// todo: add/display description
}

// todo: set source version
var MainModuleInputs = []ModuleInput{
	{
		Name:    INPUT_SOURCE,
		Default: "github.com/quintilesims/layer0/setup/module",
	},
	{
		Name:    INPUT_VERSION,
		Default: "latest",
	},
	{
		Name: INPUT_AWS_ACCESS_KEY,
	},
	{
		Name: INPUT_AWS_SECRET_KEY,
	},
	{
		Name:    INPUT_AWS_REGION,
		Default: "us-west-2",
	},
	{
		Name: INPUT_AWS_KEY_PAIR,
	},
	{
		Name:        INPUT_DOCKERCFG,
		StaticValue: "${file(\"dockercfg.json\")}",
	},
}

func (m ModuleInput) Prompt(current interface{}) (interface{}, error) {
	if m.StaticValue != nil {
		logrus.Errorf("%s: attempted to prompt an input with a static value\n", m.Name)
		return m.StaticValue, nil
	}

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
