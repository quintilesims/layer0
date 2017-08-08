package aws

import (
	"github.com/quintilesims/layer0/common/aws"
	"github.com/zpatrick/forge/common/models"
)

type Environment struct {
	*AWSEntity
}

func NewEnvironment(aws *aws.Provider, environmentID string) *Environment {
	return &Environment{
		NewAWSEntity(aws, environmentID),
	}
}

func (e *Environment) Create(req models.CreateEnvironmentRequest) error {
	// todo: if err := req.Validate(); err != nil { ... }

	// todo: createSG
	// todo: createLC
	// todo: createASG
	// todo: createCluster
	return nil
}

func (e *Environment) Delete() error {
	// todo: deleteSG
	// todo: deleteLC
	// todo: deleteASG
	// todo: deleteCluster
	return nil
}

func (e *Environment) Read() (*models.Environment, error) {
	// todo: readSG
	// todo: readLC
	// todo: readASG
	// todo: readCluster
	return nil, nil
}
