package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Update is used to update an ECS Service using the specified Update Service Request.
// The Update Service Request contains the Service ID, the Deploy ID and the scale of the DesiredCount
// of the Service. The Service's Task Definition and Desired Count are updated with two separate UpdateService
// requests to AWS.
func (s *ServiceProvider) Update(serviceID string, req models.UpdateServiceRequest) error {
	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		return err
	}

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID

	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	serviceName := fqServiceID

	if req.DeployID != nil {
		deployID := *req.DeployID
		taskDefinitionARN, err := lookupTaskDefinitionARNFromDeployID(s.TagStore, deployID)
		if err != nil {
			return err
		}

		if err := s.updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionARN); err != nil {
			return err
		}
	}

	if req.Scale != nil {
		desiredCount := *req.Scale
		if err := s.updateServiceDesiredCount(clusterName, serviceName, desiredCount); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceProvider) updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionARN string) error {
	input := &ecs.UpdateServiceInput{}
	input.SetCluster(clusterName)
	input.SetService(serviceName)
	input.SetTaskDefinition(taskDefinitionARN)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) updateServiceDesiredCount(clusterName, serviceName string, desiredCount int) error {
	input := &ecs.UpdateServiceInput{}
	input.SetCluster(clusterName)
	input.SetService(serviceName)
	input.SetDesiredCount(int64(desiredCount))

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return err
	}

	return nil
}
