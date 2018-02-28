package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// Create runs an ECS Service using the specified CreateServiceRequest.
// The CreateServiceRequest contains the DeployID, the EnvironmentID,
// the LoadBalancerID (optional), the Scale, the ServiceName, and the ServiceType.
//
// The DeployID is used to look up the ARN of the ECS TaskDefinition to run. If a
// LoadbalancerID is supplied, it will be used in conjunction with the TaskDefinition
// ARN to compare the ports specified in the TaskDefinition with those specified on
// the LoadBalancer. The ServiceType parameter is one of "stateful" or "stateless"
// and indicates which ECS LaunchType the user wishes to use ("EC2" or "FARGATE"
// respectively).
//
// Create does not generate any custom errors of its own, but will bubble up errors
// found in its helper functions as well as errors returned by AWS.
func (s *ServiceProvider) Create(req models.CreateServiceRequest) (string, error) {
	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), req.EnvironmentID)
	cluster := fqEnvironmentID
	serviceID := entityIDGenerator(req.ServiceName)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	serviceName := fqServiceID

	scale := req.Scale
	if req.Scale == 0 {
		scale = 1
	}

	taskDefinitionARN, err := lookupTaskDefinitionARNFromDeployID(s.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	var fargatePlatformVersion string
	var securityGroupIDs []*string
	var subnets []string
	if req.ServiceType == models.DeployCompatibilityStateless {
		fargatePlatformVersion = config.DefaultFargatePlatformVersion

		environmentSecurityGroupName := getEnvironmentSGName(fqEnvironmentID)
		environmentSecurityGroup, err := readSG(s.AWS.EC2, environmentSecurityGroupName)
		if err != nil {
			return "", err
		}

		securityGroupIDs = append(securityGroupIDs, environmentSecurityGroup.GroupId)

		subnets = s.Config.PrivateSubnets()
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

	if err := s.createService(
		cluster,
		req.ServiceType,
		serviceName,
		taskDefinitionARN,
		loadBalancerRole,
		fargatePlatformVersion,
		scale,
		subnets,
		securityGroupIDs,
		loadBalancer,
	); err != nil {
		return "", err
	}

	if err := s.createTags(serviceID, req.ServiceName, req.EnvironmentID, req.LoadBalancerID); err != nil {
		return "", err
	}

	return serviceID, nil
}

func (s *ServiceProvider) createService(
	cluster,
	serviceType,
	serviceName,
	taskDefinition,
	loadBalancerRole,
	fargatePlatformVersion string,
	desiredCount int,
	subnets []string,
	securityGroupIDs []*string,
	loadBalancer *ecs.LoadBalancer,
) error {
	input := &ecs.CreateServiceInput{}
	input.SetCluster(cluster)
	input.SetDesiredCount(int64(desiredCount))
	input.SetServiceName(serviceName)
	input.SetTaskDefinition(taskDefinition)

	launchType := ecs.LaunchTypeEc2
	if serviceType == models.DeployCompatibilityStateless {
		launchType = ecs.LaunchTypeFargate

		s := make([]*string, len(subnets))
		for i := range subnets {
			s[i] = aws.String(subnets[i])
		}

		awsvpcConfig := &ecs.AwsVpcConfiguration{}
		awsvpcConfig.SetAssignPublicIp(ecs.AssignPublicIpDisabled)
		awsvpcConfig.SetSecurityGroups(securityGroupIDs)
		awsvpcConfig.SetSubnets(s)

		networkConfig := &ecs.NetworkConfiguration{}
		networkConfig.SetAwsvpcConfiguration(awsvpcConfig)

		input.SetNetworkConfiguration(networkConfig)
		input.SetPlatformVersion(fargatePlatformVersion)
	}

	input.SetLaunchType(launchType)

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
