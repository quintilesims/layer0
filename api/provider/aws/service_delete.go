package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
)

// Delete a specified service within a cluster.
func (s *ServiceProvider) Delete(serviceID string) error {
	desiredCount := 0

	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.ServiceDoesNotExist {
			return nil
		}

		return err
	}

	clusterName := addLayer0Prefix(s.Context, environmentID)
	fqServiceID := addLayer0Prefix(s.Context, serviceID)

	service, err := readService(s.AWS.ECS, clusterName, fqServiceID)
	if err != nil {
		return err
	}

	taskARNs, err := s.getServiceActiveTasks(clusterName, service.Deployments)
	if err != nil {
		return err
	}

	if err := s.stopServiceTasks(clusterName, taskARNs); err != nil {
		return err
	}

	if err := s.scaleService(clusterName, fqServiceID, desiredCount); err != nil {
		return err
	}

	if err := s.deleteService(clusterName, fqServiceID); err != nil {
		return err
	}

	if err := deleteEntityTags(s.TagStore, "service", serviceID); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) getServiceActiveTasks(clusterName string, deployments []*ecs.Deployment) ([]string, error) {
	taskARNs := []string{}

	for _, deployment := range deployments {
		clusterTaskARNsRunning, err := listClusterTaskARNs(s.AWS.ECS, clusterName, aws.StringValue(deployment.Id), ecs.DesiredStatusRunning)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, clusterTaskARNsRunning...)
	}

	return taskARNs, nil
}

func (s *ServiceProvider) stopServiceTasks(clusterName string, taskARNs []string) error {
	for _, taskARN := range taskARNs {
		input := &ecs.StopTaskInput{}
		input.SetCluster(clusterName)
		input.SetTask(taskARN)

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := s.AWS.ECS.StopTask(input); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceProvider) scaleService(clusterName, serviceID string, desiredCount int) error {
	input := &ecs.UpdateServiceInput{}

	input.SetDesiredCount(int64(desiredCount))
	input.SetCluster(clusterName)
	input.SetService(serviceID)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) deleteService(clusterName, serviceID string) error {
	input := &ecs.DeleteServiceInput{}
	input.SetCluster(clusterName)
	input.SetService(serviceID)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.DeleteService(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ServiceNotFoundException" {
			return nil
		}

		return err
	}

	return nil
}
