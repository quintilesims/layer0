package command

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/setup/instance"
)

type InstanceFactory func(string) instance.Instance
type AWSClientFactory func(config *aws.Config) *setup_aws.Client

type CommandFactory struct {
	NewInstance  InstanceFactory
	NewAWSClient AWSClientFactory
}

func NewCommandFactory(i InstanceFactory, a AWSClientFactory) *CommandFactory {
	return &CommandFactory{
		NewInstance:  i,
		NewAWSClient: a,
	}
}
