package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/common/models"
)

func (s *ServiceProvider) Logs(serviceID string, tail int, start, end time.Time) ([]models.LogFile, error) {
	environmentID, err := lookupEntityEnvironmentID(s.TagStore, "service", serviceID)
	if err != nil {
			return nil, err
	}

	fqEnvironmentID := addLayer0Prefix(s.Config.Instance(), environmentID)
	fqServiceID := addLayer0Prefix(s.Config.Instance(), serviceID)
	clusterName := fqEnvironmentID
	service, err := s.readService(clusterName, fqServiceID)
	if err != nil {
		return nil, err
	}

	serviceTaskARNs := []string{}
	for _, deployment := range service.Deployments {
		startedBy := aws.StringValue(deployment.Id)
		deploymentTaskARNs, err := listClusterTaskARNs(s.AWS.ECS, clusterName, startedBy)
		if err != nil {
			return nil, err
		}

		serviceTaskARNs = append(serviceTaskARNs, deploymentTaskARNs...)
	}

	logGroupName := s.Config.LogGroupName()
	return GetLogsFromTaskARNs(s.AWS.CloudWatchLogs, logGroupName, serviceTaskARNs, tail, start, end)
}
