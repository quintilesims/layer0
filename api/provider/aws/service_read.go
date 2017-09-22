package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Read(serviceID string) (*models.Service, error) {
	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
		return nil, err
	}

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	service, err := s.readService(clusterName, serviceID)
	if err != nil {
		return nil, err
	}

	model := &models.Service{
		DesiredCount: aws.Int64Value(service.DesiredCount),
		RunningCount: aws.Int64Value(service.RunningCount),
		PendingCount: aws.Int64Value(service.PendingCount),
	}

	for _, deploy := range service.Deployments {
		deployID := aws.StringValue(deploy.TaskDefinition)
		//todo: convert deployid to layer0 deploy id

		deploy := models.Deployment{
			DeploymentID: aws.StringValue(deploy.Id),
			Created:      aws.TimeValue(deploy.CreatedAt),
			Updated:      aws.TimeValue(deploy.UpdatedAt),
			Status:       aws.StringValue(deploy.Status),
			PendingCount: aws.Int64Value(deploy.PendingCount),
			RunningCount: aws.Int64Value(deploy.RunningCount),
			DesiredCount: aws.Int64Value(deploy.DesiredCount),
			DeployID:     deployID,
		}

		model.Deployments = append(model.Deployments, deploy)
	}

	if err := s.updateWithTagInfo(model, serviceID); err != nil {
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
		return nil, err
	}

	if len(output.Services) == 0 {
		return nil, fmt.Errorf("ECS service '%s' in cluster '%s' does not exist", serviceID, clusterName)
	}

	return output.Services[0], nil
}
