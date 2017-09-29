package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// Delete a specified service within a cluster.
func (s *ServiceProvider) Delete(serviceID string) error {
	desiredCount := 0

	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		if strings.Contains(err.Error(), "ServiceDoesNotExist") {
			return nil
		}
		return err
	}

	clusterName := addLayer0Prefix(s.Config.Instance(), environmentID)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)

	service, err := s.readService(clusterName, serviceID)
	if err != nil {
		return err
	}

	taskARNs, err := s.getARNTasks(clusterName, serviceID, service.Deployments)
	if err != nil {
		return err
	}

	if err := s.stopARNTasks(clusterName, fqServiceID, taskARNs); err != nil {
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

func (s *ServiceProvider) getARNTasks(clusterName, serviceID string,
	deployments []*ecs.Deployment) ([]*string, error) {
	taskARNs := []*string{}

	for _, deployment := range deployments {
		running := "RUNNING"
		stopped := "STOPPED"

		input := &ecs.ListTasksInput{}
		input.SetCluster(clusterName)
		input.SetDesiredStatus(running)
		input.SetStartedBy(*deployment.Id)

		tasks, err := s.AWS.ECS.ListTasks(input)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, tasks.TaskArns...)

		input.SetDesiredStatus(stopped)
		tasks, err = s.AWS.ECS.ListTasks(input)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, tasks.TaskArns...)
	}
	return taskARNs, nil
}

func (s *ServiceProvider) stopARNTasks(clusterName, serviceID string, taskARNs []*string) error {
	for i := range taskARNs {
		taskARN := taskARNs[i]
		inputTask := &ecs.StopTaskInput{}
		inputTask.SetCluster(clusterName)
		inputTask.SetTask(*taskARN)

		if err := inputTask.Validate(); err != nil {
			return err
		}

		_, err := s.AWS.ECS.StopTask(inputTask)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceProvider) scaleService(clusterName, serviceID string, desiredCount int) error {
	input := &ecs.UpdateServiceInput{}
	awsDesiredCount := int64(desiredCount)

	input.DesiredCount = &awsDesiredCount
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
		return err
	}

	return nil
}
