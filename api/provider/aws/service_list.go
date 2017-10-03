package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of ECS Services and returns a list of Service Summaries.
// A Service Summary consists of the Service ID, Service name, Environment ID, and
// Environment name.
func (s *ServiceProvider) List() ([]models.ServiceSummary, error) {
	clusterNames, err := listClusterNames(s.AWS.ECS, s.Config.Instance())
	if err != nil {
		return nil, err
	}

	serviceARNs, err := s.listClusterServiceARNs(clusterNames)
	if err != nil {
		return nil, err
	}

	serviceIDs := make([]string, len(serviceARNs))
	for i, serviceARN := range serviceARNs {
		// sample service ARN:
		// arn:aws:ecs:us-west-2:856306994068:service/l0-tlakedev-guestbo80d9d
		fqServiceID := strings.Split(serviceARN, "/")[1]
		serviceID := delLayer0Prefix(s.Config.Instance(), fqServiceID)
		serviceIDs[i] = serviceID
	}

	return s.makeServiceSummaryModels(serviceIDs)
}

func (s *ServiceProvider) listClusterServiceARNs(clusterNames []string) ([]string, error) {
	var serviceARNs []string

	fn := func(output *ecs.ListServicesOutput, lastPage bool) bool {
		for _, serviceARN := range output.ServiceArns {
			serviceARNs = append(serviceARNs, aws.StringValue(serviceARN))
		}

		return !lastPage
	}

	for _, clusterName := range clusterNames {
		input := &ecs.ListServicesInput{}
		input.SetCluster(clusterName)
		if err := s.AWS.ECS.ListServicesPages(input, fn); err != nil {
			return nil, err
		}
	}

	return serviceARNs, nil
}

func (s *ServiceProvider) makeServiceSummaryModels(serviceIDs []string) ([]models.ServiceSummary, error) {
	serviceTags, err := s.TagStore.SelectByType("service")
	if err != nil {
		return nil, err
	}

	models := make([]models.ServiceSummary, len(serviceIDs))
	for i, serviceID := range serviceIDs {
		models[i].ServiceID = serviceID

		if tag, ok := serviceTags.WithID(serviceID).WithKey("name").First(); ok {
			models[i].ServiceName = tag.Value
		}

		var environmentID string
		if tag, ok := serviceTags.WithID(serviceID).WithKey("environment_id").First(); ok {
			environmentID = tag.Value
			models[i].EnvironmentID = environmentID
		}

		environmentTags, err := s.TagStore.SelectByTypeAndID("environment", environmentID)
		if err != nil {
			return nil, err
		}

		if tag, ok := environmentTags.WithKey("name").First(); ok {
			models[i].EnvironmentName = tag.Value
		}
	}

	return models, nil
}
