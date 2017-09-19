package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Update(req models.UpdateServiceRequest) error {
	serviceID := req.ServiceID

	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		return err
	}

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID

	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	serviceName := fqServiceID

	if req.DeployID != nil {
		taskDefinitionID := *req.DeployID
		if err := s.updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionID); err != nil {
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

func (s *ServiceProvider) updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionID string) error {
	input := &ecs.UpdateServiceInput{}
	input.SetCluster(clusterName)
	input.SetService(serviceName)
	input.SetTaskDefinition(taskDefinitionID)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return err
	}

	return nil
}
