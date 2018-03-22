package aws

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Logs(environmentID string, tail int, start, end time.Time) ([]models.LogFile, error) {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	logGroupName := e.Config.LogGroupName()

	return GetEnvironmentLogsFromCloudTrail(e.AWS.CloudWatchLogs, logGroupName, clusterName, tail, start, end)
}
