package aws

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Create(req models.CreateServiceRequest) (*models.Service, error) {
	serviceID := generateEntityID(req.ServiceName)

	fqDeployID := addLayer0Prefix(s.Config.Instance(), req.DeployID)
	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), req.EnvironmentID)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	desiredCount := 1

	var (
		loadBalancerContainers []ecs.LoadBalancer
		loadBalancerRole       string
	)

	if req.LoadBalancerID != "" {
		fqLoadbalancerID := addLayer0Prefix(s.Config.Instance(), req.LoadBalancerID)
		loadBalancerContainer, err := s.getLoadBalancerContainer(fqLoadBalancerID, fqDeployID)
		if err != nil {
			return nil, err
		}

		loadBalancerContainers = []ecs.LoadBalancer{loadBalancerContainer}
		loadBalancerRole = fmt.Sprintf("%s-lb", fqLoadBalancerID)
	}

	if err := s.createService(); err != nil {
		return nil, err
	}

	return s.Read(serviceID)
}

func (s *ServiceProvider) createService(cluster, serviceName, taskDefinition, loadBalancerRole string, loadBalancers []ecs.LoadBalancer) error {
	input := &ecs.CreateServiceInput{}
	input.SetCluster(cluster)
	input.SetDesiredCount(aws.Int64(desiredCount))
	input.SetServiceName(serviceName)
	if len(loadbalancers) > 0 {
		input.SetLoadBalancers

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.CreateService(input); err != nil {
		return err
	}

	return nil
}
