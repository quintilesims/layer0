package aws

import "github.com/quintilesims/layer0/common/aws"

type AWSEntity struct {
	AWS *aws.Client
	id  string
	// todo: tag.Provider
	// todo: tag.Type?
}

func NewAWSEntity(aws *aws.Client, id string) *AWSEntity {
	return &AWSEntity{
		AWS: aws,
		id:  id,
	}
}

func (e *AWSEntity) ID() string {
	return e.id
}
