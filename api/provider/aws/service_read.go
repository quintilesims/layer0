package aws

import (
	"github.com/aws/aws-sdk-go/aws"
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
	for _, deploy := range ecsService.Deployments {
		deployID := aws.StringValue(deploy.TaskDefinition)

		//todo: convert deployid to layer0 deploy id

		d := models.Deployment{
			Created:       aws.TimeValue(deploy.CreatedAt),
			DeployID:      deployID,
			DeployName:    "", // tag
			DeployVersion: "", // tag
			DesiredCount:  int(aws.Int64Value(deploy.DesiredCount)),
			PendingCount:  int(aws.Int64Value(deploy.PendingCount)),
			RunningCount:  int(aws.Int64Value(deploy.RunningCount)),
			Status:        aws.StringValue(deploy.Status),
			Updated:       aws.TimeValue(deploy.UpdatedAt),
		}

		deployments = append(deployments, d)
	}

	var loadBalancerID string
	if len(ecsService.LoadBalancers) != 0 {
		loadBalancer := ecsService.LoadBalancers[0]
		loadBalancerName := aws.StringValue(loadBalancer.LoadBalancerName)
		fqLoadBalancerID := loadBalancerName
		loadBalancerID = delLayer0Prefix(s.Config.Instance(), fqLoadBalancerID)
	}

	model := &models.Service{
		Deployments:      deployments,
		DesiredCount:     int(aws.Int64Value(ecsService.DesiredCount)),
		EnvironmentID:    environmentID,
		EnvironmentName:  "", // tag
		LoadBalancerID:   loadBalancerID,
		LoadBalancerName: "", // tag
		PendingCount:     int(aws.Int64Value(ecsService.PendingCount)),
		RunningCount:     int(aws.Int64Value(ecsService.RunningCount)),
		ServiceID:        serviceID,
		ServiceName:      "", // tag
	}

	if err := s.populateModelTags(serviceID, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ServiceProvider) readService(clusterName, serviceID string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	input.SetServices([]*string{
		aws.String(serviceID),
	})

	output, err := s.AWS.ECS.DescribeServices(input)
	if err != nil || len(output.Services) == 0 {
		return nil, errors.Newf(errors.ServiceDoesNotExist, "Service '%s' does not exist")
	}

	return output.Services[0], nil
}

func (s *ServiceProvider) populateModelTags(serviceID string, model *models.Service) error {
	tags, err := s.TagStore.SelectByTypeAndID("service", serviceID)
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		model.EnvironmentID = tag.Value
	}

	if tag, ok := tags.WithKey("load_balancer_id").First(); ok {
		model.LoadBalancerID = tag.Value
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.ServiceName = tag.Value
	}

	if model.EnvironmentID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("environment", model.EnvironmentID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.EnvironmentName = tag.Value
		}
	}

	if model.LoadBalancerID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("load_balancer", model.LoadBalancerID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.LoadBalancerName = tag.Value
		}
	}

	deployments := []models.Deployment{}
	for _, deploy := range model.Deployments {
		tags, err := s.TagStore.SelectByTypeAndID("deploy", deploy.DeployID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			deploy.DeployName = tag.Value
		}

		if tag, ok := tags.WithKey("version").First(); ok {
			deploy.DeployVersion = tag.Value
		}

		deployments = append(deployments, deploy)
	}

	model.Deployments = deployments

	return nil
}
