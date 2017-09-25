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

		deployment := models.Deployment{
			Created:       aws.TimeValue(d.CreatedAt),
			DeployID:      deployID,
			DeployName:    "", // tag
			DeployVersion: "", // tag
			DesiredCount:  int(aws.Int64Value(d.DesiredCount)),
			PendingCount:  int(aws.Int64Value(d.PendingCount)),
			RunningCount:  int(aws.Int64Value(d.RunningCount)),
			Status:        aws.StringValue(d.Status),
			Updated:       aws.TimeValue(d.UpdatedAt),
		}

		deployments = append(deployments, deployment)
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

func (s *ServiceProvider) lookupDeployIDFromTaskDefinitionARN(taskDefinitionARN string) (string, error) {
	tags, err := s.TagStore.SelectByType("deploy")
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.Newf(errors.DeployDoesNotExist, "No deploys exist")
	}

	if tag, ok := tags.WithValue(taskDefinitionARN).First(); ok {
		return tag.EntityID, nil
	}

	return "", errors.Newf(errors.DeployDoesNotExist, "Failed to find deploy with ARN '%s'", taskDefinitionARN)
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

	for _, d := range model.Deployments {
		tags, err := s.TagStore.SelectByTypeAndID("deploy", d.DeployID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("deploy_name").First(); ok {
			d.DeployName = tag.Value
		}

		if tag, ok := tags.WithKey("version").First(); ok {
			d.DeployVersion = tag.Value
		}
	}

	return nil
}
