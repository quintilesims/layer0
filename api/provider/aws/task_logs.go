package aws

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

func (t *TaskProvider) Logs(taskID string, tail int, start, end time.Time) ([]models.LogFile, error) {
	taskARN, err := t.lookupTaskARN(taskID)
	if err != nil {
		return nil, err
	}

	environmentID, err := lookupEntityEnvironmentID(t.TagStore, "task", taskID)
	if err != nil {
		return nil, err
	}

	fqEnvironmentID := addLayer0Prefix(t.Config.Instance(), environmentID)
	clusterName := fqEnvironmentID
	logGroupName := t.Config.LogGroupName()

	return GetLogsFromTaskARNs(t.AWS.CloudWatchLogs, logGroupName, clusterName, []string{taskARN}, tail, start, end)
}
