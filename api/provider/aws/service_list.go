package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

// List retrieves a list of ECS Services and returns a list of Service Summaries.
// A Service Summary consists of the Service ID, Service name, Environment ID, and
// Environment name.
func (s *ServiceProvider) List() ([]models.ServiceSummary, error) {
	serviceARNs, err := s.listServiceARNs()
	if err != nil {
		return nil, err
	}

	serviceIDs := make([]string, len(serviceARNs))
	for i, serviceARN := range serviceARNs {
		serviceID, err := lookupServiceIDFromServiceARN(s.TagStore, serviceARN)
		if err != nil {
			return nil, err
		}

		serviceIDs[i] = serviceID
	}

	return s.makeServiceSummaryModels(serviceIDs)
}

func (s *ServiceProvider) listServiceARNs() ([]string, error) {
	var serviceARNs []string
	fn := func(output *ecs.ListServicesOutput, lastPage bool) bool {
		for _, serviceARN := range output.ServiceArns {
			serviceARNs = append(serviceARNs, aws.StringValue(serviceARN))
		}

		return !lastPage
	}

	if err := s.AWS.ECS.ListServicesPages(&ecs.ListServicesInput{}, fn); err != nil {
		return nil, err
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

		environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
		if err != nil {
			return nil, err
		}

		models[i].EnvironmentID = environmentID

		environmentTags, err := s.TagStore.SelectByTypeAndID("environment", environmentID)
		if err != nil {
			return nil, err
		}

		tag, ok := environmentTags.WithKey("name").First()
		if !ok {
			return nil, errors.Newf(errors.EnvironmentDoesNotExist, "Could not resolve name for environment with id '%s'", environmentID)
		}

		models[i].EnvironmentName = tag.Value
	}

	return models, nil
}
