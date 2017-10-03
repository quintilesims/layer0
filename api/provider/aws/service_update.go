package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Update updates an ECS Service using the specified Update Service Request and returns an error.
// The Update Service Request contains the Service ID, the Deploy ID and the scale of the DesiredCount
// of the Service. The Service's Task Definition and DesiredCount are updated with two separate UpdateService
// requests to AWS.
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
		deployID := *req.DeployID
		deployName, deployVersion, err := lookupDeployNameAndVersion(s.TagStore, deployID)
		if err != nil {
			return err
		}

		taskDefinitionFamily := deployName
		taskDefinitionRevision := deployVersion

		if err := s.updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionFamily, taskDefinitionRevision); err != nil {
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

func (s *ServiceProvider) updateServiceTaskDefinition(clusterName, serviceName, taskDefinitionFamily, taskDefinitionRevision string) error {
	input := &ecs.UpdateServiceInput{}
	input.SetCluster(clusterName)
	input.SetService(serviceName)

	taskDefinition := fmt.Sprintf("%s:%s", taskDefinitionFamily, taskDefinitionRevision)
	input.SetTaskDefinition(taskDefinition)

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
