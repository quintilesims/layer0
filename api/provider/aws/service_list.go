package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) List() ([]models.ServiceSummary, error) {
	clusterArns, err := s.listClusters()
	if err != nil {
		return nil, err
	}

	clusterServices := map[*string][]*string{}

	for _, clusterArn := range clusterArns {
		clusterName := strings.Split(aws.StringValue(clusterArn), ":")[5]
		clusterName = strings.Replace(clusterName, "cluster/", "", -1)

		if !strings.HasPrefix(clusterName, s.Config.Instance()) {
			continue
		}

		serviceArns, err := s.listServices(clusterArn)
		if err != nil {
			return nil, err
		}

		clusterServices[clusterArn] = serviceArns
	}

	serviceSummaries := []models.ServiceSummary{}
	for _, serviceArns := range clusterServices {
		for _, serviceArn := range serviceArns {
			ecsServiceID := s.serviceARNToECSServiceID(aws.StringValue(serviceArn))
			service := &models.Service{
				ServiceID: ecsServiceID,
			}
			s.updateWithTagInfo(service, ecsServiceID)

			summary := models.ServiceSummary{
				ServiceID:       service.ServiceID,
				ServiceName:     service.ServiceName,
				EnvironmentID:   service.EnvironmentID,
				EnvironmentName: service.EnvironmentName,
			}

			serviceSummaries = append(serviceSummaries, summary)
		}
	}

	return serviceSummaries, nil
}

func (s *ServiceProvider) listServices(clusterArn *string) ([]*string, error) {
	serviceArns := []*string{}
	input := &ecs.ListServicesInput{
		Cluster: clusterArn,
	}

	for {
		output, err := s.AWS.ECS.ListServices(input)
		if err != nil {
			return nil, err
		}

		serviceArns = append(serviceArns, output.ServiceArns...)
		input.NextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}

	return serviceArns, nil
}

func (s *ServiceProvider) listClusters() ([]*string, error) {
	input := &ecs.ListClustersInput{}
	clusterArns := []*string{}

	for {
		output, err := s.AWS.ECS.ListClusters(input)
		if err != nil {
			return nil, err
		}

		clusterArns = append(clusterArns, output.ClusterArns...)
		input.NextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}

	return clusterArns, nil
}

func (s *ServiceProvider) updateWithTagInfo(service *models.Service, serviceID string) error {
	tags, err := s.TagStore.SelectByTypeAndID("service", serviceID)
	if err != nil {
		return err
	}
	service.ServiceID = serviceID

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		service.EnvironmentID = tag.Value
	}

	if tag, ok := tags.WithKey("load_balancer_id").First(); ok {
		service.LoadBalancerID = tag.Value
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		service.ServiceName = tag.Value
	}

	if service.EnvironmentID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("environment", service.EnvironmentID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			service.EnvironmentName = tag.Value
		}
	}

	if service.LoadBalancerID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("load_balancer", service.LoadBalancerID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			service.LoadBalancerName = tag.Value
		}
	}

	deployments := []models.Deployment{}
	for _, deploy := range service.Deployments {
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

	service.Deployments = deployments

	return nil
}

func (s *ServiceProvider) serviceARNToECSServiceID(arn string) string {
	split := strings.SplitN(arn, "/", -1)
	serviceName := split[len(split)-1]
	return serviceName
}
