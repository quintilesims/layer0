package aws

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/retry"
)

// Create runs an ECS Service using the specified CreateServiceRequest.
// The CreateServiceRequest contains the DeployID, the EnvironmentID,
// the LoadBalancerID (optional), the Scale, the ServiceName, and a
// boolean Stateful value.
//
// The DeployID is used to look up the ARN of the ECS TaskDefinition to run. The
// Stateful boolean indicates which ECS LaunchType the user wishes to use
// ("FARGATE" if false, "EC2" if true).
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

	taskDefinitionCompatibilities := make([]string, len(taskDefinition.Compatibilities))
	for i, _ := range taskDefinition.Compatibilities {
		taskDefinitionCompatibilities[i] = aws.StringValue(taskDefinition.Compatibilities[i])
	}

	if !req.Stateful && !stringInSlice(ecs.LaunchTypeFargate, taskDefinitionCompatibilities) {
		errMsg := "Cannot create stateless service using stateful deploy '%s'."
		return "", fmt.Errorf(errMsg, req.DeployID)
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

		if req.Stateful && lb.isALB {
			errMsg := "Cannot create stateful service behind application load balancer '%s'. "
			errMsg += "Use a classic load balancer instead."
			return "", fmt.Errorf(errMsg, fqLoadBalancerID)
		}

		if !req.Stateful && lb.isCLB {
			errMsg := "Cannot create stateless service behind classic load balancer '%s'. "
			errMsg += "Use an application load balancer instead."
			return "", fmt.Errorf(errMsg, fqLoadBalancerID)
		}

		if lb.isCLB && stringInSlice(ecs.LaunchTypeFargate, taskDefinitionCompatibilities) {
			errMsg := "Cannot deploy stateless deploy '%s' behind classic load balancer '%s'. "
			errMsg += "Either use a stateful deploy with a classic load balancer, "
			errMsg += "or a stateless deploy with an application load balancer."
			return "", fmt.Errorf(errMsg, req.DeployID, fqLoadBalancerID)
		}

		if lb.Scheme() == "internet-facing" {
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

					if lb.isCLB {
						loadBalancerName = lb.CLB.LoadBalancerName
					}

					// if load balancer is an application load balancer we need to assign TargetGroupArn
					// instead of the LoadBalancerName property of ecs.LoadBalancer
					if lb.isALB {
						targetGroupName := fqLoadBalancerID
						targetGroup, err := readTargetGroup(s.AWS.ALB, aws.String(targetGroupName), nil)
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

	fn := func() (shouldRetry bool, err error) {
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
			if strings.Contains(err.Error(), "Unable to assume role") {
				log.Printf("[DEBUG] Failed service create, will retry (%v)", err)
				return true, nil
			}

			log.Printf("[DEBUG] %s", reflect.TypeOf(err))
			return false, err
		}

		return false, nil
	}

	if err := retry.Retry(fn, retry.WithTimeout(time.Second*30), retry.WithDelay(time.Second)); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "InvalidParameterException" {
			return "", err
		}

		return "", errors.New(errors.EventualConsistencyError, err)
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
