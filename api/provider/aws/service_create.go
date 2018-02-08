package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Create(req models.CreateServiceRequest) (string, error) {
	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), req.EnvironmentID)
	cluster := fqEnvironmentID

	var securityGroupIDs []*string
	environmentSecurityGroupName := getEnvironmentSGName(fqEnvironmentID)
	environmentSecurityGroup, err := readSG(s.AWS.EC2, environmentSecurityGroupName)
	if err != nil {
		return "", err
	}

	securityGroupIDs = append(securityGroupIDs, environmentSecurityGroup.GroupId)

	launchType, err := getLaunchTypeFromEnvironmentID(s.TagStore, req.EnvironmentID)
	if err != nil {
		return "", err
	}

	serviceID := entityIDGenerator(req.ServiceName)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	serviceName := fqServiceID

	taskDefinitionARN, err := lookupTaskDefinitionARNFromDeployID(s.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	var loadBalancer *ecs.LoadBalancer
	var loadBalancerRole string
	if req.LoadBalancerID != "" {
		fqLoadBalancerID := addLayer0Prefix(s.Config.Instance(), req.LoadBalancerID)
		loadBalancerDescription, err := describeLoadBalancer(s.AWS.ELB, fqLoadBalancerID)
		if err != nil {
			return "", err
		}

		if aws.StringValue(loadBalancerDescription.Scheme) == "internet-facing" {
			loadBalancerSecurityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
			loadBalancerSecurityGroup, err := readSG(s.AWS.EC2, loadBalancerSecurityGroupName)
			if err != nil {
				return "", err
			}

			securityGroupIDs = append(securityGroupIDs, loadBalancerSecurityGroup.GroupId)
		}

		taskDefinition, err := describeTaskDefinition(s.AWS.ECS, taskDefinitionARN)
		if err != nil {
			return "", err
		}

		for _, containerDefinition := range taskDefinition.ContainerDefinitions {
			for _, portMapping := range containerDefinition.PortMappings {
				for _, listenerDescription := range loadBalancerDescription.ListenerDescriptions {
					listener := listenerDescription.Listener
					if aws.Int64Value(listener.InstancePort) == aws.Int64Value(portMapping.ContainerPort) {
						loadBalancer = &ecs.LoadBalancer{
							ContainerName:    containerDefinition.Name,
							ContainerPort:    portMapping.ContainerPort,
							LoadBalancerName: loadBalancerDescription.LoadBalancerName,
						}

						loadBalancerRole = fmt.Sprintf("%s-lb", fqLoadBalancerID)
					}
				}
			}
		}
	}

	scale := req.Scale
	if req.Scale == 0 {
		scale = 1
	}

	privateSubnets := s.Config.PrivateSubnets()
	subnets := make([]*string, len(privateSubnets))
	for i := range privateSubnets {
		subnets[i] = aws.String(privateSubnets[i])
	}

	if err := s.createService(
		cluster,
		launchType,
		serviceName,
		taskDefinitionARN,
		scale,
		loadBalancerRole,
		loadBalancer,
		securityGroupIDs,
		subnets); err != nil {
		return "", err
	}

	if err := s.createTags(serviceID, req.ServiceName, req.EnvironmentID, req.LoadBalancerID); err != nil {
		return "", err
	}

	return serviceID, nil
}

func (s *ServiceProvider) createService(
	cluster string,
	launchType string,
	serviceName string,
	taskDefinition string,
	desiredCount int,
	loadBalancerRole string,
	loadBalancer *ecs.LoadBalancer,
	securityGroupIDs []*string,
	subnets []*string,
) error {
	input := &ecs.CreateServiceInput{}
	input.SetCluster(cluster)
	input.SetDesiredCount(int64(desiredCount))
	input.SetLaunchType(launchType)
	input.SetServiceName(serviceName)
	input.SetTaskDefinition(taskDefinition)

	if launchType == ecs.LaunchTypeFargate {
		awsvpcConfig := &ecs.AwsVpcConfiguration{}
		awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
		awsvpcConfig.SetSecurityGroups(securityGroupIDs)
		awsvpcConfig.SetSubnets(subnets)
		networkConfig := &ecs.NetworkConfiguration{}
		networkConfig.SetAwsvpcConfiguration(awsvpcConfig)
		input.SetNetworkConfiguration(networkConfig)
		input.SetPlatformVersion("LATEST")
	}

	if loadBalancer != nil {
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

func (s *ServiceProvider) createTags(serviceID, serviceName, environmentID, loadBalancerID string) error {
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

	if loadBalancerID != "" {
		tag := models.Tag{
			EntityID:   serviceID,
			EntityType: "service",
			Key:        "load_balancer_id",
			Value:      loadBalancerID,
		}

		tags = append(tags, tag)
	}

	for _, tag := range tags {
		if err := s.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
