package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Read(serviceID string) (*models.Service, error) {
	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		return nil, err
	}

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	ecsService, err := s.readService(clusterName, serviceID)
	if err != nil {
		return nil, err
	}

	var deployments []models.Deployment
	for _, d := range ecsService.Deployments {
		taskDefinitionARN := aws.StringValue(d.TaskDefinition)
		deployID, err := s.lookupDeployIDFromTaskDefinitionARN(taskDefinitionARN)
		if err != nil {
			return nil, err
		}

		deployment, err := s.makeDeploymentModel(deployID)
		if err != nil {
			return nil, err
		}

		deployment.Created = aws.TimeValue(d.CreatedAt)
		deployment.DesiredCount = int(aws.Int64Value(d.DesiredCount))
		deployment.PendingCount = int(aws.Int64Value(d.PendingCount))
		deployment.RunningCount = int(aws.Int64Value(d.RunningCount))
		deployment.Status = aws.StringValue(d.Status)
		deployment.Updated = aws.TimeValue(d.UpdatedAt)

		deployments = append(deployments, *deployment)
	}

	var loadBalancerID string
	if len(ecsService.LoadBalancers) != 0 {
		loadBalancer := ecsService.LoadBalancers[0]
		loadBalancerName := aws.StringValue(loadBalancer.LoadBalancerName)
		fqLoadBalancerID := loadBalancerName
		loadBalancerID = delLayer0Prefix(s.Config.Instance(), fqLoadBalancerID)
	}

	model, err := s.makeServiceModel(environmentID, loadBalancerID, serviceID)
	if err != nil {
		return nil, err
	}

	model.Deployments = deployments
	model.DesiredCount = int(aws.Int64Value(ecsService.DesiredCount))
	model.PendingCount = int(aws.Int64Value(ecsService.PendingCount))
	model.RunningCount = int(aws.Int64Value(ecsService.RunningCount))

	return model, nil
}

func (s *ServiceProvider) readService(clusterName, serviceID string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	input.SetServices([]*string{
		aws.String(serviceID),
	})

	output, err := s.AWS.ECS.DescribeServices(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ServiceNotFoundException" {
			return nil, errors.Newf(errors.ServiceDoesNotExist, "Service '%s' does not exist", serviceID)
		}

		return nil, err
	}

	if len(output.Services) == 0 {
		return nil, errors.Newf(errors.ServiceDoesNotExist, "Service '%s' does not exist", serviceID)
	}

	return output.Services[0], nil
}

func (s *ServiceProvider) makeDeploymentModel(deployID string) (*models.Deployment, error) {
	model := &models.Deployment{
		DeployID: deployID,
	}

	tags, err := s.TagStore.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.DeployName = tag.Value
	}

	if tag, ok := tags.WithKey("version").First(); ok {
		model.DeployVersion = tag.Value
	}

	return model, nil
}

func (s *ServiceProvider) makeServiceModel(environmentID, loadBalancerID, serviceID string) (*models.Service, error) {
	model := &models.Service{
		EnvironmentID:  environmentID,
		LoadBalancerID: loadBalancerID,
		ServiceID:      serviceID,
	}

	tags, err := s.TagStore.SelectByTypeAndID("service", serviceID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.ServiceName = tag.Value
	}

	tags, err = s.TagStore.SelectByTypeAndID("environment", environmentID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.EnvironmentName = tag.Value
	}

	if loadBalancerID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("load_balancer", loadBalancerID)
		if err != nil {
			return nil, err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.LoadBalancerName = tag.Value
		}
	}

	return model, nil
}
