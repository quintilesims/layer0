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
	INPUT_VPC_ID         = "vpc_id"
)

const INPUT_SOURCE_DESCRIPTION = `
Source: The source input variable is the path to the Terraform module for Layer0.
By default, this points to the Layer0 github repository with the same version tag
as this l0-setup binary. Using values other than the default may result in 
undesired consequences. 
`

const INPUT_VERSION_DESCRIPTION = `
Version: The version input variable blah blah blah
blah blah blah
`

const INPUT_AWS_ACCESS_KEY_DESCRIPTION = `
AWS Access Key: The aws_access_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources. 
`

const INPUT_AWS_SECRET_KEY_DESCRIPTION = `

AWS Secret Key: The aws_secret_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources. 
`

const INPUT_AWS_REGION_DESCRIPTION = `
AWS Region: The aws_region input variable specifies which region to provision the
AWS resources required for Layer0. Note that changing this value will destroy and 
recreate any existing resources.
`

const INPUT_AWS_KEY_PAIR_DESCRIPTION = `
Version: The key_pair input variable blah blah blah
blah blah blah
`

const INPUT_VPC_ID_DESCRIPTION = `
Version: The vpc_id input variable blah blah blah
blah blah blah
`

type ModuleInput struct {
	Name        string
	Description string
	Default     interface{}
	StaticValue interface{}
	prompter    func(ModuleInput, interface{}) (interface{}, error)
}

func InitializeMainModuleInputs(version string) {
	if version == "" {
		logrus.Warningf("Version not set. Using default values for 'source' and 'version' inputs")
		return
	}

	for _, input := range MainModuleInputs {
		switch input.Name {
		case INPUT_SOURCE:
			input.Default = fmt.Sprintf("github.com/quintilesims/layer0/setup/module?ref=%s", version)
		case INPUT_VERSION:
			input.Default = version
		}
	}
}

// todo: set source version
var MainModuleInputs = []*ModuleInput{
	{
		Name:        INPUT_SOURCE,
		Description: INPUT_SOURCE_DESCRIPTION,
		Default:     "github.com/quintilesims/layer0/setup/module?ref=master",
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_VERSION,
		Description: INPUT_VERSION_DESCRIPTION,
		Default:     "latest",
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_AWS_ACCESS_KEY,
		Description: INPUT_AWS_ACCESS_KEY_DESCRIPTION,
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_AWS_SECRET_KEY,
		Description: INPUT_AWS_SECRET_KEY_DESCRIPTION,
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_AWS_REGION,
		Description: INPUT_AWS_REGION_DESCRIPTION,
		Default:     "us-west-2",
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_AWS_KEY_PAIR,
		Description: INPUT_AWS_KEY_PAIR_DESCRIPTION,
		prompter:    DefaultStringPrompter,
	},
	{
		Name:        INPUT_DOCKERCFG,
		StaticValue: "${file(\"dockercfg.json\")}",
	},
	{
		Name:        INPUT_VPC_ID,
		Description: INPUT_VPC_ID_DESCRIPTION,
		prompter:    VPCPrompter,
	},
}

func (m ModuleInput) Prompt(current interface{}) (interface{}, error) {
	return m.prompter(m, current)
}

func DefaultStringPrompter(m ModuleInput, current interface{}) (interface{}, error) {
	return prompt(m, current, func(currentOrDefault interface{}) (interface{}, error) {
		for i := 0; i < 3; i++ {
			fmt.Printf("\tInput: ")

			var input string
			fmt.Scanln(&input)

			// user pressed 'enter' with a value already in place
			if input == "" && currentOrDefault != nil {
				return currentOrDefault, nil
			}

			if input != "" {
				return input, nil
			}
		}

		return nil, fmt.Errorf("Failed to get input for '%s'", m.Name)
	})
}

func VPCPrompter(m ModuleInput, current interface{}) (interface{}, error) {
	return prompt(m, current, func(currentOrDefault interface{}) (interface{}, error) {
		fmt.Printf("\tInput: ")

		var input string
		fmt.Scanln(&input)

		// user pressed 'enter' with a value already in place
		if input == "" && currentOrDefault != nil {
			return currentOrDefault, nil
		}

		// empty input is ok for vpc
		return input, nil
	})
}

func prompt(m ModuleInput, current interface{}, fn func(interface{}) (interface{}, error)) (interface{}, error) {
	fmt.Println(m.Description)

	var display string
	var currentOrDefault interface{}
	if current != nil {
		display = fmt.Sprintf("[current: %v]\n", current)
		display += "Please enter a new value, or press 'enter' to keep the current value."
		currentOrDefault = current
	} else if m.Default != nil {
		display = fmt.Sprintf("[default: %v]\n", m.Default)
		display += "Please enter a new value, or press 'enter' to use the default value."
		currentOrDefault = m.Default
	}

	fmt.Println(display)
	return fn(currentOrDefault)
}
