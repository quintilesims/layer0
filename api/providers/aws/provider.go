package aws

import (
	"github.com/quintilesims/layer0/api/entity"
	"github.com/quintilesims/layer0/common/aws"
)

type AWSProvider struct {
	AWS *aws.Client
	// todo: Tag tag.Provider
	// todo: scheduler
}

func NewAWSProvider(a *aws.Client) *AWSProvider {
	return &AWSProvider{
		AWS: a,
	}
}

func (a *AWSProvider) ListEnvironmentIDs() ([]string, error) {
	return nil, nil
}

func (a *AWSProvider) GetEnvironment(environmentID string) entity.Environment {
	return NewEnvironment(a.AWS, environmentID)
}
