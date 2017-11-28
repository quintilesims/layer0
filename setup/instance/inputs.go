package instance

import (
	"fmt"
	"log"
)

const LAYER0_MODULE_SOURCE = "github.com/quintilesims/layer0//setup/module"

const (
	INPUT_SOURCE           = "source"
	INPUT_VERSION          = "version"
	INPUT_AWS_ACCESS_KEY   = "access_key"
	INPUT_AWS_SECRET_KEY   = "secret_key"
	INPUT_AWS_REGION       = "region"
	INPUT_AWS_SSH_KEY_PAIR = "ssh_key_pair"
	INPUT_USERNAME         = "username"
	INPUT_PASSWORD         = "password"
	INPUT_DOCKERCFG        = "dockercfg"
	INPUT_VPC_ID           = "vpc_id"
)

const INPUT_SOURCE_DESCRIPTION = `
Source: The source input variable is the path to the Terraform module for Layer0.
By default, this points to the Layer0 github repository with the same version tag
as this l0-setup binary. Using values other than the default may result in 
undesired consequences (the double slash is intentional). 
`

const INPUT_VERSION_DESCRIPTION = `
Version: The version input variable specifies the tag to use for the Layer0 
Docker images 'quintilesims/l0-api'. This value
should match the version specified in the 'source' input variable. For example,
if the source is 'github.com/quintilesims/layer0/setup?ref=v1.2.3', the 
version should be 'v1.2.3'.
`

const INPUT_AWS_ACCESS_KEY_DESCRIPTION = `
AWS Access Key: The access_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources. 
`

const INPUT_AWS_SECRET_KEY_DESCRIPTION = `

AWS Secret Key: The secret_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources. 
`

const INPUT_AWS_REGION_DESCRIPTION = `
AWS Region: The region input variable specifies which region to provision the
AWS resources required for Layer0. The following regions can be used: 

    - us-west-1
    - us-west-2
    - us-east-1
    - eu-west-1

Note that changing this value will destroy and recreate any existing resources.
`

const INPUT_USERNAME_DESCRIPTION = `
Username: The username input variable specifies the name portion of basic auth used to secure 
the Layer0 API. This value will be base64 encoded along with the password input variable to
generate your LAYER0_AUTH_TOKEN.
`

const INPUT_PASSWORD_DESCRIPTION = `
Password: The password input variable specifies the password portion of basic auth used to secure
the Layer0 API. This value will be base64 encoded along with the username input variable to
generate your LAYER0_AUTH_TOKEN.
`

const INPUT_AWS_SSH_KEY_PAIR_DESCRIPTION = `
AWS SSH Key Pair: The ssh_key_pair input variable specifies the name of the 
ssh key pair to include in EC2 instances provisioned by Layer0. This key pair must 
already exist in the AWS account. The names of existing key pairs can be found
in the EC2 dashboard. Note that changing this value will not effect instances 
that have already been provisioned. 
`

const INPUT_VPC_ID_DESCRIPTION = `
VPC ID (optional): The vpc_id input variable specifies an existing AWS VPC to provision
the AWS resources required for Layer0. If no input is specified, a new VPC will be
created for you. Existing VPCs must satisfy the following constraints:

    - Have access to the public internet (through a NAT instance or gateway)
    - Have at least 1 public and 1 private subnet
    - Each subnet must be tagged with ["Tier": "Private"] or ["Tier": "Public"]

Note that changing this value will destroy and recreate any existing resources.
`

type ModuleInput struct {
	Name        string
	Description string
	Default     interface{}
	StaticValue interface{}
	prompter    func(ModuleInput, interface{}) (interface{}, error)
}

func InitializeLayer0ModuleInputs(version string) {
	if version == "" {
		log.Printf("[WARN] Version not set. Using default values for 'source' and 'version' inputs")
		return
	}

	for _, input := range Layer0ModuleInputs {
		switch input.Name {
		case INPUT_SOURCE:
			input.Default = fmt.Sprintf("%s?ref=%s", LAYER0_MODULE_SOURCE, version)
		case INPUT_VERSION:
			input.Default = version
		}
	}
}

var Layer0ModuleInputs = []*ModuleInput{
	{
		Name:        INPUT_SOURCE,
		Description: INPUT_SOURCE_DESCRIPTION,
		Default:     fmt.Sprintf("%s?ref=master", LAYER0_MODULE_SOURCE),
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_VERSION,
		Description: INPUT_VERSION_DESCRIPTION,
		Default:     "latest",
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_AWS_ACCESS_KEY,
		Description: INPUT_AWS_ACCESS_KEY_DESCRIPTION,
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_AWS_SECRET_KEY,
		Description: INPUT_AWS_SECRET_KEY_DESCRIPTION,
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_AWS_REGION,
		Description: INPUT_AWS_REGION_DESCRIPTION,
		Default:     "us-west-2",
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_AWS_SSH_KEY_PAIR,
		Description: INPUT_AWS_SSH_KEY_PAIR_DESCRIPTION,
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_USERNAME,
		Default:     "layer0",
		Description: INPUT_USERNAME_DESCRIPTION,
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_PASSWORD,
		Description: INPUT_PASSWORD_DESCRIPTION,
		prompter:    RequiredStringPrompter,
	},
	{
		Name:        INPUT_DOCKERCFG,
		StaticValue: "${file(\"dockercfg.json\")}",
	},
	{
		Name:        INPUT_VPC_ID,
		Description: INPUT_VPC_ID_DESCRIPTION,
		prompter:    OptionalStringPrompter,
	},
}

func (m ModuleInput) Prompt(current interface{}) (interface{}, error) {
	return m.prompter(m, current)
}

func RequiredStringPrompter(m ModuleInput, current interface{}) (interface{}, error) {
	return prompt(m, current, func(currentOrDefault interface{}) (interface{}, error) {
		for i := 0; i < 3; i++ {
			fmt.Printf("\tInput: ")

			var input string
			fmt.Scanln(&input)

			// user pressed 'enter' with a value already in place
			if input == "" && currentOrDefault != nil {
				return currentOrDefault, nil
			}

			// retry on empty input
			if input == "" {
				continue
			}

			return input, nil
		}

		return nil, fmt.Errorf("Failed to get input for '%s'", m.Name)
	})
}

func OptionalStringPrompter(m ModuleInput, current interface{}) (interface{}, error) {
	return prompt(m, current, func(currentOrDefault interface{}) (interface{}, error) {
		fmt.Printf("\tInput: ")

		var input string
		fmt.Scanln(&input)

		// user pressed 'enter' with a value already in place
		if input == "" && currentOrDefault != nil {
			return currentOrDefault, nil
		}

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
	} else {
		display = fmt.Sprintf("[current: <none>]\n")
		display += "Please enter a value and press 'enter'."
	}

	fmt.Println(display)
	return fn(currentOrDefault)
}
