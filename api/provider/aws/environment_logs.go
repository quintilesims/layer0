package aws

import (
	"fmt"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Logs(environmentID string, tail int, start, end time.Time) ([]models.LogFile, error) {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	logGroupName := e.Config.LogGroupName()
	filterPattern := fmt.Sprintf("{ $.eventSource = \"ecs.amazonaws.com\" && $.requestParameters.cluster = \"%s\" }", clusterName)

	return GetLogsFromCloudTrail(e.AWS.CloudWatchLogs, logGroupName, tail, start, end, filterPattern)
}
