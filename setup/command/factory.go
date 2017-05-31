package command

import (
	"github.com/aws/aws-sdk-go/aws"
	setup_aws "github.com/quintilesims/layer0/setup/aws"
	"github.com/quintilesims/layer0/setup/instance"
)

type InstanceFactory func(string) instance.Instance
type AWSProviderFactory func(config *aws.Config) *setup_aws.Provider

type CommandFactory struct {
	NewInstance    InstanceFactory
	NewAWSProvider AWSProviderFactory
}

func NewCommandFactory(i InstanceFactory, a AWSProviderFactory) *CommandFactory {
	return &CommandFactory{
		NewInstance:    i,
		NewAWSProvider: a,
	}
}
