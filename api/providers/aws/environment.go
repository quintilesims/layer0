package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type Environment struct {
	*AWSEntity
}

func NewEnvironment(aws *awsc.Client, environmentID string) *Environment {
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

	if err := e.deleteCluster(); err != nil {
		return err
	}

	return nil
}

func (e *Environment) deleteCluster() error {
	input := &ecs.DeleteClusterInput{}
	input.SetCluster(e.id)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.ECS.DeleteCluster(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ClusterNotFoundException" {
			return nil
		}

		return err
	}

	return nil
}

func (e *Environment) Read() (*models.Environment, error) {
	if _, err := e.readCluster(); err != nil {
		return nil, err
	}

	model := &models.Environment{
		EnvironmentID: e.id,
	}

	return model, nil
}

func (e *Environment) readCluster() (*ecs.Cluster, error) {
	input := &ecs.DescribeClustersInput{}
	input.SetClusters([]*string{aws.String(e.id)})

	output, err := e.AWS.ECS.DescribeClusters(input)
	if err != nil {
		return nil, err
	}

	for _, cluster := range output.Clusters {
		if aws.StringValue(cluster.ClusterName) == e.id && aws.StringValue(cluster.Status) == "ACTIVE" {
			return cluster, nil
		}
	}

	return nil, errors.Newf(errors.EnvironmentDoesNotExist, "Environment %s does not exist", e.id)
}

func (e *Environment) Model() (*models.Environment, error) {
	if _, err := e.readCluster(); err != nil {
		return nil, err
	}

	model := &models.Environment{
		EnvironmentID: e.id,
	}

	return model, nil
}

func (e *Environment) Summary() (*models.EnvironmentSummary, error) {
	if _, err := e.readCluster(); err != nil {
		return nil, err
	}

	summary := &models.EnvironmentSummary{
		EnvironmentID: e.id,
	}

	return summary, nil
}
