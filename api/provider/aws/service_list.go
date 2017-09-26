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
			s.populateSummariesTags(service, ecsServiceID)

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

func (s *ServiceProvider) populateSummariesTags(service *models.Service, serviceID string) error {
	return nil
}

func (s *ServiceProvider) serviceARNToECSServiceID(arn string) string {
	split := strings.SplitN(arn, "/", -1)
	serviceName := split[len(split)-1]
	return serviceName
}
