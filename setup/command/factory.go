package command

import (
	"github.com/quintilesims/layer0/setup/aws"
	"github.com/quintilesims/layer0/setup/instance"
)

// todo: inject aws factory
// todo: inject instance factory
type CommandFactory struct {
	NewInstance    func(string) instance.Instance
	NewAWSProvider func(string, string, string) *aws.Provider
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		NewInstance:    defaultInstanceFactory,
		NewAWSProvider: defaultAWSProviderFactory,
	}
}

func defaultInstanceFactory(name string) instance.Instance {
	return instance.NewLocalInstance(name)
}

func defaultAWSProviderFactory(accessKey, secretKey, region string) *aws.Provider {
	return aws.NewProvider(accessKey, secretKey, region)
}
