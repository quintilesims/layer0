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

	taskDefinitionARN, err := lookupTaskDefinitionARNFromDeployID(s.TagStore, req.DeployID)
	if err != nil {
		return "", err
	}

	taskDefinition, err := describeTaskDefinition(s.AWS.ECS, taskDefinitionARN)
	if err != nil {
		return "", err
	}

	networkMode := aws.StringValue(taskDefinition.NetworkMode)

	var securityGroupIDs []*string
	var subnets []string
	if networkMode == ecs.NetworkModeAwsvpc {
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
		lb, err := describeLoadBalancer(s.AWS.ELB, s.AWS.ALB, fqLoadBalancerID)
		if err != nil {
			return "", err
		}

		if aws.StringValue(lb.GetScheme()) == "internet-facing" {
			loadBalancerSecurityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
			loadBalancerSecurityGroup, err := readSG(s.AWS.EC2, loadBalancerSecurityGroupName)
			if err != nil {
				return "", err
			}

			securityGroupIDs = append(securityGroupIDs, loadBalancerSecurityGroup.GroupId)
		}

		var assignLoadBalancer = func() error {
			for _, containerDefinition := range taskDefinition.ContainerDefinitions {
				for _, portMapping := range containerDefinition.PortMappings {
					var loadBalancerName *string
					var targetGroupArn *string

					// todo: check if this verifying this is actually needed
					// why do we need to verify at least one container exposes a port that the elb
					// has a listener for?
					if lb.isELB {
						loadBalancerName = lb.ELB.LoadBalancerName
						// for _, listenerDescription := range lb.ELB.ListenerDescriptions {
						// 	if aws.Int64Value(listenerDescription.Listener.InstancePort) != aws.Int64Value(portMapping.ContainerPort) {
						// 		continue
						// 	}
						// }
					}

					// if load balancer is an application load balancer we need to assign TargetGroupArn
					// instead of the LoadBalancerName property of ecs.LoadBalancer
					if lb.isALB {
						targetGroupName := fqLoadBalancerID
						targetGroup, err := getTargetGroupArn(s.AWS.ALB, targetGroupName)
						if err != nil {
							return err
						}

						targetGroupArn = targetGroup.TargetGroupArn
					}

					loadBalancer = &ecs.LoadBalancer{
						ContainerName:    containerDefinition.Name,
						ContainerPort:    portMapping.ContainerPort,
						LoadBalancerName: loadBalancerName,
						TargetGroupArn:   targetGroupArn,
					}

					loadBalancerRole = fmt.Sprintf("%s-lb", fqLoadBalancerID)

					return nil
				}
			}

			return nil
		}

		if err := assignLoadBalancer(); err != nil {
			return "", err
		}
	}

	if err := s.createService(
		cluster,
		serviceName,
		taskDefinitionARN,
		loadBalancerRole,
		networkMode,
		req.Stateful,
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
	serviceName,
	taskDefinitionARN,
	loadBalancerRole,
	networkMode string,
	stateful bool,
	desiredCount int,
	subnets []string,
	securityGroupIDs []*string,
	loadBalancer *ecs.LoadBalancer,
) error {
	input := &ecs.CreateServiceInput{}
	input.SetCluster(cluster)
	input.SetDesiredCount(int64(desiredCount))
	input.SetServiceName(serviceName)
	input.SetTaskDefinition(taskDefinitionARN)

	launchType := ecs.LaunchTypeEc2
	if !stateful {
		launchType = ecs.LaunchTypeFargate
		input.SetPlatformVersion(config.DefaultFargatePlatformVersion)
	}

	input.SetLaunchType(launchType)

	if networkMode == ecs.NetworkModeAwsvpc {
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
	}
	if loadBalancer != nil {
		input.SetLoadBalancers([]*ecs.LoadBalancer{loadBalancer})
		if networkMode != ecs.NetworkModeAwsvpc {
			input.SetRole(loadBalancerRole)
		}
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
