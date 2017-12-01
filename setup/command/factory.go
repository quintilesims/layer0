package command

import (
	"github.com/aws/aws-sdk-go/aws/session"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/setup/instance"
)

type InstanceFactory func(string) instance.Instance
type AWSClientFactory func(session *session.Session) *awsc.Client

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
