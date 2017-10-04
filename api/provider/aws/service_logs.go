package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
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

	taskARNs := []string{}
	for _, deployment := range service.Deployments {
		startedBy := aws.StringValue(deployment.Id)
		clusterTaskARNsStopped, err := listClusterTaskARNs(s.AWS.ECS, clusterName, startedBy, ecs.DesiredStatusStopped)
		if err != nil {
			return nil, err
		}

		clusterTaskARNsRunning, err := listClusterTaskARNs(s.AWS.ECS, clusterName, startedBy, ecs.DesiredStatusRunning)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, clusterTaskARNsStopped...)
		taskARNs = append(taskARNs, clusterTaskARNsRunning...)
	}

	logGroupName := s.Config.LogGroupName()
	return GetLogsFromTaskARNs(s.AWS.CloudWatchLogs, logGroupName, taskARNs, tail, start, end)
}
