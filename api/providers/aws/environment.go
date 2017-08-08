package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
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
	if err := req.Validate(); err != nil {
		return errors.New(errors.InvalidRequest, err)
	}

	// todo: createSG
	// todo: createLC
	// todo: createASG
	// todo: createCluster

	// todo: just using id == name as placeholder
	e.id = req.EnvironmentName
	if err := e.createCluster(); err != nil {
		return err
	}

	return nil
}

func (e *Environment) createCluster() error {
	input := &ecs.CreateClusterInput{}
	input.SetClusterName(e.id)

	if _, err := e.AWS.ECS.CreateCluster(input); err != nil {
		return err
	}

	return nil
}

func (e *Environment) Delete() error {
	// todo: deleteSG
	// todo: deleteLC
	// todo: deleteASG
	// todo: deleteCluster
	return nil
}

func (e *Environment) Model() (*models.Environment, error) {
	// todo: readSG
	// todo: readLC
	// todo: readASG
	// todo: readCluster
	return nil, nil
}

func (e *Environment) Summary() (*models.EnvironmentSummary, error) {
	// todo: readSG
	// todo: readLC
	// todo: readASG
	// todo: readCluster
	return nil, nil
}
