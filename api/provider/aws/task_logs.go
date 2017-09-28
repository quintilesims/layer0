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

	logGroupName := t.Config.LogGroupName()
	return GetLogsFromTaskARNs(t.AWS.CloudWatchLogs, logGroupName, []string{taskARN}, tail, start, end)
}
