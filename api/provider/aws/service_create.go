package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Create(req models.CreateServiceRequest) (*models.Service, error) {
	var (
		cluster          string
		desiredCount     int
		loadBalancer     *ecs.LoadBalancer
		loadBalancerRole string
		serviceName      string
		taskDefinition   string
	)

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), req.EnvironmentID)
	cluster = fqEnvironmentID

	desiredCount = 1

	serviceID := generateEntityID(req.ServiceName)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	serviceName = fqServiceID

	taskDefinitionARN, err := s.lookupTaskDefinitionARNFromDeployID(req.DeployID)
	if err != nil {
		return nil, err
	}

	taskDefinition = taskDefinitionARN

	if req.LoadBalancerID != "" {
		fqLoadBalancerID := addLayer0Prefix(s.Config.Instance(), req.LoadBalancerID)
		loadBalancerDescription, err := describeLoadBalancer(s.AWS.ELB, fqLoadBalancerID)
		if err != nil {
			return nil, err
		}

		taskDefinitionDescription, err := describeTaskDefinition(s.AWS.ECS, taskDefinitionARN)
		if err != nil {
			return nil, err
		}

		var loadBalancerContainer ecs.LoadBalancer
		for _, containerDefinition := range taskDefinitionDescription.ContainerDefinitions {
			for _, portMapping := range containerDefinition.PortMappings {
				for _, listenerDescription := range loadBalancerDescription.ListenerDescriptions {
					listener := listenerDescription.Listener
					if *listener.InstancePort == *portMapping.ContainerPort {
						loadBalancerContainer = ecs.LoadBalancer{
							ContainerName:    containerDefinition.Name,
							ContainerPort:    portMapping.ContainerPort,
							LoadBalancerName: loadBalancerDescription.LoadBalancerName,
						}
					}
				}
			}
		}

		loadBalancer = &loadBalancerContainer
		loadBalancerRole = fmt.Sprintf("%s-lb", fqLoadBalancerID)
	}

	if err := s.createService(desiredCount, cluster, serviceName, taskDefinition, loadBalancerRole, loadBalancer); err != nil {
		return nil, err
	}

	if err := s.createTags(serviceID, req.ServiceName, req.EnvironmentID); err != nil {
		return nil, err
	}

	return s.Read(serviceID)
}

func (s *ServiceProvider) createService(desiredCount int, cluster, serviceName, taskDefinition, loadBalancerRole string, loadBalancer *ecs.LoadBalancer) error {
	input := &ecs.CreateServiceInput{}
	input.SetCluster(cluster)
	input.SetDesiredCount(int64(desiredCount))
	input.SetServiceName(serviceName)
	input.SetTaskDefinition(taskDefinition)
	if loadBalancer != nil && loadBalancerRole != "" {
		input.SetLoadBalancers([]*ecs.LoadBalancer{loadBalancer})
		input.SetRole(loadBalancerRole)
	}

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.CreateService(input); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) createTags(serviceID, serviceName, environmentID string) error {
	tags := []models.Tag{
		{
			EntityID:   serviceID,
			EntityType: "service",
			Key:        "name",
			Value:      serviceName,
		},
		{
			EntityID:   serviceID,
			EntityType: "service",
			Key:        "environment_id",
			Value:      environmentID,
		},
	}

	for _, tag := range tags {
		if err := s.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceProvider) lookupTaskDefinitionARNFromDeployID(deployID string) (string, error) {
	tags, err := s.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("arn").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Could not resolve task definition ARN for deploy '%s'", deployID)
}
