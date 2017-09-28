package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// Delete a specified service within a cluster.
func (s *ServiceProvider) Delete(serviceID string) error {
	desiredCount := int64(0)

	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		return err
	}
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	fqClusterID := addLayer0Prefix(s.Config.Instance(), environmentID)

	if err := s.scaleService(fqServiceID, fqClusterID, desiredCount); err != nil {
		return fmt.Errorf("Failed to scale service: %e", err)
	}

	if err := s.deleteService(fqServiceID, fqClusterID); err != nil {
		return fmt.Errorf("Failed to delete service: %e", err)
	}

	if err := deleteEntityTags(s.TagStore, "service", serviceID); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) scaleService(serviceID, clusterID string, desiredCount int64) error {
	input := &ecs.UpdateServiceInput{}
	input.DesiredCount = &desiredCount
	input.SetCluster(clusterID)
	input.SetService(serviceID)

	if err := input.Validate(); err != nil {
		return fmt.Errorf("Error failed to validate: %e", err)
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return fmt.Errorf("Error scaling service: %e", err)
	}
	return nil
}

func (s *ServiceProvider) deleteService(serviceID, clusterID string) error {
	input := &ecs.DeleteServiceInput{}
	input.SetCluster(clusterID)
	input.SetService(serviceID)

	if err := input.Validate(); err != nil {
		return fmt.Errorf("Error failed to validate: %e", err)
	}

	if _, err := s.AWS.ECS.DeleteService(input); err != nil {
		return fmt.Errorf("Error deleting service: %e", err)
	}

	return nil
}
