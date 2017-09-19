package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Update(req models.UpdateServiceRequest) error {
	serviceID := req.ServiceID
	deployID := req.DeployID
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)

	serviceTags, err := s.TagStore.SelectByTypeAndID("service", serviceID)
	if err != nil {
		return err
	}

	tag, ok := serviceTags.WithKey("environment_id").First()
	if !ok {
		// is this canonical error handling?
		return fmt.Errorf("Cannot resolve environment_id for service %s", serviceID)
	}

	environmentID := tag.Value
	clusterName := addLayer0Prefix(s.Config.Instance(), environmentID)

	service, err := s.readService(clusterName, serviceID)
	if err != nil {
		return err
	}

	var serviceScaleCount *int64
	if req.ServiceScaleCount == nil {
		serviceScaleCount = service.DesiredCount
	} else {
		s := int64(*req.ServiceScaleCount)
		serviceScaleCount = &s
	}

	if err := s.updateService(clusterName, fqServiceID, deployID, serviceScaleCount); err != nil {
		return err
	}

	return nil
}

func (s *ServiceProvider) updateService(cluster, service string, taskDefinition *string, desiredCount *int64) error {
	input := &ecs.UpdateServiceInput{}
	input.SetCluster(cluster)
	input.SetService(service)
	if taskDefinition != nil {
		input.SetTaskDefinition(*taskDefinition)
	}

	if desiredCount != nil {
		input.SetDesiredCount(*desiredCount)
	}

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := s.AWS.ECS.UpdateService(input); err != nil {
		return err
	}

	return nil
}
