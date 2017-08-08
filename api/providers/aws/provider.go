package aws

import (
	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/common/aws"
)

type AWSProvider struct {
	AWS *aws.Provider
	// todo: Tag tag.Provider
	// todo: scheduler
}

func NewAWSProvider(a *aws.Provider) *AWSProvider {
	return &AWSProvider{
		AWS: a,
	}
}

func (a *AWSProvider) GetEnvironment(environmentID string) entity.Environment {
	return NewEnvironment(a.AWS, environmentID)
}
